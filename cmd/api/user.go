package main

import (
	"encoding/json"
	"net/http"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/auth"
	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	"github.com/julienschmidt/httprouter"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
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

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := auth.AuthenticateUser(r, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	
	decoder := json.NewDecoder(r.Body)

	type requestParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := requestParams{}
	err = decoder.Decode(&params)
		
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userId,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
}