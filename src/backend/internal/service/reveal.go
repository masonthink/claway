package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
)

// RevealResultEntry is a single entry in the reveal result API response.
type RevealResultEntry struct {
	ContributionID int64  `json:"contribution_id"`
	AuthorID       int64  `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	VoteCount      int    `json:"vote_count"`
	Rank           int    `json:"rank"`
	IsFeatured     bool   `json:"is_featured"`
}

// RevealResultResponse wraps the full reveal result for an idea.
type RevealResultResponse struct {
	IdeaID     int64               `json:"idea_id"`
	TotalVotes int                 `json:"total_votes"`
	RevealedAt string              `json:"revealed_at"`
	Results    []RevealResultEntry `json:"results"`
}

// GetRevealResult returns the reveal result for a closed idea.
func (s *Service) GetRevealResult(ctx context.Context, ideaID int64) (*RevealResultResponse, error) {
	idea, err := s.store.GetIdeaByID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("idea not found")
	}
	if idea.Status != model.IdeaStatusClosed {
		return nil, fmt.Errorf("idea has not been revealed yet")
	}

	snap, err := s.store.GetRevealSnapshotByIdeaID(ctx, ideaID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("reveal data not found")
		}
		return nil, fmt.Errorf("failed to get reveal data: %w", err)
	}

	var ranked []model.RankedResult
	if err := json.Unmarshal(snap.RankedResults, &ranked); err != nil {
		return nil, fmt.Errorf("failed to parse ranked results: %w", err)
	}

	results := make([]RevealResultEntry, 0, len(ranked))
	for _, r := range ranked {
		entry := RevealResultEntry{
			ContributionID: r.ContributionID,
			VoteCount:      r.VoteCount,
			Rank:           r.Rank,
			IsFeatured:     r.IsFeatured,
		}

		// Look up author info
		contrib, err := s.store.GetContributionByID(ctx, r.ContributionID)
		if err == nil {
			entry.AuthorID = contrib.AuthorID
			if user, err := s.store.GetUserByID(ctx, contrib.AuthorID); err == nil {
				entry.AuthorUsername = user.Username
			}
		}

		results = append(results, entry)
	}

	return &RevealResultResponse{
		IdeaID:     ideaID,
		TotalVotes: snap.TotalVotes,
		RevealedAt: snap.RevealedAt.Format("2006-01-02T15:04:05Z"),
		Results:    results,
	}, nil
}

// ProcessReveal runs the reveal logic for a single idea.
// 1. Count votes per contribution
// 2. Rank by vote count DESC, then submission time ASC
// 3. If total_votes >= 5, top 3 are featured (ties included)
// 4. Write reveal_snapshot
// 5. Close the idea
func (s *Service) ProcessReveal(ctx context.Context, ideaID int64) error {
	voteCounts, totalVotes, err := s.store.GetVoteCountsByIdea(ctx, ideaID)
	if err != nil {
		return fmt.Errorf("get vote counts: %w", err)
	}

	// Build ranked results
	ranked := make([]model.RankedResult, 0, len(voteCounts))
	for i, vc := range voteCounts {
		rank := i + 1
		// Handle ties: same vote count = same rank
		if i > 0 && vc.VoteCount == voteCounts[i-1].VoteCount {
			rank = ranked[i-1].Rank
		}

		isFeatured := false
		if totalVotes >= 5 && rank <= 3 {
			isFeatured = true
		}

		ranked = append(ranked, model.RankedResult{
			ContributionID: vc.ContributionID,
			VoteCount:      vc.VoteCount,
			Rank:           rank,
			IsFeatured:     isFeatured,
		})
	}

	rankedJSON, err := json.Marshal(ranked)
	if err != nil {
		return fmt.Errorf("marshal ranked results: %w", err)
	}

	snap := &model.RevealSnapshot{
		IdeaID:        ideaID,
		RankedResults: rankedJSON,
		TotalVotes:    totalVotes,
	}

	if _, err := s.store.CreateRevealSnapshot(ctx, snap); err != nil {
		return fmt.Errorf("create reveal snapshot: %w", err)
	}

	if err := s.store.CloseIdea(ctx, ideaID); err != nil {
		return fmt.Errorf("close idea: %w", err)
	}

	return nil
}

// RunRevealTicker starts a background goroutine that checks for expired ideas
// every interval and processes their reveals. Stops when context is cancelled.
func (s *Service) RunRevealTicker(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("reveal ticker stopped")
				return
			case <-ticker.C:
				s.processExpiredIdeas(ctx)
			}
		}
	}()
}

const maxRevealRetries = 3

func (s *Service) processExpiredIdeas(ctx context.Context) {
	ideas, err := s.store.ListExpiredOpenIdeas(ctx)
	if err != nil {
		log.Printf("error listing expired ideas: %v", err)
		return
	}

	for _, idea := range ideas {
		log.Printf("processing reveal for idea %d: %s", idea.ID, idea.Title)
		if err := s.revealWithRetry(ctx, idea.ID); err != nil {
			log.Printf("CRITICAL: reveal failed after %d retries for idea %d: %v", maxRevealRetries, idea.ID, err)
		}
	}
}

// revealWithRetry attempts ProcessReveal with exponential backoff.
// Retries up to maxRevealRetries times (3 attempts total, with 1s, 2s backoff delays).
func (s *Service) revealWithRetry(ctx context.Context, ideaID int64) error {
	var lastErr error
	for attempt := 0; attempt < maxRevealRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("retrying reveal for idea %d (attempt %d/%d, delay %v)", ideaID, attempt+1, maxRevealRetries, delay)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
		if err := s.ProcessReveal(ctx, ideaID); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return fmt.Errorf("all %d attempts failed: %w", maxRevealRetries, lastErr)
}
