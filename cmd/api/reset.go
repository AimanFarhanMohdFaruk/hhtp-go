package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	users, err := cfg.db.ListUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error fetching user lists")
		return
	}	
	chirps, err := cfg.db.ListChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error fetching user lists")
		return
	}
	for _, user := range users {
		cfg.db.DeleteUser(r.Context(), user.ID)
	}
	for _, chirp := range chirps {
		cfg.db.DeleteChirp(r.Context(), chirp.ID)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
