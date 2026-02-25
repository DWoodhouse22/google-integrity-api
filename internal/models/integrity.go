package models

type GenerateNonceResponse struct {
	Nonce string `json:"token"`
}

type VerifyNonceRequest struct {
	Nonce string `json:"token"`
}

type VerifyNonceResponse struct {
	Valid bool `json:"valid"`
}
