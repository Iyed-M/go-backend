package handlers

import "net/http"

func (cfg *ApiConfig) HandlerPostLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
