package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

type postChirpReq struct {
	Body string `json:"body"`
}

func (cfg *ApiConfig) HandlerPostChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")[len("Bearer "):]
		parsedToken, err := jwt.ParseWithClaims(
			token,
			&jwt.RegisteredClaims{},
			jwt.Keyfunc(func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			}),
		)
		if err != nil {
			utils.RespondWithError(w, 400, "invalid token")
			log.Println(5)
			log.Println(err)
			return
		}
		idStr, err := parsedToken.Claims.GetSubject()
		if err != nil {
			utils.RespondWithError(w, 400, "invalid token")
		}

		// vaidate shirp
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
			log.Println(4)
			return
		}

		chirpSent := &postChirpReq{}
		err = json.Unmarshal(reqBody, &chirpSent)
		if err != nil {
			utils.RespondWithError(w, 500, err.Error())
			log.Println(3)
			return
		}

		if len(chirpSent.Body) > 140 {
			utils.RespondWithError(w, 400, "Chirp is too long")
			log.Println(2)
			return
		}

		// save chrip in db
		// chirp, err := cfg.Db.CreateChirp(chirpSent.Body)
		chirp, err := cfg.Db.CreateChirp(idStr, chirpSent.Body)
		if err != nil {
			log.Println(1)
			utils.RespondWithError(w, 500, "error creating chirp")
			return
		}

		utils.RespondWithJSON(w, 201, chirp)
	})
}
