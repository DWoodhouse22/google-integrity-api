package repositories

import (
	"sync"
	"time"
)

const nonceTTL = 3 * time.Minute

type NonceRepository struct {
	mu     sync.RWMutex
	nonces map[string]time.Time
}

func NewNonceRepository() *NonceRepository {
	return &NonceRepository{
		nonces: make(map[string]time.Time),
	}
}

func (r *NonceRepository) Store(nonce string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nonces[nonce] = time.Now().Add(nonceTTL)
}

func (r *NonceRepository) Consume(nonce string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	expiry, ok := r.nonces[nonce]
	if !ok || time.Now().After(expiry) {
		delete(r.nonces, nonce)
		return false
	}

	delete(r.nonces, nonce)
	return true
}
