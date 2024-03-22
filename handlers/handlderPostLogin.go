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
	ID    int    `json:"id"`
	Email string `json:"email"`
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
		utils.RespondWithJSON(w, 200, loginResp{ID: id, Email: loginRequest.Email})
	})
}
