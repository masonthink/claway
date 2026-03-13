package store

import (
	"context"
	"fmt"
)

// IncrementRateLimit atomically increments the daily counter for a user+action pair.
// Returns the new count after the increment.
func (s *Store) IncrementRateLimit(ctx context.Context, userID int64, action string) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`INSERT INTO rate_limits (user_id, action, action_date, count)
		 VALUES ($1, $2, CURRENT_DATE, 1)
		 ON CONFLICT (user_id, action, action_date)
		 DO UPDATE SET count = rate_limits.count + 1
		 RETURNING count`,
		userID, action,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("increment rate limit: %w", err)
	}
	return count, nil
}

// GetRateLimitCount returns the current daily count for a user+action pair.
// Returns 0 if no record exists for today.
func (s *Store) GetRateLimitCount(ctx context.Context, userID int64, action string) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COALESCE(
			(SELECT count FROM rate_limits
			 WHERE user_id = $1 AND action = $2 AND action_date = CURRENT_DATE),
			0
		)`,
		userID, action,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get rate limit count: %w", err)
	}
	return count, nil
}
