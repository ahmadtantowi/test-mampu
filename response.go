package main

import (
	"encoding/json"
	"net/http"
)

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendJSONError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"message": message})
}
