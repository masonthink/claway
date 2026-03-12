package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

// CreateTask inserts a new task and returns it with generated fields populated.
func (s *Store) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	var t model.Task
	err := s.db.QueryRow(ctx,
		`INSERT INTO tasks (idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint,
		           status, claimed_by, claimed_at, submitted_at, approved_at, output_content, output_note,
		           quality_score, reject_reason, review_feedback, cost_usd_accumulated`,
		task.IdeaID, task.Type, task.Title, task.Description, task.AcceptanceCriteria,
		task.Dependencies, task.TokenLimitHint, task.Status,
	).Scan(
		&t.ID, &t.IdeaID, &t.Type, &t.Title, &t.Description, &t.AcceptanceCriteria,
		&t.Dependencies, &t.TokenLimitHint, &t.Status, &t.ClaimedBy, &t.ClaimedAt,
		&t.SubmittedAt, &t.ApprovedAt, &t.OutputContent, &t.OutputNote,
		&t.QualityScore, &t.RejectReason, &t.ReviewFeedback, &t.CostUSDAccumulated,
	)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return &t, nil
}

// GetTaskByID retrieves a task by ID.
func (s *Store) GetTaskByID(ctx context.Context, id int64) (*model.Task, error) {
	var t model.Task
	err := s.db.QueryRow(ctx,
		`SELECT id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint,
		        status, claimed_by, claimed_at, submitted_at, approved_at, output_content, output_note,
		        quality_score, reject_reason, review_feedback, cost_usd_accumulated
		 FROM tasks WHERE id = $1`, id,
	).Scan(
		&t.ID, &t.IdeaID, &t.Type, &t.Title, &t.Description, &t.AcceptanceCriteria,
		&t.Dependencies, &t.TokenLimitHint, &t.Status, &t.ClaimedBy, &t.ClaimedAt,
		&t.SubmittedAt, &t.ApprovedAt, &t.OutputContent, &t.OutputNote,
		&t.QualityScore, &t.RejectReason, &t.ReviewFeedback, &t.CostUSDAccumulated,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return &t, nil
}

// ListTasksByIdeaID returns all tasks for a given idea.
func (s *Store) ListTasksByIdeaID(ctx context.Context, ideaID int64) ([]*model.Task, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, idea_id, type, title, description, acceptance_criteria, dependencies, token_limit_hint,
		        status, claimed_by, claimed_at, submitted_at, approved_at, output_content, output_note,
		        quality_score, reject_reason, review_feedback, cost_usd_accumulated
		 FROM tasks WHERE idea_id = $1
		 ORDER BY id ASC`, ideaID,
	)
	if err != nil {
		return nil, fmt.Errorf("list tasks by idea: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID, &t.IdeaID, &t.Type, &t.Title, &t.Description, &t.AcceptanceCriteria,
			&t.Dependencies, &t.TokenLimitHint, &t.Status, &t.ClaimedBy, &t.ClaimedAt,
			&t.SubmittedAt, &t.ApprovedAt, &t.OutputContent, &t.OutputNote,
			&t.QualityScore, &t.RejectReason, &t.ReviewFeedback, &t.CostUSDAccumulated,
		); err != nil {
			return nil, fmt.Errorf("list tasks scan: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks rows: %w", err)
	}
	return tasks, nil
}

// ClaimTask assigns a task to a user. Only open tasks can be claimed.
func (s *Store) ClaimTask(ctx context.Context, taskID, userID int64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET status = 'claimed', claimed_by = $1, claimed_at = NOW()
		 WHERE id = $2 AND status = 'open'`, userID, taskID,
	)
	if err != nil {
		return fmt.Errorf("claim task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// UnclaimTask releases a claimed task back to open status.
func (s *Store) UnclaimTask(ctx context.Context, taskID int64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET status = 'open', claimed_by = NULL, claimed_at = NULL
		 WHERE id = $1 AND status = 'claimed'`, taskID,
	)
	if err != nil {
		return fmt.Errorf("unclaim task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// SubmitTask marks a claimed task as submitted with output content and note.
func (s *Store) SubmitTask(ctx context.Context, taskID int64, content, note string) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET status = 'submitted', output_content = $1, output_note = $2, submitted_at = NOW()
		 WHERE id = $3 AND status IN ('claimed', 'rejected', 'revision')`, content, note, taskID,
	)
	if err != nil {
		return fmt.Errorf("submit task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// ApproveTask marks a submitted task as approved with a quality score.
func (s *Store) ApproveTask(ctx context.Context, taskID int64, qualityScore float64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET status = 'approved', quality_score = $1, approved_at = NOW()
		 WHERE id = $2 AND status = 'submitted'`, qualityScore, taskID,
	)
	if err != nil {
		return fmt.Errorf("approve task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// RejectTask marks a submitted task as rejected with a reason.
func (s *Store) RejectTask(ctx context.Context, taskID int64, reason string) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET status = 'rejected', reject_reason = $1
		 WHERE id = $2 AND status = 'submitted'`, reason, taskID,
	)
	if err != nil {
		return fmt.Errorf("reject task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// RevisionTask marks a submitted task as needing revision with feedback.
func (s *Store) RevisionTask(ctx context.Context, taskID int64, feedback string) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET status = 'revision', review_feedback = $1
		 WHERE id = $2 AND status = 'submitted'`, feedback, taskID,
	)
	if err != nil {
		return fmt.Errorf("revision task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// AccumulateTaskCost atomically adds cost to a task's accumulated cost.
func (s *Store) AccumulateTaskCost(ctx context.Context, taskID int64, costUSD float64) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE tasks SET cost_usd_accumulated = cost_usd_accumulated + $1
		 WHERE id = $2`, costUSD, taskID,
	)
	if err != nil {
		return fmt.Errorf("accumulate task cost: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
