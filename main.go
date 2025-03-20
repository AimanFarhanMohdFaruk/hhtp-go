package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	// user handlers
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	// chirp handlers
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpHandler)

	server := http.Server{
		Addr: ":8080",
		Handler: mux,	
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "server error")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Bad request")
		return
	}

	stringSplit :=  strings.Fields(params.Body)
	wordsToCensor := []string{"kerfuffle", "sharbert", "fornax" }
	censorMap :=	 make(map[string]bool)
	for _, word := range wordsToCensor {
		censorMap[word] = true
	}
	for i, word := range stringSplit {
		if censorMap[strings.ToLower(word)] {
			stringSplit[i] = "****"
		}
	}

	type responseVal struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	respBody := responseVal{
		Cleaned_body: strings.Join(stringSplit, " "),
	}

	respondWithJSON(w, 200, respBody)
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