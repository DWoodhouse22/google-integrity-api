package models

type GenerateNonceResponse struct {
	Nonce string `json:"nonce"`
}

type VerifyNonceRequest struct {
	Nonce string `json:"nonce"`
}

type VerifyNonceResponse struct {
	Valid bool `json:"valid"`
}
