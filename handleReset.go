package main

import "net/http"

func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("test", "test")
		cfg.fileServerHits = 0
	})
}
