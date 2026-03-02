package models

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

type VerifyTokenResponse struct {
	Valid bool `json:"valid"`
}

type VerifyIntegrityTokenRequest struct {
	Token string `json:"token"`
}

type VerifyIntegrityTokenResponse struct {
	Verdict string `json:"verdict"`
	Reason  string `json:"reason,omitempty"`
}
