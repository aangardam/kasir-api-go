package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseSuccess(w http.ResponseWriter, payload interface{}, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "OK",
		"message": message,
		"data":    payload,
	})
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ERROR",
		"message": message,
	})
}
