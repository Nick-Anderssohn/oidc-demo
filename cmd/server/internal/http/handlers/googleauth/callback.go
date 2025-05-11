package googleauth

import (
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/helpers"
	"github.com/Nick-Anderssohn/oidc-demo/internal/oidc"
	"github.com/Nick-Anderssohn/oidc-demo/internal/sqlc/dal"
	"github.com/jackc/pgx/v5"
)

func (h *Handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	log.Printf("Full URL: %s", r.URL.String())

	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("Header: %s=%s", name, value)
		}
	}

	for name, values := range r.URL.Query() {
		for _, value := range values {
			log.Printf("Query Param: %s=%s", name, value)
		}
	}

	// Validate state
	state := r.URL.Query().Get("state")

	if state == "" {
		// TODO: Direct user to error page.
		http.Error(w, "State parameter is missing", http.StatusBadRequest)
		return
	}

	_, err := h.DepResolver.Queries.GetStateToken(r.Context(), state)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("invalid state token: %s", state)
		} else {
			log.Printf("Failed to get state token: %v", err)
		}

		http.Error(w, "Could not validate state", http.StatusBadRequest)
		return
	}

	googleConfig := createGoogleOIDCConfig(h.DepResolver)
	oauthConfig, err := helpers.GetOauthConfig(&googleConfig)
	if err != nil {
		log.Printf("Failed to get OAuth config: %v", err)
		http.Error(w, "Failed to get OAuth config", http.StatusInternalServerError)
		return
	}

	oidcConfig := oidc.Config{
		Config:       &oauthConfig,
		DiscoveryURL: googleConfig.DiscoveryURL,
	}

	tokenResp, err := oidc.ExchangeCodeForToken(r.Context(), &oidcConfig, r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	log.Printf("Token Response: %+v", tokenResp)
	if tokenResp.IDToken != "" {
		log.Printf("ID Token: %s", tokenResp.IDToken)
	}
	if tokenResp.AccessToken != "" {
		log.Printf("Access Token: %s", tokenResp.AccessToken)
	}
	if tokenResp.RefreshToken != "" {
		log.Printf("Refresh Token: %s", tokenResp.RefreshToken)
	}
	log.Printf("Expiry: %v", tokenResp.Expiry)

	for k, v := range tokenResp.IDTokenPayload {
		log.Printf("IDTokenPayload: %s=%v", k, v)
	}

	email, ok := tokenResp.IDTokenPayload["email"].(string)
	if !ok || email == "" {
		log.Printf("Email not found in ID token payload")
		http.Error(w, "Email not found in ID token payload", http.StatusBadRequest)
		return
	}

	externalId, ok := tokenResp.IDTokenPayload["sub"].(string)
	if !ok || externalId == "" {
		log.Printf("External ID not found in ID token payload")
		http.Error(w, "External ID not found in ID token payload", http.StatusBadRequest)
		return
	}

	// For now, we'll just upsert the user and identity.
	queries := h.DepResolver.Queries
	user, err := queries.UpsertUserByEmail(r.Context(), tokenResp.IDTokenPayload["email"].(string))
	if err != nil {
		log.Printf("Failed to upsert user: %v", err)
		http.Error(w, "Failed to upsert user", http.StatusInternalServerError)
		return
	}

	_, err = queries.UpsertIdentity(r.Context(), dal.UpsertIdentityParams{
		UserID:             user.ID,
		IdentityProviderID: dal.IdentityProviderIDGoogle,
		ExternalID:         externalId,
	})

	if err != nil {
		log.Printf("Failed to upsert identity: %v", err)
		http.Error(w, "Failed to upsert identity", http.StatusInternalServerError)
		return
	}

	// Redirect back to home page.
	http.Redirect(w, r, h.DepResolver.Config.APIConfig.BaseURL+"/", http.StatusFound)
}
