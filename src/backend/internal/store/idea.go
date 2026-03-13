package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

const ideaColumns = `id, initiator_id, title, description, target_user, core_problem, out_of_scope, status, deadline, revealed_at, created_at`

func scanIdea(row pgx.Row) (*model.Idea, error) {
	var i model.Idea
	err := row.Scan(
		&i.ID, &i.InitiatorID, &i.Title, &i.Description,
		&i.TargetUser, &i.CoreProblem, &i.OutOfScope,
		&i.Status, &i.Deadline, &i.RevealedAt, &i.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &i, nil
}

// CreateIdea inserts a new idea and returns it with generated fields populated.
func (s *Store) CreateIdea(ctx context.Context, idea *model.Idea) (*model.Idea, error) {
	i, err := scanIdea(s.db.QueryRow(ctx,
		`INSERT INTO ideas (initiator_id, title, description, target_user, core_problem, out_of_scope, status, deadline)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING `+ideaColumns,
		idea.InitiatorID, idea.Title, idea.Description,
		idea.TargetUser, idea.CoreProblem, idea.OutOfScope,
		idea.Status, idea.Deadline,
	))
	if err != nil {
		return nil, fmt.Errorf("create idea: %w", err)
	}
	return i, nil
}

// GetIdeaByID retrieves an idea by ID.
func (s *Store) GetIdeaByID(ctx context.Context, id int64) (*model.Idea, error) {
	i, err := scanIdea(s.db.QueryRow(ctx,
		`SELECT `+ideaColumns+` FROM ideas WHERE id = $1`, id))
	if err != nil {
		return nil, fmt.Errorf("get idea by id: %w", err)
	}
	return i, nil
}

// ListIdeas returns a paginated list of ideas filtered by status.
// If status is empty, all ideas are returned.
func (s *Store) ListIdeas(ctx context.Context, status string, limit, offset int) ([]*model.Idea, int, error) {
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

	var rows pgx.Rows
	var err error
	if status != "" {
		rows, err = s.db.Query(ctx,
			`SELECT `+ideaColumns+` FROM ideas WHERE status = $1
			 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			status, limit, offset,
		)
	} else {
		rows, err = s.db.Query(ctx,
			`SELECT `+ideaColumns+` FROM ideas
			 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
			limit, offset,
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
			&i.ID, &i.InitiatorID, &i.Title, &i.Description,
			&i.TargetUser, &i.CoreProblem, &i.OutOfScope,
			&i.Status, &i.Deadline, &i.RevealedAt, &i.CreatedAt,
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

// ListIdeasByInitiator returns all ideas created by a specific user.
func (s *Store) ListIdeasByInitiator(ctx context.Context, userID int64, limit, offset int) ([]*model.Idea, int, error) {
	var total int
	if err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM ideas WHERE initiator_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("list ideas by initiator count: %w", err)
	}

	rows, err := s.db.Query(ctx,
		`SELECT `+ideaColumns+` FROM ideas WHERE initiator_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list ideas by initiator query: %w", err)
	}
	defer rows.Close()

	var ideas []*model.Idea
	for rows.Next() {
		var i model.Idea
		if err := rows.Scan(
			&i.ID, &i.InitiatorID, &i.Title, &i.Description,
			&i.TargetUser, &i.CoreProblem, &i.OutOfScope,
			&i.Status, &i.Deadline, &i.RevealedAt, &i.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("list ideas by initiator scan: %w", err)
		}
		ideas = append(ideas, &i)
	}
	return ideas, total, nil
}

// ListExpiredOpenIdeas returns all open ideas past their deadline (for reveal processing).
func (s *Store) ListExpiredOpenIdeas(ctx context.Context) ([]*model.Idea, error) {
	rows, err := s.db.Query(ctx,
		`SELECT `+ideaColumns+` FROM ideas
		 WHERE status = 'open' AND deadline <= $1
		 ORDER BY deadline ASC`,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("list expired open ideas: %w", err)
	}
	defer rows.Close()

	var ideas []*model.Idea
	for rows.Next() {
		var i model.Idea
		if err := rows.Scan(
			&i.ID, &i.InitiatorID, &i.Title, &i.Description,
			&i.TargetUser, &i.CoreProblem, &i.OutOfScope,
			&i.Status, &i.Deadline, &i.RevealedAt, &i.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("list expired open ideas scan: %w", err)
		}
		ideas = append(ideas, &i)
	}
	return ideas, nil
}

// CloseIdea sets an idea's status to closed and records the reveal timestamp.
func (s *Store) CloseIdea(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE ideas SET status = 'closed', revealed_at = NOW()
		 WHERE id = $1 AND status = 'open'`, id,
	)
	if err != nil {
		return fmt.Errorf("close idea: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
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
