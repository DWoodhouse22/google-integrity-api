package services

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/DWoodhouse22/google-integrity-api/internal/repositories"
)

type TokenService struct {
	repo *repositories.TokenRepository
}

func NewTokenService(repo *repositories.TokenRepository) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) Generate() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(b)
	s.repo.Store(token)
	return token, nil
}

func (s *TokenService) Consume(token string) bool {
	return s.repo.Consume(token)
}
