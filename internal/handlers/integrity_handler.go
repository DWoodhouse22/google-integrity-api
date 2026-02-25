package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DWoodhouse22/google-integrity-api/internal/httputil"
	"github.com/DWoodhouse22/google-integrity-api/internal/models"
	"github.com/DWoodhouse22/google-integrity-api/internal/services"
)

type IntegrityHandler struct {
	nonceService *services.NonceService
}

func NewIntegrityHandler(nonceService *services.NonceService) *IntegrityHandler {
	return &IntegrityHandler{
		nonceService: nonceService,
	}
}

func (h *IntegrityHandler) GenerateNonce(w http.ResponseWriter, r *http.Request) {
	nonce, err := h.nonceService.Generate()
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to generate nonce")
		return
	}

	response := models.GenerateNonceResponse{
		Nonce: nonce,
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *IntegrityHandler) VerifyNonce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.VerifyNonceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Nonce == "" {
		httputil.WriteError(w, http.StatusBadRequest, "nonce required")
		return
	}

	valid := h.nonceService.Consume(req.Nonce)
	response := models.VerifyNonceResponse{
		Valid: valid,
	}

	statusCode := http.StatusOK
	if valid == false {
		statusCode = http.StatusUnauthorized
	}

	httputil.WriteJSON(w, statusCode, response)
}
