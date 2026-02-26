package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DWoodhouse22/google-integrity-api/internal/api/response"
	"github.com/DWoodhouse22/google-integrity-api/internal/models"
	"github.com/DWoodhouse22/google-integrity-api/internal/services"
)

type IntegrityHandler struct {
	nonceService     *services.NonceService
	integrityService *services.IntegrityService
}

func NewIntegrityHandler(nonceService *services.NonceService, integrityService *services.IntegrityService) *IntegrityHandler {
	return &IntegrityHandler{
		nonceService:     nonceService,
		integrityService: integrityService,
	}
}

func (h *IntegrityHandler) GenerateNonce(w http.ResponseWriter, r *http.Request) {
	nonce, err := h.nonceService.Generate()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate one-time token")
		return
	}

	response.WriteSuccess(w, http.StatusOK, models.GenerateNonceResponse{
		Nonce: nonce,
	})
}

func (h *IntegrityHandler) VerifyNonce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.VerifyNonceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Nonce == "" {
		response.WriteError(w, http.StatusBadRequest, "one-time token required")
		return
	}

	valid := h.nonceService.Consume(req.Nonce)
	response.WriteSuccess(w, http.StatusOK, models.VerifyNonceResponse{
		Valid: valid,
	})
}

func (h *IntegrityHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.VerifyIntegrityTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		response.WriteError(w, http.StatusBadRequest, "token required")
		return
	}

	decodedToken, err := h.integrityService.DecodeToken(req.Token)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "failed to decode integrity token")
		return
	}

	requestHash := decodedToken.TokenPayloadExternal.RequestDetails.RequestHash
	if !h.nonceService.Consume(requestHash) {
		response.WriteSuccess(w, http.StatusOK, models.VerifyIntegrityTokenResponse{
			Verdict: "invalid",
			Reason:  "invalid request hash",
		})
		return
	}

	if ok, reason := h.integrityService.ValidateToken(decodedToken); !ok {
		response.WriteSuccess(w, http.StatusOK, models.VerifyIntegrityTokenResponse{
			Verdict: "invalid",
			Reason:  reason,
		})
		return
	}

	response.WriteSuccess(w, http.StatusOK, models.VerifyIntegrityTokenResponse{
		Verdict: "valid",
	})
}
