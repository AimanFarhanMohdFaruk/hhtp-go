package main

import (
	"encoding/json"
	"net/http"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/auth"
	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type webhookData struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	polkaAPIKey, err  := auth.GetPolkaAPIKey(r.Header)
	println(polkaAPIKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if polkaAPIKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Event string `json:"event"`
		Data webhookData `json:"data"`
	}
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, map[string]interface{}{} )
		return
	}

	err = cfg.db.UpdateUserChirpyRed(r.Context(), database.UpdateUserChirpyRedParams{
		ID: params.Data.UserID,
		IsChirpyRed: true,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]interface{}{} )
}