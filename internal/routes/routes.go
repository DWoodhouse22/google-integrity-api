package routes

import (
	"net/http"

	"github.com/DWoodhouse22/google-integrity-api/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.IntegrityHandler) {
	mux.HandleFunc("/token/generate", handler.GenerateNonce)
	mux.HandleFunc("/token/verify", handler.VerifyNonce)
}
