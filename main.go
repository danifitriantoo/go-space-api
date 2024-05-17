package main

import (
	"go-space-api/db"
	"go-space-api/router"
	"go-space-api/storage"
	"log"
	"net/http"

	_ "go-space-api/docs" // swagger docs

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Solar System API
// @version 1.0
// @description This is an API for managing planets in the solar system.
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url http://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
func main() {

	db.Init()
	// Initialize storage with the MongoDB client
	storage.Init(db.Client)
	r := router.AppRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
