package services

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/DWoodhouse22/google-integrity-api/internal/repositories"
)

type NonceService struct {
	repo *repositories.NonceRepository
}

func NewNonceService(repo *repositories.NonceRepository) *NonceService {
	return &NonceService{repo: repo}
}

func (s *NonceService) Generate() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	nonce := base64.StdEncoding.EncodeToString(b)
	s.repo.Store(nonce)
	return nonce, nil
}

func (s *NonceService) Consume(nonce string) bool {
	return s.repo.Consume(nonce)
}
