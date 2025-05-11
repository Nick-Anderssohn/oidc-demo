package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/api"
	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/googleauth"
	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
	"github.com/Nick-Anderssohn/oidc-demo/internal/session"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	log.Println("Starting server...")

	backgroundCtx := context.Background()

	resolver, err := deps.InitDepsResolver(backgroundCtx)
	if err != nil {
		panic(err)
	}
	defer resolver.Close()

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	sessionSVC := session.Service{Resolver: &resolver}
	router.Use(sessionSVC.SessionMiddleware)

	// Serve static files
	router.Handle("/*", http.FileServer(http.Dir("./static")))

	apiHandlers := api.Handlers{
		DepResolver: &resolver,
	}

	registerAuthEndpoints(router, &apiHandlers)

	// Routes under /private/api require the user to be authenticated
	router.Route("/private/api", func(r chi.Router) {
		r.Use(sessionSVC.RequireSessionMiddleware)
		r.Use(contentTypeJsonMiddleware)

		r.Get("/me", apiHandlers.Me)
	})

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func contentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func registerAuthEndpoints(router *chi.Mux, apiHandlers *api.Handlers) {
	googleHandlers := googleauth.Handlers{
		DepResolver: apiHandlers.DepResolver,
	}

	// Endpoints that handle redirecting to the authorization server
	router.Route("/login", func(r chi.Router) {
		r.Get("/google", googleHandlers.RedirectToAuthorizationServer)
	})

	// Endpoints that handle the callback from the authorization server
	router.Route("/callbacks", func(r chi.Router) {
		r.Get("/google", googleHandlers.HandleCallback)
	})

	// Endpoints that handle logging out
	router.Get("/logout", apiHandlers.Logout)
}
