package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func (cfg *apiConfig) middlewareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
		fmt.Println("hi")
	})
}

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	const rootPath = "."

	mux.HandleFunc("GET /api/healthz", hanlderReadiness)

	mux.Handle("POST /api/chirps", handlerPostChirps())
	mux.Handle("GET /api/chirps", handlerGetChirps())

	mux.Handle("GET /app/*", apiCfg.middlewareMetricsIncrement(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath)))))
	mux.Handle("GET /api/reset", apiCfg.handlerReset())
	mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())

	corsMux := middlewareCors(mux)
	s := &http.Server{
		Handler: corsMux,
		Addr:    "localhost:8080",
	}

	log.Fatal(s.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Cache-Control", "no-cache")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
