package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
)

// CreateContributionRequest represents the request body for creating a contribution draft.
type CreateContributionRequest struct {
	Content     string          `json:"content"`
	DecisionLog json.RawMessage `json:"decision_log"`
}

// UpdateContributionRequest represents the request body for updating a draft.
type UpdateContributionRequest struct {
	Content     string          `json:"content"`
	DecisionLog json.RawMessage `json:"decision_log,omitempty"`
}

// ContributionResponse wraps a contribution with optional author info.
type ContributionResponse struct {
	ID          int64              `json:"id"`
	IdeaID      int64              `json:"idea_id"`
	AuthorID    *int64             `json:"author_id,omitempty"`    // hidden pre-reveal
	AuthorName  string             `json:"author_name,omitempty"`  // hidden pre-reveal
	Content     string             `json:"content"`
	Preview     string             `json:"preview,omitempty"`
	DecisionLog json.RawMessage    `json:"decision_log,omitempty"`
	Status      model.ContributionStatus `json:"status"`
	ViewCount   int                `json:"view_count"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	SubmittedAt *string            `json:"submitted_at,omitempty"`
	PreviewURL  string             `json:"preview_url,omitempty"`
}

// CreateContribution creates a new draft contribution for an idea.
func (s *Service) CreateContribution(ctx context.Context, userID int64, ideaID int64, req CreateContributionRequest) (*ContributionResponse, error) {
	// Check idea exists and is open
	idea, err := s.store.GetIdeaByID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("idea not found")
	}
	if idea.Status != model.IdeaStatusOpen {
		return nil, fmt.Errorf("idea is not open for contributions")
	}

	// Check for existing contribution
	existing, err := s.store.GetContributionByIdeaAndAuthor(ctx, ideaID, userID)
	if err == nil {
		// User already has a contribution for this idea
		if existing.Status == model.ContributionStatusSubmitted {
			return nil, fmt.Errorf("you have already submitted a contribution for this idea")
		}
		// Return existing draft
		return s.toContributionResponse(ctx, existing, true, true), nil
	}
	if !errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing contribution: %w", err)
	}

	decisionLog := req.DecisionLog
	if decisionLog == nil {
		decisionLog = json.RawMessage("[]")
	}

	contrib := &model.Contribution{
		IdeaID:      ideaID,
		AuthorID:    userID,
		Content:     req.Content,
		DecisionLog: decisionLog,
	}

	created, err := s.store.CreateContribution(ctx, contrib)
	if err != nil {
		return nil, fmt.Errorf("failed to create contribution: %w", err)
	}

	return s.toContributionResponse(ctx, created, true, true), nil
}

// UpdateContribution updates a draft contribution.
func (s *Service) UpdateContribution(ctx context.Context, userID int64, contribID int64, req UpdateContributionRequest) (*ContributionResponse, error) {
	contrib, err := s.store.GetContributionByID(ctx, contribID)
	if err != nil {
		return nil, fmt.Errorf("contribution not found")
	}

	if contrib.AuthorID != userID {
		return nil, fmt.Errorf("not authorized to update this contribution")
	}
	if contrib.Status != model.ContributionStatusDraft {
		return nil, fmt.Errorf("only draft contributions can be updated")
	}

	// Check idea is still open
	idea, err := s.store.GetIdeaByID(ctx, contrib.IdeaID)
	if err != nil {
		return nil, fmt.Errorf("idea not found")
	}
	if idea.Status != model.IdeaStatusOpen {
		return nil, fmt.Errorf("idea is no longer open")
	}

	updated, err := s.store.UpdateContributionDraft(ctx, contribID, req.Content, req.DecisionLog)
	if err != nil {
		return nil, fmt.Errorf("failed to update contribution: %w", err)
	}

	return s.toContributionResponse(ctx, updated, true, true), nil
}

// SubmitContribution locks a draft contribution.
func (s *Service) SubmitContribution(ctx context.Context, userID int64, contribID int64) (*ContributionResponse, error) {
	contrib, err := s.store.GetContributionByID(ctx, contribID)
	if err != nil {
		return nil, fmt.Errorf("contribution not found")
	}

	if contrib.AuthorID != userID {
		return nil, fmt.Errorf("not authorized to submit this contribution")
	}
	if contrib.Status != model.ContributionStatusDraft {
		return nil, fmt.Errorf("contribution is already submitted")
	}

	// Check idea is still open
	idea, err := s.store.GetIdeaByID(ctx, contrib.IdeaID)
	if err != nil {
		return nil, fmt.Errorf("idea not found")
	}
	if idea.Status != model.IdeaStatusOpen {
		return nil, fmt.Errorf("idea is no longer open, cannot submit")
	}

	// Content must not be empty for submission
	if contrib.Content == "" {
		return nil, fmt.Errorf("cannot submit an empty contribution")
	}

	submitted, err := s.store.SubmitContribution(ctx, contribID)
	if err != nil {
		return nil, fmt.Errorf("failed to submit contribution: %w", err)
	}

	return s.toContributionResponse(ctx, submitted, true, true), nil
}

// GetContribution returns a single contribution.
// Access rules: drafts are only visible to the author; submitted contributions are public.
func (s *Service) GetContribution(ctx context.Context, userID int64, contribID int64) (*ContributionResponse, error) {
	contrib, err := s.store.GetContributionByID(ctx, contribID)
	if err != nil {
		return nil, fmt.Errorf("contribution not found")
	}

	isAuthor := contrib.AuthorID == userID
	if contrib.Status == model.ContributionStatusDraft && !isAuthor {
		return nil, fmt.Errorf("contribution not found")
	}

	// Determine if we should reveal author info
	idea, _ := s.store.GetIdeaByID(ctx, contrib.IdeaID)
	revealed := idea != nil && idea.Status != model.IdeaStatusOpen

	// Increment view count for submitted contributions
	if contrib.Status == model.ContributionStatusSubmitted {
		_ = s.store.IncrementContributionViewCount(ctx, contribID)
	}

	return s.toContributionResponse(ctx, contrib, revealed || isAuthor, isAuthor), nil
}

// ListContributions returns contributions for an idea.
// Pre-reveal: random order, anonymous. Post-reveal: ordered, with author info.
func (s *Service) ListContributions(ctx context.Context, ideaID int64) ([]*ContributionResponse, error) {
	idea, err := s.store.GetIdeaByID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("idea not found")
	}

	contributions, err := s.store.ListContributionsByIdea(ctx, ideaID, string(idea.Status))
	if err != nil {
		return nil, fmt.Errorf("failed to list contributions: %w", err)
	}

	revealed := idea.Status != model.IdeaStatusOpen
	result := make([]*ContributionResponse, 0, len(contributions))
	for _, c := range contributions {
		result = append(result, s.toContributionResponse(ctx, c, revealed, revealed))
	}
	return result, nil
}

// ListMyContributions returns all contributions by the authenticated user.
func (s *Service) ListMyContributions(ctx context.Context, userID int64, limit, offset int) ([]*ContributionResponse, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	contributions, total, err := s.store.ListContributionsByAuthor(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list my contributions: %w", err)
	}

	result := make([]*ContributionResponse, 0, len(contributions))
	for _, c := range contributions {
		result = append(result, s.toContributionResponse(ctx, c, true, true))
	}
	return result, total, nil
}

// GetDraftPreview returns a draft contribution for the preview page (author only).
func (s *Service) GetDraftPreview(ctx context.Context, userID int64, contribID int64) (*ContributionResponse, error) {
	contrib, err := s.store.GetContributionByID(ctx, contribID)
	if err != nil {
		return nil, fmt.Errorf("contribution not found")
	}

	if contrib.AuthorID != userID {
		return nil, fmt.Errorf("not authorized")
	}

	return s.toContributionResponse(ctx, contrib, true, true), nil
}

// toContributionResponse converts a model.Contribution to a ContributionResponse.
// showAuthor controls whether author info is included.
// showFull controls whether full content or just preview is returned.
func (s *Service) toContributionResponse(ctx context.Context, c *model.Contribution, showAuthor, showFull bool) *ContributionResponse {
	resp := &ContributionResponse{
		ID:        c.ID,
		IdeaID:    c.IdeaID,
		Status:    c.Status,
		ViewCount: c.ViewCount,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if showAuthor {
		resp.AuthorID = &c.AuthorID
		if user, err := s.store.GetUserByID(ctx, c.AuthorID); err == nil {
			resp.AuthorName = user.Username
		}
	}

	if showFull {
		resp.Content = c.Content
		resp.DecisionLog = c.DecisionLog
	} else {
		// Generate preview: first 200 chars
		preview := c.Content
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		resp.Preview = preview
	}

	if c.SubmittedAt.Valid {
		t := c.SubmittedAt.Time.Format("2006-01-02T15:04:05Z")
		resp.SubmittedAt = &t
	}

	if c.Status == model.ContributionStatusDraft {
		resp.PreviewURL = fmt.Sprintf("%s/draft/%d", s.cfg.FrontendURL, c.ID)
	}

	return resp
}
