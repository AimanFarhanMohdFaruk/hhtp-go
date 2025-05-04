package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/auth"
	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}
	decoder := json.NewDecoder(r.Body)

	type requestParams struct {
		Body string `json:"body"`
	}

	params := requestParams{}
	err = decoder.Decode(&params)
		
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirpList, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Error fetching chirps")
	}
	respondWithJSON(w, 200, chirpList)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
			respondWithError(w, http.StatusBadRequest, "Misisng Chirp ID")
			return
	}
	
	parseUUID, err := uuid.Parse(id)
	if err !=  nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parseUUID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	respondWithJSON(w, 200, chirp)
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