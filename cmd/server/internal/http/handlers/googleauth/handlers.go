package googleauth

import (
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/helpers"
	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
)

type Handlers struct {
	DepResolver *deps.Resolver
}

func (h *Handlers) RedirectToAuthorizationServer(w http.ResponseWriter, r *http.Request) {
	cfg := createGoogleOIDCConfig(h.DepResolver)

	helpers.RedirectToAuthorizationServer(
		h.DepResolver,
		&cfg,
		w,
		r,
	)
}

func (h *Handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	config := createGoogleOIDCConfig(h.DepResolver)

	helpers.HandleOIDCCallback(
		h.DepResolver,
		&config,
		w,
		r,
	)

	// Redirect back to home page.
	http.Redirect(w, r, h.DepResolver.Config.APIConfig.BaseURL+"/", http.StatusFound)
}

func createGoogleOIDCConfig(depResolver *deps.Resolver) helpers.OIDCConfig {
	googleCfg := depResolver.Config.GoogleOIDCConfig

	return helpers.OIDCConfig{
		ClientID:     googleCfg.ClientID,
		ClientSecret: googleCfg.ClientSecret,
		RedirectURL:  depResolver.Config.APIConfig.BaseURL + "/callbacks/google",
		DiscoveryURL: "https://accounts.google.com/.well-known/openid-configuration",
		Scopes:       []string{"openid", "email"},
	}
}
