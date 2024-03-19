package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(response)
	return nil
}

func RespondWithError(w http.ResponseWriter, status int, message string) error {
	return RespondWithJSON(w, status, map[string]string{"error": message})
}
