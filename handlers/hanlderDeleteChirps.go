package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Iyed-M/go-backend/database"
	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *ApiConfig) HanlderDeleteChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirpIDStr := r.PathValue("chirpID")
		token := r.Header.Get("Authorization")[len("Bearer "):]

		// parse token
		parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, jwt.Keyfunc(func(t *jwt.Token) (interface{}, error) {
			if issuer, _ := t.Claims.GetIssuer(); issuer != "chirpy-access" {
				log.Print(issuer)
				return nil, errors.New("NIIIICEk")
			}
			return []byte(cfg.JWTSecret), nil
		}))
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, 500, "something went wrong")
		}

		// get chirp by id
		chirp, err := cfg.Db.GetChirpByID(chirpIDStr)
		if err == database.ErrChirpNotFound {
			utils.RespondWithError(w, 400, "chirp not found")
			log.Println(err)
			return
		}
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			log.Println(err)
			return
		}

		authorID, err := parsedToken.Claims.GetSubject()
		if err != nil {
			utils.RespondWithError(w, 500, "something went wrong")
			log.Println(err)
			return
		}
		if chirp.AuthorId != authorID {
			utils.RespondWithError(w, 403, "unauthorized access")
			return
		}
		cfg.Db.DeleteChirp(chirpIDStr)
		w.WriteHeader(200)
	})
}
