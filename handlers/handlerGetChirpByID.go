package handlers

import (
	"net/http"
	"strconv"

	"github.com/Iyed-M/go-backend/utils"
)

func (cfg *ApiConfig) HandlerGetChirpByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get id from path
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
		}

		chirps, err := cfg.Db.GetChirps()
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
		}

		for _, chirp := range chirps {
			if chirp.ID == id {
				utils.RespondWithJSON(w, 200, chirp)
				return
			}
		}
		// case ID doesnt exit
		utils.RespondWithError(w, 404, "Chirp Not Found")
	})
}
