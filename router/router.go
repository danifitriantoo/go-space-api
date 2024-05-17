package router

import (
	"go-space-api/handlers"

	"github.com/gorilla/mux"
)

func AppRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/planets", handlers.GetPlanets).Methods("GET")
	router.HandleFunc("/planets/{id}", handlers.GetPlanet).Methods("GET")
	router.HandleFunc("/planets", handlers.CreatePlanet).Methods("POST")
	router.HandleFunc("/planets/{id}", handlers.UpdatePlanet).Methods("PUT")
	router.HandleFunc("/planets/{id}", handlers.DeletePlanet).Methods("DELETE")

	return router
}
