package service

import (
	"context"
	"fmt"

	"github.com/clawbeach/server/internal/store"
)

// UserComputeResponse contains a user's total compute usage.
type UserComputeResponse struct {
	TotalCostUSD float64 `json:"total_cost_usd"`
}

// GetMyCompute returns the current user's total compute cost.
func (s *Service) GetMyCompute(ctx context.Context, userID int64) (*UserComputeResponse, error) {
	total, err := s.store.GetUserComputeTotal(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get compute total: %w", err)
	}
	return &UserComputeResponse{TotalCostUSD: total}, nil
}

// UserIdeaComputeResponse contains a user's compute cost for a specific idea.
type UserIdeaComputeResponse struct {
	IdeaID       int64   `json:"idea_id"`
	TotalCostUSD float64 `json:"total_cost_usd"`
}

// GetMyIdeaCompute returns the current user's compute cost for a specific idea.
func (s *Service) GetMyIdeaCompute(ctx context.Context, userID, ideaID int64) (*UserIdeaComputeResponse, error) {
	total, err := s.store.GetUserIdeaCompute(ctx, userID, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get idea compute: %w", err)
	}
	return &UserIdeaComputeResponse{IdeaID: ideaID, TotalCostUSD: total}, nil
}

// IdeaComputeResponse contains the compute breakdown for an idea.
type IdeaComputeResponse struct {
	IdeaID    int64                `json:"idea_id"`
	Breakdown []store.ComputeEntry `json:"breakdown"`
}

// GetIdeaCompute returns the compute breakdown for an idea (per contributor).
func (s *Service) GetIdeaCompute(ctx context.Context, ideaID int64) (*IdeaComputeResponse, error) {
	breakdown, err := s.store.GetIdeaComputeBreakdown(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get idea compute breakdown: %w", err)
	}
	return &IdeaComputeResponse{IdeaID: ideaID, Breakdown: breakdown}, nil
}

// TaskComputeResponse contains the compute breakdown for a task.
type TaskComputeResponse struct {
	TaskID    int64                `json:"task_id"`
	Breakdown []store.ComputeEntry `json:"breakdown"`
}

// GetTaskCompute returns the compute breakdown for a task.
func (s *Service) GetTaskCompute(ctx context.Context, taskID int64) (*TaskComputeResponse, error) {
	breakdown, err := s.store.GetTaskComputeBreakdown(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task compute breakdown: %w", err)
	}
	return &TaskComputeResponse{TaskID: taskID, Breakdown: breakdown}, nil
}

// PlatformComputeResponse contains platform-wide compute statistics.
type PlatformComputeResponse struct {
	TotalCostUSD float64 `json:"total_cost_usd"`
	TotalUsers   int     `json:"total_users"`
}

// GetPlatformCompute returns platform-wide compute statistics.
func (s *Service) GetPlatformCompute(ctx context.Context) (*PlatformComputeResponse, error) {
	total, users, err := s.store.GetPlatformComputeTotal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform compute: %w", err)
	}
	return &PlatformComputeResponse{TotalCostUSD: total, TotalUsers: users}, nil
}
