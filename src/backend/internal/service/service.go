package service

import (
	"github.com/claway/server/internal/config"
	"github.com/claway/server/internal/store"
)

// Service contains business logic and delegates to store for data access.
type Service struct {
	store *store.Store
	cfg   *config.Config
}

// New creates a new Service instance.
func New(s *store.Store, cfg *config.Config) *Service {
	return &Service{store: s, cfg: cfg}
}
