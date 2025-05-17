package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/auth"
	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := auth.AuthenticateUser(r, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
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

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	currentUserID, err := auth.AuthenticateUser(r, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	id := ps.ByName("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Chirp ID")
		return
	}
	
	parseUUID, err := uuid.Parse(id)
	if err !=  nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parseUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if chirp.UserID != currentUserID {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	respondWithJSON(w, http.StatusNoContent, chirp.ID)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryParams := r.URL.Query()
	authorID, err := uuid.Parse(queryParams.Get("author_id"))
	sortParams := queryParams.Get("sort")
	
	chirpList, err := cfg.db.ListChirps(r.Context(), uuid.NullUUID{
		UUID: authorID,
		Valid: err == nil,
	})
		
	if err != nil {
		respondWithError(w, 400, "Error fetching chirps")
	}

	if sortParams != "" {
		if sortParams == "asc" {
			sort.Slice(chirpList, func(i, j int) bool { return chirpList[i].CreatedAt.Time.Before(chirpList[j].CreatedAt.Time) })
		}
		if sortParams == "desc" {
			sort.Slice(chirpList, func(i, j int) bool { return chirpList[i].CreatedAt.Time.After(chirpList[j].CreatedAt.Time) })
		}
	}
	
	respondWithJSON(w, 200, chirpList)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Chirp ID")
		return
	}
	
	parseUUID, err := uuid.Parse(id)
	if err !=  nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parseUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	respondWithJSON(w, 200, chirp)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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