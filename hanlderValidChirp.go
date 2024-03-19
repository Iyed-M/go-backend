package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
)

type chirpRequest struct {
	Body string `json:"body"`
}

func handlerValidChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			if err = utils.RespondWithError(w, 400, "Something went wrong"); err != nil {
				utils.RespondWithError(w, 500, "Something went wrong")
			}
		}

		chirp := &chirpRequest{}
		err = json.Unmarshal(reqBody, chirp)
		if err != nil {
			if err = utils.RespondWithError(w, 400, "Something went wrong"); err != nil {
				utils.RespondWithError(w, 500, "Something went wrong")
			}
			return
		}

		if len(chirp.Body) > 140 {
			if err = utils.RespondWithError(w, 400, "Chirp is too long"); err != nil {
				utils.RespondWithError(w, 500, "Something went wrong")
			}
			return
		}

		validChirp := utils.CleanChirp(chirp.Body)
		utils.RespondWithJSON(w, 200, map[string]string{"cleaned_body": validChirp})
	})
}
