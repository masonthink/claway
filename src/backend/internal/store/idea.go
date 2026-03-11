package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/clawbeach/server/internal/model"
	"github.com/jackc/pgx/v5"
)

// CreateIdea inserts a new idea and returns it with generated fields populated.
func (s *Store) CreateIdea(ctx context.Context, idea *model.Idea) (*model.Idea, error) {
	var i model.Idea
	err := s.db.QueryRow(ctx,
		`INSERT INTO ideas (title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, deadline)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at, deadline`,
		idea.Title, idea.Description, idea.TargetUserHint, idea.ProblemDefinition,
		idea.InitiatorID, idea.InitiatorCutPercent, idea.PackageType, idea.Status, idea.Deadline,
	).Scan(
		&i.ID, &i.Title, &i.Description, &i.TargetUserHint, &i.ProblemDefinition,
		&i.InitiatorID, &i.InitiatorCutPercent, &i.PackageType, &i.Status, &i.CreatedAt, &i.Deadline,
	)
	if err != nil {
		return nil, fmt.Errorf("create idea: %w", err)
	}
	return &i, nil
}

// GetIdeaByID retrieves an idea by ID.
func (s *Store) GetIdeaByID(ctx context.Context, id int64) (*model.Idea, error) {
	var i model.Idea
	err := s.db.QueryRow(ctx,
		`SELECT id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at, deadline
		 FROM ideas WHERE id = $1`, id,
	).Scan(
		&i.ID, &i.Title, &i.Description, &i.TargetUserHint, &i.ProblemDefinition,
		&i.InitiatorID, &i.InitiatorCutPercent, &i.PackageType, &i.Status, &i.CreatedAt, &i.Deadline,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get idea by id: %w", err)
	}
	return &i, nil
}

// ListIdeas returns a paginated list of ideas filtered by status.
// If status is empty, all ideas are returned.
// Returns the matching ideas and the total count (for pagination).
func (s *Store) ListIdeas(ctx context.Context, status string, limit, offset int) ([]*model.Idea, int, error) {
	// Count total matching rows.
	var total int
	var countErr error
	if status != "" {
		countErr = s.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM ideas WHERE status = $1`, status,
		).Scan(&total)
	} else {
		countErr = s.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM ideas`,
		).Scan(&total)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("list ideas count: %w", countErr)
	}

	// Fetch the page.
	var rows pgx.Rows
	var err error
	if status != "" {
		rows, err = s.db.Query(ctx,
			`SELECT id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at, deadline
			 FROM ideas WHERE status = $1
			 ORDER BY created_at DESC
			 LIMIT $2 OFFSET $3`, status, limit, offset,
		)
	} else {
		rows, err = s.db.Query(ctx,
			`SELECT id, title, description, target_user_hint, problem_definition, initiator_id, initiator_cut_percent, package_type, status, created_at, deadline
			 FROM ideas
			 ORDER BY created_at DESC
			 LIMIT $1 OFFSET $2`, limit, offset,
		)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("list ideas query: %w", err)
	}
	defer rows.Close()

	var ideas []*model.Idea
	for rows.Next() {
		var i model.Idea
		if err := rows.Scan(
			&i.ID, &i.Title, &i.Description, &i.TargetUserHint, &i.ProblemDefinition,
			&i.InitiatorID, &i.InitiatorCutPercent, &i.PackageType, &i.Status, &i.CreatedAt, &i.Deadline,
		); err != nil {
			return nil, 0, fmt.Errorf("list ideas scan: %w", err)
		}
		ideas = append(ideas, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list ideas rows: %w", err)
	}
	return ideas, total, nil
}

// UpdateIdeaStatus updates the status of an idea.
func (s *Store) UpdateIdeaStatus(ctx context.Context, id int64, status string) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE ideas SET status = $1 WHERE id = $2`, status, id,
	)
	if err != nil {
		return fmt.Errorf("update idea status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
