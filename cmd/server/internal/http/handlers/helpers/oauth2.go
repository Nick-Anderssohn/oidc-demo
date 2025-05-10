package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
	"github.com/Nick-Anderssohn/oidc-demo/internal/oidc"
	"golang.org/x/oauth2"
)

type AuthRedirectConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	DiscoveryURL string
	Scopes       []string
}

func RedirectToAuthorizationServer(
	depResolver *deps.Resolver,
	config *AuthRedirectConfig,
	w http.ResponseWriter,
	r *http.Request,
) {
	discoveryData, err := oidc.GetDiscoveryData(config.DiscoveryURL)
	if err != nil {
		http.Error(w, "Configuration error", http.StatusInternalServerError)
		return
	}

	oauthConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,

		Endpoint: oauth2.Endpoint{
			AuthURL:  discoveryData.AuthorizationEndpoint,
			TokenURL: discoveryData.TokenEndpoint,
		},

		RedirectURL: config.RedirectURL,
		Scopes:      config.Scopes,
	}

	stateToken := generateStateToken()

	err = depResolver.Queries.InsertStateToken(r.Context(), stateToken)
	if err != nil {
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

func generateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
