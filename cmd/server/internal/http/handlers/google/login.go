package google

import (
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/helpers"
	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
)

type Handlers struct {
	DepResolver *deps.Resolver
}

func (h *Handlers) RedirectToAuthorizationServer(w http.ResponseWriter, r *http.Request) {
	googleCfg := h.DepResolver.Config.GoogleOIDCConfig

	cfg := helpers.AuthRedirectConfig{
		ClientID:     googleCfg.ClientID,
		ClientSecret: googleCfg.ClientSecret,
		RedirectURL:  h.DepResolver.Config.APIConfig.BaseURL + "/callbacks/google",
		DiscoveryURL: "https://accounts.google.com/.well-known/openid-configuration",
		Scopes:       []string{"openid", "email"},
	}

	helpers.RedirectToAuthorizationServer(
		h.DepResolver,
		&cfg,
		w,
		r,
	)
}
