package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

type loginReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type loginResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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

		accessToken := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer:   "chirpy-access",
				IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
				ExpiresAt: jwt.NewNumericDate(
					time.Now().UTC().Add(time.Hour * time.Duration(cfg.AccessTokenLifeHours)),
				),
				Subject: fmt.Sprintf("%d", id),
			},
		)
		signedAccessToken, err := accessToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Print("accessToken", err)
			utils.RespondWithError(w, 500, "error signing token")
			return
		}

		refreshToken := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer:   "chirpy-refresh",
				IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
				ExpiresAt: jwt.NewNumericDate(
					time.Now().UTC().Add(time.Duration(cfg.RefreshTokenLifeDays) * time.Hour * 24),
				),
				Subject: fmt.Sprintf("%d", id),
			},
		)

		signedRefreshToken, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Print("refreshToken", err)
			utils.RespondWithError(w, 500, "error signing token")
			return
		}
		utils.RespondWithJSON(w, 200, loginResp{Token: signedAccessToken, RefreshToken: signedRefreshToken})
	})
}
