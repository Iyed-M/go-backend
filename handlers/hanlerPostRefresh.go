package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

type postRefreshResp struct {
	AccessToken string `json:"token"`
}

func (cfg *ApiConfig) HandlerPostRefresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		refreshToken := r.Header.Get("Authorization")[len("Bearer "):]

		parsedToken, err := jwt.ParseWithClaims(
			refreshToken,
			&jwt.RegisteredClaims{},
			jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			}),
		)
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			log.Println(err)
			return
		}

		issuer, err := parsedToken.Claims.GetIssuer()
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			return
		}

		if issuer != "chirpy-refresh" {
			utils.RespondWithError(w, 401, "invalid refresh token")
			return
		}

		idStr, err := parsedToken.Claims.GetSubject()
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			return
		}

		if cfg.Db.IsTokenRevoked(idStr, refreshToken) {
			utils.RespondWithError(w, 401, "revoked refresh tokenk")
			return
		}

		id, _ := strconv.Atoi(idStr)

		signedAccessToken, err := newSignedToken("chirpy-access", cfg.AccessTokenLifeHours, id, cfg.JWTSecret)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
			return
		}
		utils.RespondWithJSON(w, 200, postRefreshResp{AccessToken: signedAccessToken})
	})
}
