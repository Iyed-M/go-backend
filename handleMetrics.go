package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Hits : %d\n", cfg.fileServerHits)
		htmlData := fmt.Sprintf(" <html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html> ", cfg.fileServerHits)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlData))
	})
}
