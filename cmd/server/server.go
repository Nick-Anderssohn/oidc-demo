package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	log.Println("Starting server...")
	// Initialize the server

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Handle("/*", http.FileServer(http.Dir("./static")))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
