package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend-go/internal/api"
	"backend-go/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Încărcăm variabilele de mediu din fișierul .env (pentru dezvoltare locală)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Inițializăm configurația
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Creem router-ul pentru API
	router := api.SetupRoutes(cfg)

	// Determinăm portul pentru server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port implicit
	}

	// Pornire server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
