package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type chirpRequest struct {
	Body string `json:"body"`
}

type errorResp struct {
	Error string `json:"error"`
}

func handlerValidChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirp := &chirpRequest{}
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal(reqBody, chirp)
		if err != nil {
			log.Print(err)
		}

		if len(chrip.Body) > 140 {
			invalidResp := &errorResp{""}
		}
	})
}

func RespondWihtJSON(w http.ResponseWriter, r *http.Request) {

}
