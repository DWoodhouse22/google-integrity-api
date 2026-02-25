package routes

import (
	"net/http"

	"github.com/DWoodhouse22/google-integrity-api/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.IntegrityHandler) {
	mux.HandleFunc("/nonce/generate", handler.GenerateNonce)
	mux.HandleFunc("/nonce/verify", handler.VerifyNonce)
}
