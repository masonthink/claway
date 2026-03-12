package store

import (
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Common sentinel errors for store operations.
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

// ComputeEntry represents aggregated compute usage per user.
type ComputeEntry struct {
	UserID         int64   `json:"user_id"`
	Username       string  `json:"username"`
	TotalCost      float64 `json:"total_cost_usd"`
	TotalTokensIn  int     `json:"total_tokens_in"`
	TotalTokensOut int     `json:"total_tokens_out"`
	CallCount      int     `json:"call_count"`
}

// Store wraps database access.
type Store struct {
	db *pgxpool.Pool

	// authSessions holds in-memory auth sessions for agent-based login flows.
	// Using sync.Map because sessions are short-lived (5 min) and don't need persistence.
	authSessions sync.Map
}

// New creates a new Store instance.
func New(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}
