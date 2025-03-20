package main

import (
	"encoding/json"
	"net/http"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type requestParams struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := requestParams{}
	err := decoder.Decode(&params)
		
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
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