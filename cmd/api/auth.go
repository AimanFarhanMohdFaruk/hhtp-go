package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "invalid request")
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "invalid request")
		return
	}

	type loginResponse struct{
		ID uuid.UUID
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		Email string
		Token string
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: jwtToken,
	})
}