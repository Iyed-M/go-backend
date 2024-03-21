package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Iyed-M/go-backend/database"
	"github.com/Iyed-M/go-backend/utils"
)

func incrementID() func() int {
	id := 1
	return func() int {
		id++
		return id
	}
}

func (cfg *apiConfig) handlerPostChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// vaidate shirp
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
			return
		}

		chirpSent := database.Chirp{}
		err = json.Unmarshal(reqBody, &chirpSent)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
			return
		}

		if len(chirpSent.Body) > 140 {
			utils.RespondWithError(w, 400, "Chirp is too long")
			return
		}

		// save chrip in db
		chirp := cfg.db.CreateChirp(chirpSent.Body)

		err = utils.RespondWithJSON(w, 200, chirp)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
		}
	})
}

func (cfg *apiConfig) handlerGetChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cfg.db.GetChirps()
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		utils.RespondWithJSON(w, 200, chirps)
	})
}
