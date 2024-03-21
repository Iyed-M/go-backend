package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type postUserReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type postUserResp struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (cfg *apiConfig) handlerPostUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read Request Body
		resp, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error reading request body")
			return
		}
		r.Body.Close()

		// Parse Request Body
		parsedReq := postUserReq{}

		err = json.Unmarshal(resp, &parsedReq)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error parsing request body")
			return
		}

		// create user save it to db and retrun it
		hashedPassord, err := bcrypt.GenerateFromPassword([]byte(parsedReq.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondWithError(w, 500, "Something went wrong")
		}
		// User data stoted in db
		// the password is returned and hashed
		usr_, err := cfg.db.CreateUser(parsedReq.Email, string(hashedPassord))

		// POST RESPONSE
		usrResp := postUserResp{
			Email: usr_.Email,
			ID:    usr_.ID,
		}

		if err != nil && err.Error() != "empty file" {
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating user")
			return
		}

		utils.RespondWithJSON(w, 201, usrResp)
	})
}
