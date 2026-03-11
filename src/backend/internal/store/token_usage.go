package store

import (
	"context"
	"fmt"

	"github.com/clawbeach/server/internal/model"
)

// CreateTokenUsageLog records a single LLM API call's token usage.
func (s *Store) CreateTokenUsageLog(ctx context.Context, log *model.TokenUsageLog) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO token_usage_logs (user_id, task_id, model, tokens_in, tokens_out, cost_usd)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		log.UserID, log.TaskID, log.Model, log.TokensIn, log.TokensOut, log.CostUSD,
	)
	if err != nil {
		return fmt.Errorf("create token usage log: %w", err)
	}
	return nil
}

// GetUserComputeTotal returns the total compute cost for a user across all tasks.
func (s *Store) GetUserComputeTotal(ctx context.Context, userID int64) (float64, error) {
	var total float64
	err := s.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(cost_usd), 0) FROM token_usage_logs WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get user compute total: %w", err)
	}
	return total, nil
}

// GetUserIdeaCompute returns the total compute cost for a user on a specific idea.
func (s *Store) GetUserIdeaCompute(ctx context.Context, userID, ideaID int64) (float64, error) {
	var total float64
	err := s.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(tul.cost_usd), 0)
		 FROM token_usage_logs tul
		 JOIN tasks t ON t.id = tul.task_id
		 WHERE tul.user_id = $1 AND t.idea_id = $2`, userID, ideaID,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get user idea compute: %w", err)
	}
	return total, nil
}

// GetIdeaComputeBreakdown returns per-user compute usage for an idea.
func (s *Store) GetIdeaComputeBreakdown(ctx context.Context, ideaID int64) ([]ComputeEntry, error) {
	rows, err := s.db.Query(ctx,
		`SELECT tul.user_id, u.username,
		        COALESCE(SUM(tul.cost_usd), 0) AS total_cost,
		        COALESCE(SUM(tul.tokens_in), 0) AS total_tokens_in,
		        COALESCE(SUM(tul.tokens_out), 0) AS total_tokens_out,
		        COUNT(*) AS call_count
		 FROM token_usage_logs tul
		 JOIN tasks t ON t.id = tul.task_id
		 JOIN users u ON u.id = tul.user_id
		 WHERE t.idea_id = $1
		 GROUP BY tul.user_id, u.username
		 ORDER BY total_cost DESC`, ideaID,
	)
	if err != nil {
		return nil, fmt.Errorf("get idea compute breakdown: %w", err)
	}
	defer rows.Close()

	var entries []ComputeEntry
	for rows.Next() {
		var e ComputeEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.TotalCost, &e.TotalTokensIn, &e.TotalTokensOut, &e.CallCount); err != nil {
			return nil, fmt.Errorf("idea compute breakdown scan: %w", err)
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("idea compute breakdown rows: %w", err)
	}
	return entries, nil
}

// GetTaskComputeBreakdown returns per-user compute usage for a specific task.
func (s *Store) GetTaskComputeBreakdown(ctx context.Context, taskID int64) ([]ComputeEntry, error) {
	rows, err := s.db.Query(ctx,
		`SELECT tul.user_id, u.username,
		        COALESCE(SUM(tul.cost_usd), 0) AS total_cost,
		        COALESCE(SUM(tul.tokens_in), 0) AS total_tokens_in,
		        COALESCE(SUM(tul.tokens_out), 0) AS total_tokens_out,
		        COUNT(*) AS call_count
		 FROM token_usage_logs tul
		 JOIN users u ON u.id = tul.user_id
		 WHERE tul.task_id = $1
		 GROUP BY tul.user_id, u.username
		 ORDER BY total_cost DESC`, taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("get task compute breakdown: %w", err)
	}
	defer rows.Close()

	var entries []ComputeEntry
	for rows.Next() {
		var e ComputeEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.TotalCost, &e.TotalTokensIn, &e.TotalTokensOut, &e.CallCount); err != nil {
			return nil, fmt.Errorf("task compute breakdown scan: %w", err)
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("task compute breakdown rows: %w", err)
	}
	return entries, nil
}

// GetPlatformComputeTotal returns the total compute cost and token counts across the entire platform.
func (s *Store) GetPlatformComputeTotal(ctx context.Context) (float64, int, error) {
	var totalCost float64
	var totalTokens int
	err := s.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(cost_usd), 0), COALESCE(SUM(tokens_in + tokens_out), 0)
		 FROM token_usage_logs`,
	).Scan(&totalCost, &totalTokens)
	if err != nil {
		return 0, 0, fmt.Errorf("get platform compute total: %w", err)
	}
	return totalCost, totalTokens, nil
}
