package main

import (
	"log"
	"net/http"

	"backend-go/api"

	"github.com/gorilla/handlers"
)

func main() {
	router := api.SetupRoutes()

	// Adaugăm middleware-ul CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),  // Permitem cereri de la Angular
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}), // Metode permise
		handlers.AllowedHeaders([]string{"Content-Type"}),           // Headere permise
	)(router)

	log.Println("Server pornit pe port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
