package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db database.Queries
}

func main() {
	mux := http.NewServeMux()

	db, err := sql.Open("postgres", "user=aimanfarhan dbname=chirpy sslmode=disable")
	if err != nil {
		log.Fatal("DB Connection failed")
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db: *dbQueries,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	server := http.Server{
		Addr: ":8080",
		Handler: apiCfg.routes(),	
	}
	log.Printf("Serving files from %s on port: %s\n", ".", ":8080")
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
			return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}