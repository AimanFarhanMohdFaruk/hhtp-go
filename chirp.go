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

	respondWithJSON(w, http.StatusCreated, database.Chirp{
		ID: chirp.ID,
		Body: chirp.Body,
		UserID: chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}