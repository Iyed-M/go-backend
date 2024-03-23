package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Iyed-M/go-backend/database"
	"github.com/Iyed-M/go-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type putUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *ApiConfig) HandlerPutUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]

		// parse token
		parsedToken, err := jwt.ParseWithClaims(
			token,
			&jwt.RegisteredClaims{},
			jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			}),
		)
		if err != nil {
			utils.RespondWithError(w, 401, "invalid token")
			log.Println(err)
			return
		}
		strID, err := parsedToken.Claims.GetSubject()
		if err != nil {
			utils.RespondWithError(w, 401, "invalid token")
			log.Println(err)
			return
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			utils.RespondWithError(w, 400, "somthing went wrong")
			log.Println(err)
			return
		}
		paresedRequestBody := &putUserReq{}
		err = json.NewDecoder(r.Body).Decode(paresedRequestBody)
		if err != nil {
			utils.RespondWithError(w, 400, "somthing went wrong")
			log.Println(err)
			return
		}

		hashedNewPass, err := bcrypt.GenerateFromPassword(
			[]byte(paresedRequestBody.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			utils.RespondWithError(w, 400, "somthing went wrong")
			log.Println(err)
			return
		}
		log.Printf("new User :%+v\n", *paresedRequestBody)

		// update the user data in the database
		newUserData := database.User{
			ID:       id,
			Email:    paresedRequestBody.Email,
			Password: string(hashedNewPass),
		}
		cfg.Db.UpdateUser(newUserData)

		// Respond with the new user data
		utils.RespondWithJSON(
			w,
			200,
			map[string]any{"email": newUserData.Email, "id": newUserData.ID},
		)
	})
}
