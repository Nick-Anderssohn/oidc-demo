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

func (h *Handlers) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, err := session.UserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		http.Error(w, "Failed to get user ID from context", http.StatusInternalServerError)
		return
	}

	err = h.DepResolver.Queries.DeleteUser(r.Context(), userID)
	if err != nil {
		log.Printf("could not delete user %v", userID.String())
		http.Error(w, "could not delete user", http.StatusInternalServerError)
		return
	}

	sessionSVC := session.Service{}
	sessionSVC.DeleteSessionCookie(w)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	sessionSVC := session.Service{
		Resolver: h.DepResolver,
	}

	err := sessionSVC.Logout(r.Context(), w)
	if err != nil {
		log.Printf("Failed to logout: %v", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
