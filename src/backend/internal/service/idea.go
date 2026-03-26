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

// Input length limits for idea fields.
const (
	maxIdeaTitleLen       = 200
	maxIdeaDescriptionLen = 5000
	maxIdeaTargetUserLen  = 500
	maxIdeaCoreProblemLen = 2000
	maxIdeaOutOfScopeLen  = 2000
)

// CreateIdea validates input, checks daily rate limit, and creates an idea with a 7-day deadline.
func (s *Service) CreateIdea(ctx context.Context, userID int64, req CreateIdeaRequest) (*IdeaResponse, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if len(req.Title) > maxIdeaTitleLen {
		return nil, fmt.Errorf("title must be at most %d characters", maxIdeaTitleLen)
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if len(req.Description) > maxIdeaDescriptionLen {
		return nil, fmt.Errorf("description must be at most %d characters", maxIdeaDescriptionLen)
	}
	if req.TargetUser == "" {
		return nil, fmt.Errorf("target_user is required")
	}
	if len(req.TargetUser) > maxIdeaTargetUserLen {
		return nil, fmt.Errorf("target_user must be at most %d characters", maxIdeaTargetUserLen)
	}
	if req.CoreProblem == "" {
		return nil, fmt.Errorf("core_problem is required")
	}
	if len(req.CoreProblem) > maxIdeaCoreProblemLen {
		return nil, fmt.Errorf("core_problem must be at most %d characters", maxIdeaCoreProblemLen)
	}
	if len(req.OutOfScope) > maxIdeaOutOfScopeLen {
		return nil, fmt.Errorf("out_of_scope must be at most %d characters", maxIdeaOutOfScopeLen)
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

	return &ListIdeasResponse{
		Ideas: s.enrichIdeas(ctx, ideas),
		Total: total,
	}, nil
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

	return &ListIdeasResponse{
		Ideas: s.enrichIdeas(ctx, ideas),
		Total: total,
	}, nil
}

// enrichIdea adds contribution count, voter count, and initiator username to a single idea.
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

// enrichIdeas batch-enriches a list of ideas with 3 queries instead of 3*N.
func (s *Service) enrichIdeas(ctx context.Context, ideas []*model.Idea) []*IdeaResponse {
	if len(ideas) == 0 {
		return []*IdeaResponse{}
	}

	// Collect IDs
	ideaIDs := make([]int64, len(ideas))
	userIDSet := make(map[int64]struct{})
	for i, idea := range ideas {
		ideaIDs[i] = idea.ID
		userIDSet[idea.InitiatorID] = struct{}{}
	}
	userIDs := make([]int64, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	// Batch queries (3 queries total instead of 3*N)
	contribCounts, _ := s.store.CountContributionsByIdeaIDs(ctx, ideaIDs)
	voterCounts, _ := s.store.CountVotersByIdeaIDs(ctx, ideaIDs)
	usersMap, _ := s.store.GetUsersByIDs(ctx, userIDs)

	result := make([]*IdeaResponse, len(ideas))
	for i, idea := range ideas {
		username := ""
		if u, ok := usersMap[idea.InitiatorID]; ok {
			username = u.Username
		}
		result[i] = &IdeaResponse{
			Idea:              idea,
			ContributionCount: contribCounts[idea.ID],
			VoterCount:        voterCounts[idea.ID],
			InitiatorUsername:  username,
		}
	}
	return result
}

// GetIdeaStats returns platform-level stats.
type PlatformStats struct {
	OpenIdeas         int `json:"open_ideas"`
	ClosedIdeas       int `json:"closed_ideas"`
	TotalContributions int `json:"total_contributions"`
}

// GetPlatformStats returns aggregate platform statistics with minimal queries.
func (s *Service) GetPlatformStats(ctx context.Context) (*PlatformStats, error) {
	_, openTotal, err := s.store.ListIdeas(ctx, "open", 1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to count open ideas: %w", err)
	}
	_, closedTotal, err := s.store.ListIdeas(ctx, "closed", 1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to count closed ideas: %w", err)
	}
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
	featuredCount, _ := s.store.CountFeaturedByUser(ctx, user.ID)

	// Sanitize avatar URL before returning to client
	user.AvatarURL = sanitizeAvatarURL(user.AvatarURL)

	return &UserProfile{
		User:              user,
		IdeaCount:         ideaCount,
		ContributionCount: contribCount,
		FeaturedCount:     featuredCount,
	}, nil
}
