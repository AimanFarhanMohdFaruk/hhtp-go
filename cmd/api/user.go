package main

import (
	"encoding/json"
	"net/http"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email string `json:"email"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
		
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, database.User{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}