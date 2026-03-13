package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
)

// CreateIdeaRequest represents the request body for creating an idea.
type CreateIdeaRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetUser  string `json:"target_user"`
	CoreProblem string `json:"core_problem"`
	OutOfScope  string `json:"out_of_scope"`
}

// IdeaResponse wraps an idea with aggregated counts.
type IdeaResponse struct {
	*model.Idea
	ContributionCount int    `json:"contribution_count"`
	VoterCount        int    `json:"voter_count"`
	InitiatorUsername string `json:"initiator_username"`
}

// CreateIdea validates input, checks daily rate limit, and creates an idea with a 7-day deadline.
func (s *Service) CreateIdea(ctx context.Context, userID int64, req CreateIdeaRequest) (*IdeaResponse, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if req.TargetUser == "" {
		return nil, fmt.Errorf("target_user is required")
	}
	if req.CoreProblem == "" {
		return nil, fmt.Errorf("core_problem is required")
	}

	// Check daily rate limit: max 2 ideas per day
	count, err := s.store.IncrementRateLimit(ctx, userID, "post_idea")
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}
	if count > 2 {
		return nil, fmt.Errorf("daily idea limit reached (max 2 per day)")
	}

	idea := &model.Idea{
		InitiatorID: userID,
		Title:       req.Title,
		Description: req.Description,
		TargetUser:  req.TargetUser,
		CoreProblem: req.CoreProblem,
		Status:      model.IdeaStatusOpen,
		Deadline:    time.Now().Add(7 * 24 * time.Hour),
	}
	if req.OutOfScope != "" {
		idea.OutOfScope = model.NullString{}
		idea.OutOfScope.String = req.OutOfScope
		idea.OutOfScope.Valid = true
	}

	created, err := s.store.CreateIdea(ctx, idea)
	if err != nil {
		return nil, fmt.Errorf("failed to create idea: %w", err)
	}

	user, _ := s.store.GetUserByID(ctx, userID)
	username := ""
	if user != nil {
		username = user.Username
	}

	return &IdeaResponse{
		Idea:              created,
		ContributionCount: 0,
		VoterCount:        0,
		InitiatorUsername:  username,
	}, nil
}

// ListIdeasResponse wraps a list of ideas with total count for pagination.
type ListIdeasResponse struct {
	Ideas []*IdeaResponse `json:"ideas"`
	Total int             `json:"total"`
}

// ListIdeas returns ideas filtered by status with pagination, enriched with counts.
func (s *Service) ListIdeas(ctx context.Context, status string, limit, offset int) (*ListIdeasResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	ideas, total, err := s.store.ListIdeas(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list ideas: %w", err)
	}

	resp := &ListIdeasResponse{
		Ideas: make([]*IdeaResponse, 0, len(ideas)),
		Total: total,
	}
	for _, idea := range ideas {
		resp.Ideas = append(resp.Ideas, s.enrichIdea(ctx, idea))
	}
	return resp, nil
}

// GetIdea returns an idea by ID enriched with counts.
func (s *Service) GetIdea(ctx context.Context, id int64) (*IdeaResponse, error) {
	idea, err := s.store.GetIdeaByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get idea: %w", err)
	}
	return s.enrichIdea(ctx, idea), nil
}

// ListMyIdeas returns ideas created by the authenticated user.
func (s *Service) ListMyIdeas(ctx context.Context, userID int64, limit, offset int) (*ListIdeasResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	ideas, total, err := s.store.ListIdeasByInitiator(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list my ideas: %w", err)
	}

	resp := &ListIdeasResponse{
		Ideas: make([]*IdeaResponse, 0, len(ideas)),
		Total: total,
	}
	for _, idea := range ideas {
		resp.Ideas = append(resp.Ideas, s.enrichIdea(ctx, idea))
	}
	return resp, nil
}

// enrichIdea adds contribution count, voter count, and initiator username to an idea.
func (s *Service) enrichIdea(ctx context.Context, idea *model.Idea) *IdeaResponse {
	contribCount, _ := s.store.CountContributionsByIdea(ctx, idea.ID)
	voterCount, _ := s.store.CountVotersByIdea(ctx, idea.ID)

	username := ""
	if user, err := s.store.GetUserByID(ctx, idea.InitiatorID); err == nil {
		username = user.Username
	}

	return &IdeaResponse{
		Idea:              idea,
		ContributionCount: contribCount,
		VoterCount:        voterCount,
		InitiatorUsername:  username,
	}
}

// GetIdeaStats returns platform-level stats.
type PlatformStats struct {
	OpenIdeas         int `json:"open_ideas"`
	ClosedIdeas       int `json:"closed_ideas"`
	TotalContributions int `json:"total_contributions"`
}

// GetPlatformStats returns aggregate platform statistics.
func (s *Service) GetPlatformStats(ctx context.Context) (*PlatformStats, error) {
	openIdeas, _, err := s.store.ListIdeas(ctx, "open", 1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to count open ideas: %w", err)
	}
	closedIdeas, _, err := s.store.ListIdeas(ctx, "closed", 1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to count closed ideas: %w", err)
	}

	_, openTotal, _ := s.store.ListIdeas(ctx, "open", 1, 0)
	_, closedTotal, _ := s.store.ListIdeas(ctx, "closed", 1, 0)
	_ = openIdeas
	_ = closedIdeas

	totalContribs, _ := s.store.CountAllSubmittedContributions(ctx)

	return &PlatformStats{
		OpenIdeas:          openTotal,
		ClosedIdeas:        closedTotal,
		TotalContributions: totalContribs,
	}, nil
}

// GetUserProfile returns public profile data for a user by username.
type UserProfile struct {
	User              *model.User `json:"user"`
	IdeaCount         int         `json:"idea_count"`
	ContributionCount int         `json:"contribution_count"`
	FeaturedCount     int         `json:"featured_count"`
}

func (s *Service) GetUserProfile(ctx context.Context, username string) (*UserProfile, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	_, ideaCount, _ := s.store.ListIdeasByInitiator(ctx, user.ID, 1, 0)
	_, contribCount, _ := s.store.ListContributionsByAuthor(ctx, user.ID, 1, 0)

	// TODO: count featured contributions from reveal_snapshots

	return &UserProfile{
		User:              user,
		IdeaCount:         ideaCount,
		ContributionCount: contribCount,
		FeaturedCount:     0,
	}, nil
}
