package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Iyed-M/go-backend/database"
	"github.com/Iyed-M/go-backend/utils"
)

func (cfg *ApiConfig) HandlerPostChirps() http.Handler {
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
		chirp, err := cfg.Db.CreateChirp(chirpSent.Body)
		if err != nil {
			utils.RespondWithError(w, 500, "error creating chirp")
		}

		err = utils.RespondWithJSON(w, 201, chirp)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
		}
	})
}

func (cfg *ApiConfig) HandlerGetChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cfg.Db.GetChirps()
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		utils.RespondWithJSON(w, 200, chirps)
	})
}
