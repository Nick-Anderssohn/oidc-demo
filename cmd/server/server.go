package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/google"
	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
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

	// Serve static files
	router.Handle("/*", http.FileServer(http.Dir("./static")))

	registerAuthEndpoints(router, &resolver)

	// Routes under /private/api require the user to be authenticated
	router.Route("/private/api", func(r chi.Router) {
		r.Use(contentTypeJsonMiddleware)

		// TODO: Get identities.
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

func registerAuthEndpoints(router *chi.Mux, resolver *deps.Resolver) {
	googleHandlers := google.Handlers{
		DepResolver: resolver,
	}

	// Login redirects
	router.Route("/login", func(r chi.Router) {
		r.Get("/google", googleHandlers.RedirectToAuthorizationServer)
	})

	router.Route("/callbacks", func(r chi.Router) {
		r.Get("/google", func(w http.ResponseWriter, r *http.Request) {
			// Handle Google callback here
			json.NewEncoder(w).Encode(Message{Text: "Google callback received!"})
		})
	})
}
