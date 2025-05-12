package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (cfg *apiConfig) routes() *httprouter.Router{
	router := httprouter.New()
	
	router.HandlerFunc(http.MethodGet, "/api/healthz", cfg.logHandler(readinessHandler))
	router.GET("/admin/metrics", cfg.metricsHandler)
	
	router.POST("/admin/reset", cfg.resetHandler)
	router.POST("/api/validate_chirp", validateChirpHandler)
	router.POST("/api/users", cfg.createUserHandler)
	router.PUT("/api/users", cfg.updateUserHandler)

	router.POST("/api/login", cfg.loginHandler)
	router.POST("/api/refresh", cfg.refreshTokenHandler)
	router.POST( "/api/revoke", cfg.revokeRefreshTokenHandler)

	router.POST("/api/chirps", cfg.createChirpHandler)
	router.GET("/api/chirps", cfg.getChirpsHandler)
	router.GET("/api/chirps/:id", cfg.getChirpHandler)
	router.DELETE("/api/chirps/:id", cfg.deleteChirpHandler)

	router.POST("/api/polka/webhooks", cfg.polkaWebhookHandler)

	return router
}
