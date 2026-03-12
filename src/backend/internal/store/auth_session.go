package store

import (
	"context"
	"fmt"
	"time"

	"github.com/claway/server/internal/model"
)

// Auth sessions are stored in-memory (sync.Map on the Store struct) rather than
// in PostgreSQL because they are ephemeral — each session lives at most 5 minutes,
// is polled a handful of times, then discarded. A database table would add latency
// and require a background job to GC expired rows, with no durability benefit since
// sessions cannot survive a server restart anyway (the agent will just retry).

// CreateAuthSession stores a new pending auth session.
func (s *Store) CreateAuthSession(_ context.Context, session *model.AuthSession) error {
	s.authSessions.Store(session.ID, session)
	return nil
}

// GetAuthSession retrieves an auth session by ID, returning ErrNotFound if
// the session does not exist or has expired.
func (s *Store) GetAuthSession(_ context.Context, id string) (*model.AuthSession, error) {
	val, ok := s.authSessions.Load(id)
	if !ok {
		return nil, ErrNotFound
	}

	session := val.(*model.AuthSession)

	// Lazy expiry check — delete and return not-found if expired.
	if time.Now().After(session.ExpiresAt) {
		s.authSessions.Delete(id)
		return nil, ErrNotFound
	}

	return session, nil
}

// CompleteAuthSession marks a session as completed and stores the JWT token.
func (s *Store) CompleteAuthSession(_ context.Context, id string, token string) error {
	val, ok := s.authSessions.Load(id)
	if !ok {
		return ErrNotFound
	}

	session := val.(*model.AuthSession)
	if time.Now().After(session.ExpiresAt) {
		s.authSessions.Delete(id)
		return fmt.Errorf("auth session expired")
	}

	session.Token = token
	session.Status = "completed"
	// Re-store to ensure visibility (sync.Map Store is the publish barrier).
	s.authSessions.Store(id, session)
	return nil
}

// CleanupExpiredAuthSessions removes all expired sessions from the map.
// Called periodically from a background goroutine.
func (s *Store) CleanupExpiredAuthSessions() {
	now := time.Now()
	s.authSessions.Range(func(key, value any) bool {
		session := value.(*model.AuthSession)
		if now.After(session.ExpiresAt) {
			s.authSessions.Delete(key)
		}
		return true
	})
}
