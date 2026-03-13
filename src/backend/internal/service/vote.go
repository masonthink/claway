package service

import (
	"context"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
)

// CastVoteRequest represents the request body for casting a vote.
type CastVoteRequest struct {
	ContributionID int64 `json:"contribution_id"`
}

// CastVote validates and records a vote.
func (s *Service) CastVote(ctx context.Context, userID int64, ideaID int64, req CastVoteRequest) (*model.Vote, error) {
	if req.ContributionID == 0 {
		return nil, fmt.Errorf("contribution_id is required")
	}

	// Check idea exists and is open
	idea, err := s.store.GetIdeaByID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("idea not found")
	}
	if idea.Status != model.IdeaStatusOpen {
		return nil, fmt.Errorf("voting has ended for this idea")
	}

	// Check contribution exists, belongs to this idea, and is submitted
	contrib, err := s.store.GetContributionByID(ctx, req.ContributionID)
	if err != nil {
		return nil, fmt.Errorf("contribution not found")
	}
	if contrib.IdeaID != ideaID {
		return nil, fmt.Errorf("contribution does not belong to this idea")
	}
	if contrib.Status != model.ContributionStatusSubmitted {
		return nil, fmt.Errorf("cannot vote for a draft contribution")
	}

	// Self-vote check (application layer)
	if contrib.AuthorID == userID {
		return nil, fmt.Errorf("cannot vote for your own contribution")
	}

	// Check if user already voted on this idea
	alreadyVoted, err := s.store.HasUserVotedOnIdea(ctx, ideaID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing vote: %w", err)
	}
	if alreadyVoted {
		return nil, fmt.Errorf("you have already voted on this idea")
	}

	// Check daily vote rate limit: max 10 per day
	count, err := s.store.IncrementRateLimit(ctx, userID, "vote")
	if err != nil {
		return nil, fmt.Errorf("failed to check vote rate limit: %w", err)
	}
	if count > 10 {
		return nil, fmt.Errorf("daily vote limit reached (max 10 per day)")
	}

	vote := &model.Vote{
		IdeaID:         ideaID,
		VoterID:        userID,
		ContributionID: req.ContributionID,
	}

	created, err := s.store.CreateVote(ctx, vote)
	if err != nil {
		if err == store.ErrConflict {
			return nil, fmt.Errorf("you have already voted on this idea")
		}
		return nil, fmt.Errorf("failed to cast vote: %w", err)
	}

	return created, nil
}

// VoteResponse wraps a vote for API responses.
type VoteResponse struct {
	ID             int64  `json:"id"`
	IdeaID         int64  `json:"idea_id"`
	ContributionID int64  `json:"contribution_id"`
	VotedAt        string `json:"voted_at"`
}

// ListMyVotes returns all votes cast by the authenticated user.
func (s *Service) ListMyVotes(ctx context.Context, userID int64, limit, offset int) ([]*VoteResponse, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	votes, total, err := s.store.ListVotesByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list my votes: %w", err)
	}

	result := make([]*VoteResponse, 0, len(votes))
	for _, v := range votes {
		result = append(result, &VoteResponse{
			ID:             v.ID,
			IdeaID:         v.IdeaID,
			ContributionID: v.ContributionID,
			VotedAt:        v.VotedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result, total, nil
}
