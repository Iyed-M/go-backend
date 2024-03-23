package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Iyed-M/go-backend/utils"
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

		signedAccessToken, err := newSignedToken("chirpy-access", cfg.AccessTokenLifeHours, id, cfg.JWTSecret)
		if err != nil {
			log.Print("accessToken", err)
			utils.RespondWithError(w, 500, "error signing token")
			return
		}

		signedRefreshToken, err := newSignedToken("chirpy-refresh", cfg.RefreshTokenLifeDays*24, id, cfg.JWTSecret)
		if err != nil {
			log.Print("refreshToken", err)
			utils.RespondWithError(w, 500, "error signing token")
			return
		}
		utils.RespondWithJSON(w, 200, loginResp{Token: signedAccessToken, RefreshToken: signedRefreshToken})
	})
}
