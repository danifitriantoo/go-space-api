package handlers

import (
	"encoding/json"
	"go-space-api/models"
	"go-space-api/storage"
	"go-space-api/utils"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPlanets(w http.ResponseWriter, r *http.Request) {
	planets, err := storage.GetAllPlanets()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch planets")
		return
	}
	utils.RespondSuccess(w, http.StatusOK, planets)
}

func GetPlanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	planet, found := storage.GetPlanetById(params["id"])
	if !found {
		utils.RespondWithError(w, http.StatusNotFound, "Planet not found")
		return
	}
	utils.RespondSuccess(w, http.StatusOK, planet)
}

func CreatePlanet(w http.ResponseWriter, r *http.Request) {
	var planet models.Planet
	if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := storage.AddPlanet(planet); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create planet")
		return
	}
	utils.RespondSuccess(w, http.StatusCreated, planet)
}

func UpdatePlanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedPlanet models.Planet
	if err := json.NewDecoder(r.Body).Decode(&updatedPlanet); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid planet ID")
		return
	}
	updatedPlanet.ID = objID

	if err := storage.UpdatePlanet(updatedPlanet); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Planet not found")
		return
	}
	utils.RespondSuccess(w, http.StatusOK, updatedPlanet)
}

func DeletePlanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := storage.DeletePlanet(params["id"]); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Planet not found")
		return
	}
	utils.RespondSuccess(w, http.StatusNoContent, nil)
}
