package api

import (
	"net/http"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"backend-go/internal/api/handlers"
	"backend-go/internal/config"
	"backend-go/internal/services"
	"backend-go/internal/storage"
)

// SetupRoutes configurează rutele API
func SetupRoutes(cfg *config.Config) http.Handler {
	// Inițializare servicii
	blobStore, err := storage.NewBlobStorage(cfg)
	if err != nil {
		panic(err)
	}

	sqlStore, err := storage.NewSQLStorage(cfg)
	if err != nil {
		panic(err)
	}

	iaService := services.NewEntityExtractor(cfg)

	// Inițializare handlers
	fileHandler := handlers.NewFileHandler(cfg, blobStore, sqlStore, iaService)

	// Creare router
	router := mux.NewRouter()

	// API Endpoints
	router.HandleFunc("/api/files", fileHandler.UploadFile).Methods("POST")
	router.HandleFunc("/api/files", fileHandler.ListFiles).Methods("GET")
	router.HandleFunc("/api/files/{id}", fileHandler.GetFileByID).Methods("GET")

	// Adaugă middleware pentru CORS
	corsMiddleware := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{"*"}),
		gorillahandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Handler pentru servirea fișierelor statice pentru frontend
	staticHandler := http.FileServer(http.Dir("./frontend/dist"))
	router.PathPrefix("/").Handler(staticHandler)

	return corsMiddleware(router)
}
