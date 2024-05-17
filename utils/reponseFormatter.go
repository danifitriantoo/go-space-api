package utils

import (
	"encoding/json"
	"go-space-api/models"
	"net/http"
)

func RespondSuccess(w http.ResponseWriter, status int, data interface{}) {
	response := models.Response{Status: "success", Data: data}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	response := models.Response{Status: "error", Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
