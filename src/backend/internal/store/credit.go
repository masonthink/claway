package store

import (
	"context"
	"fmt"

	"github.com/clawbeach/server/internal/model"
)

// CreateCreditTransaction records a credit transaction for a user.
func (s *Store) CreateCreditTransaction(ctx context.Context, tx *model.CreditTransaction) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO credit_transactions (user_id, type, amount, reference_type, reference_id, description)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		tx.UserID, tx.Type, tx.Amount, tx.ReferenceType, tx.ReferenceID, tx.Description,
	)
	if err != nil {
		return fmt.Errorf("create credit transaction: %w", err)
	}
	return nil
}

// GetCreditTransactionsByUserID returns paginated credit transactions for a user,
// ordered by most recent first.
func (s *Store) GetCreditTransactionsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.CreditTransaction, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, type, amount, reference_type, reference_id, description, created_at
		 FROM credit_transactions WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get credit transactions: %w", err)
	}
	defer rows.Close()

	var txns []*model.CreditTransaction
	for rows.Next() {
		var t model.CreditTransaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &t.Amount, &t.ReferenceType, &t.ReferenceID, &t.Description, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("credit transactions scan: %w", err)
		}
		txns = append(txns, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("credit transactions rows: %w", err)
	}
	return txns, nil
}
