package api

import (
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Rute pentru mesaje
	router.HandleFunc("/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/messages", GetMessages).Methods("GET")

	return router
}
