package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/cmd/server/internal/http/handlers/google"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	log.Println("Starting server...")
	// Initialize the server

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(contentTypeJsonMiddleware)

	router.Handle("/*", http.FileServer(http.Dir("./static")))

	router.Route("/api", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(Message{Text: "Welcome to the API!"})
		})

		r.Route("/google", func(r chi.Router) {
			r.Get("/discovery", google.GetDiscoveryData)
		})
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
