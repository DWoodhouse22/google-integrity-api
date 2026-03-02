package routes

import (
	"net/http"

	"github.com/DWoodhouse22/google-integrity-api/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.IntegrityHandler) {
	mux.HandleFunc("/token/generate", handler.GenerateToken)
	mux.HandleFunc("/integrity/verify-token", handler.VerifyIntegrityToken)
}
