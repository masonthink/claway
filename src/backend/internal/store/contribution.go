package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

const contributionColumns = `id, idea_id, author_id, content, decision_log, status, view_count, created_at, updated_at, submitted_at`

func scanContribution(row pgx.Row) (*model.Contribution, error) {
	var c model.Contribution
	err := row.Scan(
		&c.ID, &c.IdeaID, &c.AuthorID, &c.Content, &c.DecisionLog,
		&c.Status, &c.ViewCount, &c.CreatedAt, &c.UpdatedAt, &c.SubmittedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

// CreateContribution inserts a new contribution (draft) for an idea.
func (s *Store) CreateContribution(ctx context.Context, c *model.Contribution) (*model.Contribution, error) {
	result, err := scanContribution(s.db.QueryRow(ctx,
		`INSERT INTO contributions (idea_id, author_id, content, decision_log, status)
		 VALUES ($1, $2, $3, $4, 'draft')
		 RETURNING `+contributionColumns,
		c.IdeaID, c.AuthorID, c.Content, c.DecisionLog,
	))
	if err != nil {
		return nil, fmt.Errorf("create contribution: %w", err)
	}
	return result, nil
}

// GetContributionByID retrieves a contribution by ID.
func (s *Store) GetContributionByID(ctx context.Context, id int64) (*model.Contribution, error) {
	c, err := scanContribution(s.db.QueryRow(ctx,
		`SELECT `+contributionColumns+` FROM contributions WHERE id = $1`, id))
	if err != nil {
		return nil, fmt.Errorf("get contribution by id: %w", err)
	}
	return c, nil
}

// GetContributionByIdeaAndAuthor returns the existing contribution (draft or submitted) for a user+idea pair.
func (s *Store) GetContributionByIdeaAndAuthor(ctx context.Context, ideaID, authorID int64) (*model.Contribution, error) {
	c, err := scanContribution(s.db.QueryRow(ctx,
		`SELECT `+contributionColumns+` FROM contributions
		 WHERE idea_id = $1 AND author_id = $2`, ideaID, authorID))
	if err != nil {
		return nil, fmt.Errorf("get contribution by idea and author: %w", err)
	}
	return c, nil
}

// UpdateContributionDraft updates a draft contribution's content and decision log.
func (s *Store) UpdateContributionDraft(ctx context.Context, id int64, content string, decisionLog []byte) (*model.Contribution, error) {
	c, err := scanContribution(s.db.QueryRow(ctx,
		`UPDATE contributions
		 SET content = $1, decision_log = COALESCE($2, decision_log), updated_at = NOW()
		 WHERE id = $3 AND status = 'draft'
		 RETURNING `+contributionColumns,
		content, decisionLog, id,
	))
	if err != nil {
		return nil, fmt.Errorf("update contribution draft: %w", err)
	}
	return c, nil
}

// SubmitContribution transitions a draft contribution to submitted status.
func (s *Store) SubmitContribution(ctx context.Context, id int64) (*model.Contribution, error) {
	c, err := scanContribution(s.db.QueryRow(ctx,
		`UPDATE contributions
		 SET status = 'submitted', submitted_at = NOW(), updated_at = NOW()
		 WHERE id = $1 AND status = 'draft'
		 RETURNING `+contributionColumns,
		id,
	))
	if err != nil {
		return nil, fmt.Errorf("submit contribution: %w", err)
	}
	return c, nil
}

// ListContributionsByIdea returns submitted contributions for an idea.
// When the idea is open (pre-reveal), results are randomly ordered.
// When closed (post-reveal), results are ordered by submitted_at.
func (s *Store) ListContributionsByIdea(ctx context.Context, ideaID int64, ideaStatus string) ([]*model.Contribution, error) {
	orderClause := "ORDER BY submitted_at ASC"
	if ideaStatus == "open" {
		orderClause = "ORDER BY RANDOM()"
	}

	rows, err := s.db.Query(ctx,
		`SELECT `+contributionColumns+` FROM contributions
		 WHERE idea_id = $1 AND status = 'submitted'
		 `+orderClause,
		ideaID,
	)
	if err != nil {
		return nil, fmt.Errorf("list contributions by idea: %w", err)
	}
	defer rows.Close()

	var contributions []*model.Contribution
	for rows.Next() {
		var c model.Contribution
		if err := rows.Scan(
			&c.ID, &c.IdeaID, &c.AuthorID, &c.Content, &c.DecisionLog,
			&c.Status, &c.ViewCount, &c.CreatedAt, &c.UpdatedAt, &c.SubmittedAt,
		); err != nil {
			return nil, fmt.Errorf("list contributions scan: %w", err)
		}
		contributions = append(contributions, &c)
	}
	return contributions, nil
}

// ListContributionsByAuthor returns all contributions by a specific user.
func (s *Store) ListContributionsByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*model.Contribution, int, error) {
	var total int
	if err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM contributions WHERE author_id = $1`, authorID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("list contributions by author count: %w", err)
	}

	rows, err := s.db.Query(ctx,
		`SELECT `+contributionColumns+` FROM contributions
		 WHERE author_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		authorID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list contributions by author query: %w", err)
	}
	defer rows.Close()

	var contributions []*model.Contribution
	for rows.Next() {
		var c model.Contribution
		if err := rows.Scan(
			&c.ID, &c.IdeaID, &c.AuthorID, &c.Content, &c.DecisionLog,
			&c.Status, &c.ViewCount, &c.CreatedAt, &c.UpdatedAt, &c.SubmittedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("list contributions by author scan: %w", err)
		}
		contributions = append(contributions, &c)
	}
	return contributions, total, nil
}

// CountContributionsByIdeaIDs returns submitted contribution counts for multiple ideas in one query.
func (s *Store) CountContributionsByIdeaIDs(ctx context.Context, ideaIDs []int64) (map[int64]int, error) {
	if len(ideaIDs) == 0 {
		return make(map[int64]int), nil
	}
	rows, err := s.db.Query(ctx,
		`SELECT idea_id, COUNT(*) FROM contributions
		 WHERE idea_id = ANY($1) AND status = 'submitted'
		 GROUP BY idea_id`, ideaIDs)
	if err != nil {
		return nil, fmt.Errorf("count contributions by idea ids: %w", err)
	}
	defer rows.Close()

	result := make(map[int64]int, len(ideaIDs))
	for rows.Next() {
		var id int64
		var count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, fmt.Errorf("count contributions by idea ids scan: %w", err)
		}
		result[id] = count
	}
	return result, nil
}

// CountContributionsByIdea returns the number of submitted contributions for an idea.
func (s *Store) CountContributionsByIdea(ctx context.Context, ideaID int64) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM contributions WHERE idea_id = $1 AND status = 'submitted'`, ideaID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count contributions by idea: %w", err)
	}
	return count, nil
}

// IncrementContributionViewCount atomically increments the view count.
func (s *Store) IncrementContributionViewCount(ctx context.Context, id int64) error {
	_, err := s.db.Exec(ctx,
		`UPDATE contributions SET view_count = view_count + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("increment contribution view count: %w", err)
	}
	return nil
}

// CountAllSubmittedContributions returns the total number of submitted contributions platform-wide.
func (s *Store) CountAllSubmittedContributions(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM contributions WHERE status = 'submitted'`,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count all submitted contributions: %w", err)
	}
	return count, nil
}
