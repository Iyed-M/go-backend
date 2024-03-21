package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
)

type postUserResp struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerPostUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error reading request body")
			return
		}
		r.Body.Close()

		parsedEmail := postUserResp{}
		err = json.Unmarshal(resp, &parsedEmail)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error parsing request body")
			return
		}

		usr, err := cfg.db.CreateUser(parsedEmail.Email)
		if err != nil && err.Error() != "empty file" {
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating user")
			return
		}

		utils.RespondWithJSON(w, 201, usr)
	})
}
