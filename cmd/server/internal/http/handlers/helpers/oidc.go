package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
	"github.com/Nick-Anderssohn/oidc-demo/internal/oidc"
	"github.com/Nick-Anderssohn/oidc-demo/internal/session"
	"github.com/Nick-Anderssohn/oidc-demo/internal/sqlc/dal"
	"github.com/Nick-Anderssohn/oidc-demo/internal/util"
	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	DiscoveryURL string
	Scopes       []string
}

func RedirectToAuthorizationServer(
	depResolver *deps.Resolver,
	config *OIDCConfig,
	w http.ResponseWriter,
	r *http.Request,
) {
	oauthConfig, err := getOIDCConfig(config)

	if err != nil {
		http.Error(w, "Configuration error", http.StatusInternalServerError)
		return
	}

	stateToken, err := util.GenerateSecureID()
	if err != nil {
		log.Printf("could not generate state token")
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
	}

	err = depResolver.Queries.InsertStateToken(r.Context(), stateToken)
	if err != nil {
		log.Printf("Failed to insert state token: %v", err)
		http.Error(w, "Failed to insert state token", http.StatusInternalServerError)
		return
	}

	// We're just going to reuse the state token for the nonce too
	// for this demo app. In a real app, you might want this to work
	// differently depending on how/why you are using oauth2.
	nonceOption := oauth2.SetAuthURLParam("nonce", stateToken)

	authUrl := oauthConfig.AuthCodeURL(stateToken, nonceOption)

	http.Redirect(w, r, authUrl, http.StatusFound)
}

func HandleOIDCCallback(
	depResolver *deps.Resolver,
	config *OIDCConfig,
	w http.ResponseWriter,
	r *http.Request,
) {
	err := validateState(depResolver, r)
	if err != nil {
		log.Printf("State validation failed: %v", err)
		http.Error(w, "State validation failed", http.StatusBadRequest)
		return
	}

	oidcConfig, err := getOIDCConfig(config)
	if err != nil {
		log.Printf("Failed to get OIDC config: %v", err)
		http.Error(w, "Failed to get OIDC config", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Code parameter is missing")
		http.Error(w, "Code parameter is missing", http.StatusBadRequest)
		return
	}

	tokenResp, err := oidc.ExchangeCodeForToken(r.Context(), &oidcConfig, code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	nonce, ok := tokenResp.IDTokenPayload["nonce"].(string)
	if !ok || nonce == "" {
		log.Println("missing nonce")
		http.Error(w, "missing nonce", http.StatusInternalServerError)
		return
	}

	if err = checkNonce(depResolver, r.Context(), nonce); err != nil {
		log.Println("nonce already seen before! you trying a replay attack?!")
		http.Error(w, "nonce already seen before! you trying a replay attack?!", http.StatusInternalServerError)
		return
	}

	user, err := upsertUserAndIdentity(depResolver, &tokenResp, r.Context())
	if err != nil {
		log.Printf("Failed to upsert user and identity: %v", err)
		http.Error(w, "Failed to upsert user and identity", http.StatusInternalServerError)
		return
	}

	sessionSVC := session.Service{Resolver: depResolver}
	err = sessionSVC.SaveNewSessionCookie(r.Context(), user.ID, w)
	if err != nil {
		log.Printf("Failed to save session cookie: %v", err)
		http.Error(w, "Failed to save session cookie", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func validateState(
	depResolver *deps.Resolver,
	r *http.Request,
) error {
	state := r.URL.Query().Get("state")

	if state == "" {
		return fmt.Errorf("state parameter is missing")
	}

	_, err := depResolver.Queries.GetStateToken(r.Context(), state)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("invalid state token: %s", state)
		} else {
			return fmt.Errorf("failed to get state token: %v", err)
		}
	}

	return nil
}

func checkNonce(depResolver *deps.Resolver, ctx context.Context, nonce string) error {
	queries := depResolver.Queries

	// We currently are using the state token to track what nonces we generated
	_, err := queries.GetStateToken(ctx, nonce)
	if err != nil {
		return fmt.Errorf("invalid nonce")
	}

	err = depResolver.Queries.InsertNonce(ctx, nonce)
	if err != nil {
		return fmt.Errorf("nonce already seen. this could be a replay attack!")
	}

	return nil
}

func getOIDCConfig(config *OIDCConfig) (oidc.Config, error) {
	discoveryData, err := oidc.GetDiscoveryData(config.DiscoveryURL)
	if err != nil {
		return oidc.Config{}, err
	}

	return oidc.Config{
		Config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,

			Endpoint: oauth2.Endpoint{
				AuthURL:  discoveryData.AuthorizationEndpoint,
				TokenURL: discoveryData.TokenEndpoint,
			},

			RedirectURL: config.RedirectURL,
			Scopes:      config.Scopes,
		},
		DiscoveryURL: config.DiscoveryURL,
	}, nil
}

func upsertUserAndIdentity(
	depResolver *deps.Resolver,
	tokenResp *oidc.TokenResponse,
	ctx context.Context,
) (dal.DemoUser, error) {
	queries := depResolver.Queries

	// Extract external ID from token payload
	externalID, ok := tokenResp.IDTokenPayload["sub"].(string)
	if !ok || externalID == "" {
		return dal.DemoUser{}, fmt.Errorf("external ID not found in ID token payload")
	}

	// Check if a user exists with the given external ID
	user, err := queries.GetUserByIdentityExternalID(ctx, dal.GetUserByIdentityExternalIDParams{
		IdentityProviderID: dal.IdentityProviderIDGoogle,
		ExternalID:         externalID,
	})

	existingUserFound := err == nil
	if err != nil && err != pgx.ErrNoRows {
		return dal.DemoUser{}, fmt.Errorf("failed to get user by external ID: %v", err)
	}

	// Check if the user is logged in
	loggedInUserID, err := session.UserIDFromContext(ctx)
	userIsLoggedIn := err == nil

	// If user is logged in, ensure the account is not already linked to another user
	if userIsLoggedIn && existingUserFound && user.ID != loggedInUserID {
		return dal.DemoUser{}, fmt.Errorf("account is already linked to another user")
	}

	// Make sure to grab the user if they are logged in via another account
	if userIsLoggedIn && !existingUserFound {
		user, err = queries.GetUser(ctx, loggedInUserID)
		if err != nil {
			return dal.DemoUser{}, fmt.Errorf("failed to get logged-in user: %v", err)
		}
	}

	// If no existing user is found and the user is not logged in, create a new user
	if !userIsLoggedIn && !existingUserFound {
		email, ok := tokenResp.IDTokenPayload["email"].(string)
		if !ok || email == "" {
			return dal.DemoUser{}, fmt.Errorf("email not found in ID token payload")
		}

		user, err = queries.UpsertUserByEmail(ctx, email)
		if err != nil {
			return dal.DemoUser{}, fmt.Errorf("failed to upsert user: %v", err)
		}
	}

	// Marshal ID token payload
	idTokenJSON, err := json.Marshal(tokenResp.IDTokenPayload)
	if err != nil {
		return dal.DemoUser{}, fmt.Errorf("failed to marshal ID token payload: %v", err)
	}

	// Upsert identity record
	if _, err := queries.UpsertIdentity(ctx, dal.UpsertIdentityParams{
		UserID:             user.ID,
		IdentityProviderID: dal.IdentityProviderIDGoogle,
		ExternalID:         externalID,
		MostRecentIDToken:  idTokenJSON,
	}); err != nil {
		return dal.DemoUser{}, fmt.Errorf("failed to upsert identity: %v", err)
	}

	return user, nil
}
