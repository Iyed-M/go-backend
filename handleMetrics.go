package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits)
	})
}
