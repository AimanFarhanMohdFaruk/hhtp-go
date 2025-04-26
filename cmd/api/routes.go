package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (cfg *apiConfig) routes() *httprouter.Router{
	router := httprouter.New()
	
	router.HandlerFunc(http.MethodGet, "/api/healthz", readinessHandler)
	router.HandlerFunc(http.MethodGet,"/admin/metrics", cfg.metricsHandler)
	
	router.HandlerFunc(http.MethodPost,"/admin/reset", cfg.resetHandler)
	router.HandlerFunc(http.MethodPost,"/api/validate_chirp", validateChirpHandler)
	router.HandlerFunc(http.MethodPost,"/api/users", cfg.createUserHandler)

	router.HandlerFunc(http.MethodPost,"/api/login", cfg.loginHandler)

	router.HandlerFunc(http.MethodPost,"/api/chirps", cfg.createChirpHandler)
	router.HandlerFunc(http.MethodGet,"/api/chirps", cfg.getChirpsHandler)
	router.HandlerFunc(http.MethodGet,"/api/chirps/:id", cfg.getChirpHandler)

	return router
}

