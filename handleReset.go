package main

import "net/http"

func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits = 0
	})
}
