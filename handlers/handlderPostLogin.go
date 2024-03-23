package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

type loginReq struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}
type loginResp struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) HandlerPostLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqDat, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong processing your request")
			return
		}

		// parse the request data
		loginRequest := loginReq{}
		err = json.Unmarshal(reqDat, &loginRequest)
		if err != nil {
			utils.RespondWithError(w, 500, "something wernt wrong parsing your request")
			return
		}
		id, err := cfg.Db.GetUserIDByEmail(loginRequest.Email, loginRequest.Password)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, 401, "invalid email or password")
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:   "chirpy",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().UTC().Add(time.Second * time.Duration(loginRequest.ExpiresInSeconds)),
			),
			Subject: fmt.Sprintf("%d", id),
		})

		signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Print(err)
			utils.RespondWithError(w, 500, "error signing token")
		}
		utils.RespondWithJSON(
			w,
			200,
			loginResp{Token: signedToken},
		)
	})
}
