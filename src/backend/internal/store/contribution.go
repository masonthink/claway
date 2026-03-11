package store

import (
	"context"
	"fmt"

	"github.com/clawbeach/server/internal/model"
)

// CreateContribution records a user's contribution to a task within an idea.
func (s *Store) CreateContribution(ctx context.Context, c *model.Contribution) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO contributions (idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.IdeaID, c.TaskID, c.UserID, c.CostUSD, c.QualityScore, c.WeightedScore, c.WeightPercent,
	)
	if err != nil {
		return fmt.Errorf("create contribution: %w", err)
	}
	return nil
}

// GetContributionsByUserID returns all contributions for a user.
func (s *Store) GetContributionsByUserID(ctx context.Context, userID int64) ([]*model.Contribution, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent
		 FROM contributions WHERE user_id = $1
		 ORDER BY id DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get contributions by user: %w", err)
	}
	defer rows.Close()

	var contributions []*model.Contribution
	for rows.Next() {
		var c model.Contribution
		if err := rows.Scan(&c.ID, &c.IdeaID, &c.TaskID, &c.UserID, &c.CostUSD, &c.QualityScore, &c.WeightedScore, &c.WeightPercent); err != nil {
			return nil, fmt.Errorf("contributions by user scan: %w", err)
		}
		contributions = append(contributions, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("contributions by user rows: %w", err)
	}
	return contributions, nil
}

// GetContributionsByIdeaID returns all contributions for an idea.
func (s *Store) GetContributionsByIdeaID(ctx context.Context, ideaID int64) ([]*model.Contribution, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, idea_id, task_id, user_id, cost_usd, quality_score, weighted_score, weight_percent
		 FROM contributions WHERE idea_id = $1
		 ORDER BY id ASC`, ideaID,
	)
	if err != nil {
		return nil, fmt.Errorf("get contributions by idea: %w", err)
	}
	defer rows.Close()

	var contributions []*model.Contribution
	for rows.Next() {
		var c model.Contribution
		if err := rows.Scan(&c.ID, &c.IdeaID, &c.TaskID, &c.UserID, &c.CostUSD, &c.QualityScore, &c.WeightedScore, &c.WeightPercent); err != nil {
			return nil, fmt.Errorf("contributions by idea scan: %w", err)
		}
		contributions = append(contributions, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("contributions by idea rows: %w", err)
	}
	return contributions, nil
}
