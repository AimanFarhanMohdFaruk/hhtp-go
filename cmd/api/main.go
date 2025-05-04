package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db database.Queries
	jwtSecret string
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	mux := http.NewServeMux()

	db, err := sql.Open("postgres", "user=aimanfarhan dbname=chirpy sslmode=disable")
	if err != nil {
		log.Fatal("DB Connection failed")
	}
	dbQueries := database.New(db)
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("Missing env variable: JWT_SECRET")
	}

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db: *dbQueries,
		jwtSecret: jwtSecret,
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

// func (cfg *apiConfig) authRequired(next http.HandlerFunc) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token, err := auth.GetBearerToken(r.Header)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return
// 		}

// 		_, err = auth.ValidateJWT(token, cfg.jwtSecret)

// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 			return
// 		}
// 		next.ServeHTTP(w,r)
// 	})
// }

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