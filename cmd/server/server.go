package main

import (
	"encoding/json"
	"log"
	"net/http"

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

	router.Handle("/*", http.FileServer(http.Dir("./static")))

	router.Route("/api", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(Message{Text: "Welcome to the API!"})
		})
		r.Post("/echo", func(w http.ResponseWriter, r *http.Request) {
			var msg Message
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(msg)
		})
	})

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
