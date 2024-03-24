package handlers

import (
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
)

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
