package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type chirpRequest struct {
	Body string `json:"body"`
}

func handlerValidChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirp := &chirpRequest{}
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			if err = RespondWithError(w, 400, "Something went wrong"); err != nil {
				RespondWithError(w, 500, "Something went wrong")
			}
		}

		err = json.Unmarshal(reqBody, chirp)
		if err != nil {
			if err = RespondWithError(w, 400, "Something went wrong"); err != nil {
				RespondWithError(w, 500, "Something went wrong")
			}
			return
		}

		if len(chirp.Body) > 140 {
			if err = RespondWithError(w, 400, "Chirp is too long"); err != nil {
				RespondWithError(w, 500, "Something went wrong")
			}
			return
		}

		if err = RespondWithJSON(w, 200, map[string]bool{"valid": true}); err != nil {
			RespondWithError(w, 500, "Something went wrong")
		}
	})
}

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
