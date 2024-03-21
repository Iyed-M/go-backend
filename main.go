package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Iyed-M/go-backend/database"
	"github.com/Iyed-M/go-backend/handlers"
)

var DataBasePath = "database.json"

func main() {
	apiCfg := &handlers.ApiConfig{
		Db: database.NewDB(DataBasePath),
	}
	mux := http.NewServeMux()
	const rootPath = "."

	dbg := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *dbg {
		deleteDBafterTest(DataBasePath)
	}

	mux.HandleFunc("GET /api/healthz", hanlderReadiness)

	mux.Handle("POST /api/login", apiCfg.HandlerPostLogin())

	mux.Handle("POST /api/users", apiCfg.HandlerPostUsers())

	mux.Handle("POST /api/chirps", apiCfg.HandlerPostChirps())
	mux.Handle("GET /api/chirps", apiCfg.HandlerGetChirps())
	mux.Handle("GET /api/chirps/{id}", apiCfg.HandlerGetChirpByID())

	// mux.Handle("GET /app/*", apiCfg.middlewareMetricsIncrement(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath)))))
	// mux.Handle("GET /api/reset", apiCfg.handlerReset())
	// mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())

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

func deleteDBafterTest(path string) {
	os.Remove(path)
}
