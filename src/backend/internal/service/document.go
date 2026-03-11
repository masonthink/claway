package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/claway/server/internal/model"
)

// GetDocument returns the current document for a task.
func (s *Service) GetDocument(ctx context.Context, taskID int64) (*model.Document, error) {
	doc, err := s.store.GetDocumentByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

// ListDocumentVersions returns all versions of a task's document.
func (s *Service) ListDocumentVersions(ctx context.Context, taskID int64) ([]*model.DocumentVersion, error) {
	doc, err := s.store.GetDocumentByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	versions, err := s.store.ListDocumentVersions(ctx, doc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}
	return versions, nil
}

// GetDocumentVersion returns a specific version of a task's document.
func (s *Service) GetDocumentVersion(ctx context.Context, taskID int64, version int) (*model.DocumentVersion, error) {
	doc, err := s.store.GetDocumentByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	ver, err := s.store.GetDocumentVersion(ctx, doc.ID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get document version: %w", err)
	}
	return ver, nil
}

// UpdateDocumentRequest represents a request to update a document.
type UpdateDocumentRequest struct {
	Content string `json:"content"`
}

// UpdateDocument updates a task's document content and creates a new version.
// Only the claimer of the task can update the document.
func (s *Service) UpdateDocument(ctx context.Context, taskID, userID int64, req UpdateDocumentRequest) error {
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}

	task, err := s.store.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Only the claimer can update the document
	if !task.ClaimedBy.Valid || task.ClaimedBy.Int64 != userID {
		return fmt.Errorf("you are not the claimer of this task")
	}

	if task.Status != model.TaskStatusClaimed && task.Status != model.TaskStatusRejected {
		return fmt.Errorf("task document cannot be updated in status: %s", task.Status)
	}

	doc, err := s.store.GetDocumentByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	newVersion := doc.CurrentVersion + 1
	diff := "" // TODO: compute diff

	if err := s.store.UpdateDocumentContent(ctx, doc.ID, req.Content, newVersion); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	if _, err := s.store.CreateDocumentVersion(ctx, doc.ID, newVersion, req.Content, diff, userID); err != nil {
		return fmt.Errorf("failed to create document version: %w", err)
	}

	return nil
}

// PublishPRD aggregates all approved task documents into a single PRD.
// Only the idea initiator can publish. All tasks must be approved.
func (s *Service) PublishPRD(ctx context.Context, ideaID, userID int64) (*model.PRD, error) {
	idea, err := s.store.GetIdeaByID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get idea: %w", err)
	}

	if idea.InitiatorID != userID {
		return nil, fmt.Errorf("only the idea initiator can publish the PRD")
	}

	tasks, err := s.store.ListTasksByIdeaID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	// Verify all tasks are approved and calculate total cost
	var totalCostUSD float64
	var sections []string

	for _, task := range tasks {
		if task.Status != model.TaskStatusApproved {
			return nil, fmt.Errorf("task %s (%s) is not approved yet", task.Type, task.Title)
		}
		totalCostUSD += task.CostUSDAccumulated

		doc, err := s.store.GetDocumentByTaskID(ctx, task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get document for task %s: %w", task.Type, err)
		}

		section := fmt.Sprintf("# %s - %s\n\n%s", task.Type, task.Title, doc.Content)
		sections = append(sections, section)
	}

	// Merge all sections into one PRD document
	content := fmt.Sprintf("# %s - 完整产品需求文档\n\n%s\n\n---\n\n%s",
		idea.Title,
		idea.Description,
		strings.Join(sections, "\n\n---\n\n"),
	)

	// Calculate price: total_cost_usd * 2 * 1000 credits (2x markup)
	priceCredits := totalCostUSD * 2 * 1000

	prd, err := s.store.CreatePRD(ctx, ideaID, content, priceCredits)
	if err != nil {
		return nil, fmt.Errorf("failed to create PRD: %w", err)
	}

	// Mark idea as completed
	if err := s.store.UpdateIdeaStatus(ctx, ideaID, string(model.IdeaStatusCompleted)); err != nil {
		return nil, fmt.Errorf("failed to update idea status: %w", err)
	}

	return prd, nil
}
