package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"go-space-api/db"
	"go-space-api/models"
	"go-space-api/router"
	"go-space-api/storage"

	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func setupRouter() *mux.Router {
	db.Init()
	storage.Init(db.Client)
	cleanUpDB()
	return router.AppRouter()
}

func cleanUpDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := storage.GetPlanetCollection().Drop(ctx)
	if err != nil {
		log.Fatalf("Failed to clean up database: %v", err)
	}
}

func TestGetPlanets(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/planets", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.Response
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.Empty(t, response.Data)

	t.Logf("TestGetPlanets passed")
}

func TestGetPlanet_NotFound(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/planet/nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response models.Response
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Planet not found", response.Message)

	t.Logf("TestGetPlanet_NotFound passed")
}

func TestCreatePlanet(t *testing.T) {
	r := setupRouter()
	planet := models.Planet{
		Name:     "Earth",
		Diameter: 12742,
		Moons:    1,
		Distance: 149,
	}
	planetJSON, _ := json.Marshal(planet)
	req, _ := http.NewRequest("POST", "/planet", bytes.NewBuffer(planetJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.Response
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)

	var createdPlanet models.Planet
	data, _ := json.Marshal(response.Data)
	json.Unmarshal(data, &createdPlanet)
	assert.NotEmpty(t, createdPlanet.ID)
	assert.Equal(t, planet.Name, createdPlanet.Name)
	assert.Equal(t, planet.Diameter, createdPlanet.Diameter)
	assert.Equal(t, planet.Moons, createdPlanet.Moons)
	assert.Equal(t, planet.Distance, createdPlanet.Distance)

	t.Logf("TestCreatePlanet passed")
}

func TestUpdatePlanet(t *testing.T) {
	r := setupRouter()
	initialPlanet := models.Planet{
		Name:     "Earth",
		Diameter: 12742,
		Moons:    1,
		Distance: 149,
	}
	storage.AddPlanet(initialPlanet)

	var createdPlanet models.Planet
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	storage.GetPlanetCollection().FindOne(ctx, bson.M{"name": initialPlanet.Name}).Decode(&createdPlanet)

	updatedPlanet := models.Planet{
		ID:       createdPlanet.ID,
		Name:     "Earth Updated",
		Diameter: 12742,
		Moons:    1,
		Distance: 149,
	}
	updatedPlanetJSON, _ := json.Marshal(updatedPlanet)
	req, _ := http.NewRequest("PUT", "/planet/"+createdPlanet.ID.Hex(), bytes.NewBuffer(updatedPlanetJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.Response
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)

	var updatedPlanetResponse models.Planet
	data, _ := json.Marshal(response.Data)
	json.Unmarshal(data, &updatedPlanetResponse)
	assert.Equal(t, updatedPlanet.Name, updatedPlanetResponse.Name)
	assert.Equal(t, updatedPlanet.Diameter, updatedPlanetResponse.Diameter)
	assert.Equal(t, updatedPlanet.Moons, updatedPlanetResponse.Moons)
	assert.Equal(t, updatedPlanet.Distance, updatedPlanetResponse.Distance)

	t.Logf("TestUpdatePlanet passed")
}

func TestDeletePlanet(t *testing.T) {
	r := setupRouter()
	planet := models.Planet{
		Name:     "Earth",
		Diameter: 12742,
		Moons:    1,
		Distance: 149,
	}
	storage.AddPlanet(planet)

	var createdPlanet models.Planet
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	storage.GetPlanetCollection().FindOne(ctx, bson.M{"name": planet.Name}).Decode(&createdPlanet)

	req, _ := http.NewRequest("DELETE", "/planet/"+createdPlanet.ID.Hex(), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	req, _ = http.NewRequest("GET", "/planet/"+createdPlanet.ID.Hex(), nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response models.Response
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Planet not found", response.Message)

	t.Logf("TestDeletePlanet passed")
}
