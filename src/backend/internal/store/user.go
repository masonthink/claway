package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

// GetUserByOpenClawID retrieves a user by their OpenClaw ID.
func (s *Store) GetUserByOpenClawID(ctx context.Context, openclawID string) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx,
		`SELECT id, openclaw_id, username, agent_api_key, credits_balance, created_at, updated_at
		 FROM users WHERE openclaw_id = $1`, openclawID,
	).Scan(&u.ID, &u.OpenClawID, &u.Username, &u.AgentAPIKey, &u.CreditsBalance, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by openclaw_id: %w", err)
	}
	return &u, nil
}

// GetUserByID retrieves a user by their internal ID.
func (s *Store) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx,
		`SELECT id, openclaw_id, username, agent_api_key, credits_balance, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.OpenClawID, &u.Username, &u.AgentAPIKey, &u.CreditsBalance, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

// CreateUser inserts a new user and returns the created record.
func (s *Store) CreateUser(ctx context.Context, openclawID, username string) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx,
		`INSERT INTO users (openclaw_id, username)
		 VALUES ($1, $2)
		 RETURNING id, openclaw_id, username, agent_api_key, credits_balance, created_at, updated_at`,
		openclawID, username,
	).Scan(&u.ID, &u.OpenClawID, &u.Username, &u.AgentAPIKey, &u.CreditsBalance, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

// UpdateCreditsBalance atomically adjusts a user's credits balance by delta.
// Delta can be positive (credit) or negative (debit).
func (s *Store) UpdateCreditsBalance(ctx context.Context, userID int64, delta float64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE users SET credits_balance = credits_balance + $1, updated_at = NOW()
		 WHERE id = $2`, delta, userID,
	)
	if err != nil {
		return fmt.Errorf("update credits balance: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
