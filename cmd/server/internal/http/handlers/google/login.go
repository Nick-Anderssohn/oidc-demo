package google

import (
	"net/http"
	"os"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/helpers"
	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
)

type Handlers struct {
	DepResolver *deps.Resolver
}

func (h *Handlers) RedirectToAuthorizationServer(w http.ResponseWriter, r *http.Request) {
	cfg := helpers.AuthRedirectConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callbacks/google",
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
