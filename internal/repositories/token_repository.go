package repositories

import (
	"sync"
	"time"
)

const tokenTTL = 3 * time.Minute

type TokenRepository struct {
	mu     sync.RWMutex
	tokens map[string]time.Time
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{
		tokens: make(map[string]time.Time),
	}
}

func (r *TokenRepository) Store(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token] = time.Now().Add(tokenTTL)
}

func (r *TokenRepository) Consume(token string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	expiry, ok := r.tokens[token]
	if !ok || time.Now().After(expiry) {
		delete(r.tokens, token)
		return false
	}

	delete(r.tokens, token)
	return true
}
