package handlers

import (
	"log"
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *ApiConfig) HandlerPostRevoke() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		idStr, err := parsedToken.Claims.GetSubject()
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			log.Println(err)
			return
		}
		if err = cfg.Db.AddRevokeToken(idStr, refreshToken); err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			log.Println(err)
			return
		}
		w.WriteHeader(200)
	})
}
