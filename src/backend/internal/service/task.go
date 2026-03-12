package service

import (
	"context"
	"fmt"

	"github.com/claway/server/internal/model"
)

// GetTask returns a task by ID.
func (s *Service) GetTask(ctx context.Context, id int64) (*model.Task, error) {
	task, err := s.store.GetTaskByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

// ListTasks returns all tasks for an idea.
func (s *Service) ListTasks(ctx context.Context, ideaID int64) ([]*model.Task, error) {
	tasks, err := s.store.ListTasksByIdeaID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	return tasks, nil
}

// ClaimTask allows a user to claim an open task.
// Validates that the task is open and its dependencies are approved.
func (s *Service) ClaimTask(ctx context.Context, taskID, userID int64) error {
	task, err := s.store.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.Status != model.TaskStatusOpen {
		return fmt.Errorf("task is not open for claiming (current status: %s)", task.Status)
	}

	// Verify dependency tasks are all approved before allowing claim
	if err := s.checkDependencies(ctx, task); err != nil {
		return fmt.Errorf("cannot claim task: %w", err)
	}

	if err := s.store.ClaimTask(ctx, taskID, userID); err != nil {
		return fmt.Errorf("failed to claim task: %w", err)
	}

	return nil
}

// UnclaimTask allows the claimer to release a claimed task back to open.
func (s *Service) UnclaimTask(ctx context.Context, taskID, userID int64) error {
	task, err := s.store.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.Status != model.TaskStatusClaimed {
		return fmt.Errorf("task is not in claimed status")
	}

	if !task.ClaimedBy.Valid || task.ClaimedBy.Int64 != userID {
		return fmt.Errorf("you are not the claimer of this task")
	}

	if err := s.store.UnclaimTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to unclaim task: %w", err)
	}

	return nil
}

// TokenUsageReport is the self-reported LLM usage for a task submission.
type TokenUsageReport struct {
	Model    string  `json:"model"`
	TokensIn int     `json:"tokens_in"`
	TokensOut int    `json:"tokens_out"`
	CostUSD  float64 `json:"cost_usd"`
}

// SubmitTaskRequest represents the request body for submitting a task.
type SubmitTaskRequest struct {
	Content    string            `json:"content"`
	Note       string            `json:"note"`
	TokenUsage *TokenUsageReport `json:"token_usage"`
}

// SubmitTask submits a task's output. Only the claimer can submit.
// Dependencies must be approved before submission.
func (s *Service) SubmitTask(ctx context.Context, taskID, userID int64, req SubmitTaskRequest) error {
	task, err := s.store.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.Status != model.TaskStatusClaimed && task.Status != model.TaskStatusRejected && task.Status != model.TaskStatusRevision {
		return fmt.Errorf("task cannot be submitted (current status: %s)", task.Status)
	}

	if !task.ClaimedBy.Valid || task.ClaimedBy.Int64 != userID {
		return fmt.Errorf("you are not the claimer of this task")
	}

	if req.Content == "" {
		return fmt.Errorf("content is required")
	}

	if len(req.Note) > 200 {
		return fmt.Errorf("note must be 200 characters or fewer")
	}

	// Verify dependencies are approved
	if err := s.checkDependencies(ctx, task); err != nil {
		return fmt.Errorf("cannot submit: %w", err)
	}

	// Update the task status to submitted
	if err := s.store.SubmitTask(ctx, taskID, req.Content, req.Note); err != nil {
		return fmt.Errorf("failed to submit task: %w", err)
	}

	// Update document content and create a new version
	doc, err := s.store.GetDocumentByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	newVersion := doc.CurrentVersion + 1
	diff := computeUnifiedDiff(doc.Content, req.Content)
	if err := s.store.UpdateDocumentContent(ctx, doc.ID, req.Content, newVersion); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	if _, err := s.store.CreateDocumentVersion(ctx, doc.ID, newVersion, req.Content, diff, userID); err != nil {
		return fmt.Errorf("failed to create document version: %w", err)
	}

	// Record self-reported token usage
	if req.TokenUsage != nil && req.TokenUsage.CostUSD > 0 {
		log := &model.TokenUsageLog{
			UserID:    userID,
			TaskID:    taskID,
			Model:     req.TokenUsage.Model,
			TokensIn:  req.TokenUsage.TokensIn,
			TokensOut: req.TokenUsage.TokensOut,
			CostUSD:   req.TokenUsage.CostUSD,
		}
		if err := s.store.CreateTokenUsageLog(ctx, log); err != nil {
			return fmt.Errorf("failed to record token usage: %w", err)
		}

		// Accumulate cost on the task
		if err := s.store.AccumulateTaskCost(ctx, taskID, req.TokenUsage.CostUSD); err != nil {
			return fmt.Errorf("failed to update task cost: %w", err)
		}
	}

	return nil
}

// ReviewTaskRequest represents the request body for reviewing a task.
type ReviewTaskRequest struct {
	Action       string  `json:"action"`        // "approve", "reject", or "revision"
	QualityScore float64 `json:"quality_score"` // 1.0, 1.2, or 1.5 (for approve)
	RejectReason string  `json:"reject_reason"` // required for reject
	Feedback     string  `json:"feedback"`      // required for revision
}

// ReviewTask allows the idea initiator to approve or reject a submitted task.
// On approval, creates a contribution record and awards credits.
func (s *Service) ReviewTask(ctx context.Context, taskID, userID int64, req ReviewTaskRequest) error {
	task, err := s.store.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.Status != model.TaskStatusSubmitted {
		return fmt.Errorf("task is not in submitted status")
	}

	// Verify the reviewer is the idea initiator
	idea, err := s.store.GetIdeaByID(ctx, task.IdeaID)
	if err != nil {
		return fmt.Errorf("failed to get idea: %w", err)
	}

	if idea.InitiatorID != userID {
		return fmt.Errorf("only the idea initiator can review tasks")
	}

	switch req.Action {
	case "approve":
		if req.QualityScore != 1.0 && req.QualityScore != 1.2 && req.QualityScore != 1.5 {
			return fmt.Errorf("quality_score must be 1.0, 1.2, or 1.5")
		}

		if err := s.store.ApproveTask(ctx, taskID, req.QualityScore); err != nil {
			return fmt.Errorf("failed to approve task: %w", err)
		}

		// Calculate credits: cost_usd * quality_score * 1000
		credits := task.CostUSDAccumulated * req.QualityScore * 1000
		weightedScore := task.CostUSDAccumulated * req.QualityScore

		// Create contribution record
		contrib := &model.Contribution{
			IdeaID:        task.IdeaID,
			TaskID:        task.ID,
			UserID:        task.ClaimedBy.Int64,
			CostUSD:       task.CostUSDAccumulated,
			QualityScore:  req.QualityScore,
			WeightedScore: weightedScore,
		}
		if err := s.store.CreateContribution(ctx, contrib); err != nil {
			return fmt.Errorf("failed to create contribution: %w", err)
		}

		// Award credits to the contributor
		if credits > 0 {
			if err := s.store.UpdateCreditsBalance(ctx, task.ClaimedBy.Int64, credits); err != nil {
				return fmt.Errorf("failed to update credits balance: %w", err)
			}

			creditTx := &model.CreditTransaction{
				UserID:        task.ClaimedBy.Int64,
				Type:          "earn_contribute",
				Amount:        credits,
				ReferenceType: "task",
				ReferenceID:   task.ID,
				Description:   fmt.Sprintf("Contribution reward for task %s (quality: %.1f)", task.Type, req.QualityScore),
			}
			if err := s.store.CreateCreditTransaction(ctx, creditTx); err != nil {
				return fmt.Errorf("failed to create credit transaction: %w", err)
			}
		}

	case "revision":
		if req.Feedback == "" {
			return fmt.Errorf("feedback is required when requesting revision")
		}

		if err := s.store.RevisionTask(ctx, taskID, req.Feedback); err != nil {
			return fmt.Errorf("failed to set task to revision: %w", err)
		}

	case "reject":
		if req.RejectReason == "" {
			return fmt.Errorf("reject_reason is required when rejecting")
		}

		if err := s.store.RejectTask(ctx, taskID, req.RejectReason); err != nil {
			return fmt.Errorf("failed to reject task: %w", err)
		}

	default:
		return fmt.Errorf("action must be 'approve', 'reject', or 'revision'")
	}

	return nil
}
