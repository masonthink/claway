package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

// CreatePRD inserts a new PRD for an idea.
func (s *Store) CreatePRD(ctx context.Context, ideaID int64, content string, priceCredits float64) (*model.PRD, error) {
	var p model.PRD
	err := s.db.QueryRow(ctx,
		`INSERT INTO prds (idea_id, content, price_credits, published_at)
		 VALUES ($1, $2, $3, NOW())
		 RETURNING id, idea_id, content, published_at, price_credits, read_count`,
		ideaID, content, priceCredits,
	).Scan(&p.ID, &p.IdeaID, &p.Content, &p.PublishedAt, &p.PriceCredits, &p.ReadCount)
	if err != nil {
		return nil, fmt.Errorf("create prd: %w", err)
	}
	return &p, nil
}

// GetPRDByIdeaID retrieves the PRD for an idea.
func (s *Store) GetPRDByIdeaID(ctx context.Context, ideaID int64) (*model.PRD, error) {
	var p model.PRD
	err := s.db.QueryRow(ctx,
		`SELECT id, idea_id, content, published_at, price_credits, read_count
		 FROM prds WHERE idea_id = $1`, ideaID,
	).Scan(&p.ID, &p.IdeaID, &p.Content, &p.PublishedAt, &p.PriceCredits, &p.ReadCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get prd by idea_id: %w", err)
	}
	return &p, nil
}

// GetPRDByID retrieves a PRD by its ID.
func (s *Store) GetPRDByID(ctx context.Context, id int64) (*model.PRD, error) {
	var p model.PRD
	err := s.db.QueryRow(ctx,
		`SELECT id, idea_id, content, published_at, price_credits, read_count
		 FROM prds WHERE id = $1`, id,
	).Scan(&p.ID, &p.IdeaID, &p.Content, &p.PublishedAt, &p.PriceCredits, &p.ReadCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get prd by id: %w", err)
	}
	return &p, nil
}

// IncrementPRDReadCount atomically increments the read count of a PRD.
func (s *Store) IncrementPRDReadCount(ctx context.Context, prdID int64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE prds SET read_count = read_count + 1 WHERE id = $1`, prdID,
	)
	if err != nil {
		return fmt.Errorf("increment prd read count: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasUserPurchasedPRD checks whether a user has a credit transaction
// recording a purchase for the given PRD.
func (s *Store) HasUserPurchasedPRD(ctx context.Context, userID, prdID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM credit_transactions
			WHERE user_id = $1 AND reference_type = 'prd' AND reference_id = $2
		)`, userID, prdID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("has user purchased prd: %w", err)
	}
	return exists, nil
}
