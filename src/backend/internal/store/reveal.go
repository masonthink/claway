package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

// CreateRevealSnapshot stores the frozen ranking data for an idea.
func (s *Store) CreateRevealSnapshot(ctx context.Context, snap *model.RevealSnapshot) (*model.RevealSnapshot, error) {
	var result model.RevealSnapshot
	err := s.db.QueryRow(ctx,
		`INSERT INTO reveal_snapshots (idea_id, ranked_results, total_votes)
		 VALUES ($1, $2, $3)
		 RETURNING id, idea_id, ranked_results, total_votes, revealed_at`,
		snap.IdeaID, snap.RankedResults, snap.TotalVotes,
	).Scan(&result.ID, &result.IdeaID, &result.RankedResults, &result.TotalVotes, &result.RevealedAt)
	if err != nil {
		return nil, fmt.Errorf("create reveal snapshot: %w", err)
	}
	return &result, nil
}

// CountFeaturedByUser counts how many times a user's contributions were featured.
func (s *Store) CountFeaturedByUser(ctx context.Context, userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM reveal_snapshots rs
		 CROSS JOIN LATERAL jsonb_array_elements(rs.ranked_results) AS elem
		 JOIN contributions c ON c.id = (elem->>'contribution_id')::bigint
		 WHERE (elem->>'is_featured')::boolean = true
		   AND c.author_id = $1`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count featured by user: %w", err)
	}
	return count, nil
}

// GetRevealSnapshotByIdeaID retrieves the reveal snapshot for an idea.
func (s *Store) GetRevealSnapshotByIdeaID(ctx context.Context, ideaID int64) (*model.RevealSnapshot, error) {
	var snap model.RevealSnapshot
	err := s.db.QueryRow(ctx,
		`SELECT id, idea_id, ranked_results, total_votes, revealed_at
		 FROM reveal_snapshots WHERE idea_id = $1`, ideaID,
	).Scan(&snap.ID, &snap.IdeaID, &snap.RankedResults, &snap.TotalVotes, &snap.RevealedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get reveal snapshot: %w", err)
	}
	return &snap, nil
}
