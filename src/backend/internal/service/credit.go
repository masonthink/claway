package service

import (
	"context"
	"fmt"

	"github.com/clawbeach/server/internal/model"
)

// CreditsResponse contains user balance and recent transactions.
type CreditsResponse struct {
	Balance      float64                    `json:"balance"`
	Transactions []*model.CreditTransaction `json:"transactions"`
}

// GetMyCredits returns the current user's credit balance and recent transactions.
func (s *Service) GetMyCredits(ctx context.Context, userID int64, limit, offset int) (*CreditsResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	txs, err := s.store.GetCreditTransactionsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return &CreditsResponse{
		Balance:      user.CreditsBalance,
		Transactions: txs,
	}, nil
}

// GetMyContributions returns all contributions for the current user.
func (s *Service) GetMyContributions(ctx context.Context, userID int64) ([]*model.Contribution, error) {
	contribs, err := s.store.GetContributionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributions: %w", err)
	}
	return contribs, nil
}

// PurchasePRD allows a user to spend credits to unlock a PRD.
// Credits are distributed: 10% platform, initiator_cut_percent to initiator,
// remainder split among contributors by weight.
func (s *Service) PurchasePRD(ctx context.Context, userID int64, prdID int64) error {
	prd, err := s.store.GetPRDByID(ctx, prdID)
	if err != nil {
		return fmt.Errorf("failed to get PRD: %w", err)
	}

	// Check if already purchased
	purchased, err := s.store.HasUserPurchasedPRD(ctx, userID, prdID)
	if err != nil {
		return fmt.Errorf("failed to check purchase status: %w", err)
	}
	if purchased {
		return fmt.Errorf("you have already purchased this PRD")
	}

	// Check balance
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user.CreditsBalance < prd.PriceCredits {
		return fmt.Errorf("insufficient credits: need %.0f, have %.0f", prd.PriceCredits, user.CreditsBalance)
	}

	idea, err := s.store.GetIdeaByID(ctx, prd.IdeaID)
	if err != nil {
		return fmt.Errorf("failed to get idea: %w", err)
	}

	price := prd.PriceCredits

	// Deduct credits from buyer
	if err := s.store.UpdateCreditsBalance(ctx, userID, -price); err != nil {
		return fmt.Errorf("failed to deduct credits: %w", err)
	}

	// Record buyer's transaction
	buyerTx := &model.CreditTransaction{
		UserID:        userID,
		Type:          "spend_read",
		Amount:        -price,
		ReferenceType: "prd",
		ReferenceID:   prdID,
		Description:   fmt.Sprintf("Purchased PRD: %s", idea.Title),
	}
	if err := s.store.CreateCreditTransaction(ctx, buyerTx); err != nil {
		return fmt.Errorf("failed to record buyer transaction: %w", err)
	}

	// Platform cut: 10%
	platformCut := price * 0.10

	// Initiator cut
	initiatorCut := price * (idea.InitiatorCutPercent / 100.0)

	// Contributors pool
	contributorsPool := price - platformCut - initiatorCut

	// Award initiator
	if initiatorCut > 0 {
		if err := s.store.UpdateCreditsBalance(ctx, idea.InitiatorID, initiatorCut); err != nil {
			return fmt.Errorf("failed to award initiator: %w", err)
		}
		initiatorTx := &model.CreditTransaction{
			UserID:        idea.InitiatorID,
			Type:          "earn_read_share",
			Amount:        initiatorCut,
			ReferenceType: "prd",
			ReferenceID:   prdID,
			Description:   fmt.Sprintf("Initiator share from PRD purchase: %s", idea.Title),
		}
		if err := s.store.CreateCreditTransaction(ctx, initiatorTx); err != nil {
			return fmt.Errorf("failed to record initiator transaction: %w", err)
		}
	}

	// Distribute to contributors by weight
	contribs, err := s.store.GetContributionsByIdeaID(ctx, prd.IdeaID)
	if err != nil {
		return fmt.Errorf("failed to get contributions: %w", err)
	}

	// Calculate total weighted score
	var totalWeighted float64
	for _, c := range contribs {
		totalWeighted += c.WeightedScore
	}

	if totalWeighted > 0 && contributorsPool > 0 {
		for _, c := range contribs {
			share := contributorsPool * (c.WeightedScore / totalWeighted)
			if share <= 0 {
				continue
			}

			if err := s.store.UpdateCreditsBalance(ctx, c.UserID, share); err != nil {
				return fmt.Errorf("failed to award contributor %d: %w", c.UserID, err)
			}

			contribTx := &model.CreditTransaction{
				UserID:        c.UserID,
				Type:          "earn_read_share",
				Amount:        share,
				ReferenceType: "prd",
				ReferenceID:   prdID,
				Description:   fmt.Sprintf("Contributor share from PRD purchase: %s", idea.Title),
			}
			if err := s.store.CreateCreditTransaction(ctx, contribTx); err != nil {
				return fmt.Errorf("failed to record contributor transaction: %w", err)
			}
		}
	}

	// Increment read count
	if err := s.store.IncrementPRDReadCount(ctx, prdID); err != nil {
		return fmt.Errorf("failed to increment read count: %w", err)
	}

	return nil
}
