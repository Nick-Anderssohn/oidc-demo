package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
	"github.com/Nick-Anderssohn/oidc-demo/internal/session"
	"github.com/Nick-Anderssohn/oidc-demo/internal/user"
)

type Handlers struct {
	DepResolver *deps.Resolver
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	userSVC := user.Service{
		Resolver: h.DepResolver,
	}

	userID, err := session.UserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		http.Error(w, "Failed to get user ID from context", http.StatusInternalServerError)
		return
	}
	userData, err := userSVC.GetUserData(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(userData); err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		return
	}
}
