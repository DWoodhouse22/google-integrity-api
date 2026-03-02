package main

import (
	"log"
	"net/http"

	"github.com/DWoodhouse22/google-integrity-api/internal/handlers"
	"github.com/DWoodhouse22/google-integrity-api/internal/repositories"
	"github.com/DWoodhouse22/google-integrity-api/internal/routes"
	"github.com/DWoodhouse22/google-integrity-api/internal/services"
)

func main() {
	tokenRepo := repositories.NewTokenRepository()
	tokenService := services.NewTokenService(tokenRepo)
	integrityService, err := services.NewIntegrityService()
	if err != nil {
		panic(err)
	}
	handler := handlers.NewIntegrityHandler(tokenService, integrityService)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, handler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
