package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/claway/server/internal/model"
)

// taskTemplate defines a deliverable template for auto-generation.
type taskTemplate struct {
	Type               model.TaskType
	Title              string
	Description        string
	AcceptanceCriteria string
	Dependencies       string
	TokenLimitHint     int
}

// lightTemplates are the task templates for the light package (D1-D5).
var lightTemplates = []taskTemplate{
	{
		Type:               model.TaskTypeD1,
		Title:              "竞品分析报告",
		Description:        "Research and analyze the competitive landscape. Identify at least 3 direct competitors and 2 indirect competitors, analyze their strengths and weaknesses, and identify differentiation opportunities.",
		AcceptanceCriteria: ">=3 direct competitors + >=2 indirect competitors analyzed, with differentiation space analysis",
		Dependencies:       "",
		TokenLimitHint:     80000,
	},
	{
		Type:               model.TaskTypeD2,
		Title:              "目标用户画像",
		Description:        "Define target user personas with core pain points and usage scenarios. Create 2-3 detailed user personas with narrative scenarios.",
		AcceptanceCriteria: "2-3 user personas, each with >=2 narrative scenarios and current solution limitations",
		Dependencies:       "",
		TokenLimitHint:     60000,
	},
	{
		Type:               model.TaskTypeD3,
		Title:              "用户旅程地图",
		Description:        "Create a product requirements document with user stories, acceptance criteria, and P0/P1 priority classification.",
		AcceptanceCriteria: "User story format, each feature has acceptance criteria, P0 features <=10",
		Dependencies:       "D1,D2",
		TokenLimitHint:     100000,
	},
	{
		Type:               model.TaskTypeD4,
		Title:              "功能需求文档",
		Description:        "Design the business model including revenue model, pricing strategy, and cold-start path.",
		AcceptanceCriteria: "Revenue model + pricing rationale + cold-start path",
		Dependencies:       "D1,D2",
		TokenLimitHint:     60000,
	},
	{
		Type:               model.TaskTypeD5,
		Title:              "信息架构图",
		Description:        "Define success metrics including north star metric, process metrics, and measurement plans.",
		AcceptanceCriteria: "1 north star metric + 3-5 process metrics, each with target value and measurement plan",
		Dependencies:       "D1,D2",
		TokenLimitHint:     40000,
	},
}

// standardExtraTemplates are the additional templates for standard package (D6-D9).
var standardExtraTemplates = []taskTemplate{
	{
		Type:               model.TaskTypeD6,
		Title:              "页面流程图",
		Description:        "Design information architecture including page structure, navigation system, and permission matrix.",
		AcceptanceCriteria: "Complete page list (routes + content) + navigation structure + permission matrix",
		Dependencies:       "D3",
		TokenLimitHint:     60000,
	},
	{
		Type:               model.TaskTypeD7,
		Title:              "交互设计规范",
		Description:        "Map core user flows covering all P0 features with normal paths and exception paths.",
		AcceptanceCriteria: "Covers all P0 features, normal path + >=2 exception paths, includes empty/error/loading states",
		Dependencies:       "D3",
		TokenLimitHint:     80000,
	},
	{
		Type:               model.TaskTypeD8,
		Title:              "视觉设计规范",
		Description:        "Create design specifications including component library selection, color palette, typography, and custom component list.",
		AcceptanceCriteria: "Component library selection + color palette + typography + spacing + custom component list",
		Dependencies:       "D3",
		TokenLimitHint:     60000,
	},
	{
		Type:               model.TaskTypeD9,
		Title:              "技术可行性评估",
		Description:        "Evaluate technical feasibility including technology stack recommendations, key risk points, and feasibility conclusions.",
		AcceptanceCriteria: "Tech stack recommendations (with rationale) + key risk points + clear feasible/infeasible conclusion",
		Dependencies:       "D3",
		TokenLimitHint:     60000,
	},
}

// CreateIdeaRequest represents the request body for creating an idea.
type CreateIdeaRequest struct {
	Title               string           `json:"title"`
	Description         string           `json:"description"`
	TargetUserHint      string           `json:"target_user_hint"`
	ProblemDefinition   string           `json:"problem_definition"`
	InitiatorCutPercent float64          `json:"initiator_cut_percent"`
	PackageType         model.PackageType `json:"package_type"`
}

// CreateIdea validates input, creates an idea, auto-generates tasks, and creates documents.
func (s *Service) CreateIdea(ctx context.Context, userID int64, req CreateIdeaRequest) (*model.Idea, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if req.PackageType != model.PackageLight && req.PackageType != model.PackageStandard {
		return nil, fmt.Errorf("package_type must be 'light' or 'standard'")
	}
	if req.InitiatorCutPercent < 10 || req.InitiatorCutPercent > 30 {
		return nil, fmt.Errorf("initiator_cut_percent must be between 10 and 30")
	}

	idea := &model.Idea{
		Title:               req.Title,
		Description:         req.Description,
		TargetUserHint:      req.TargetUserHint,
		ProblemDefinition:   req.ProblemDefinition,
		InitiatorID:         userID,
		InitiatorCutPercent: req.InitiatorCutPercent,
		PackageType:         req.PackageType,
		Status:              model.IdeaStatusActive,
	}

	idea, err := s.store.CreateIdea(ctx, idea)
	if err != nil {
		return nil, fmt.Errorf("failed to create idea: %w", err)
	}

	// Select task templates based on package type
	templates := make([]taskTemplate, len(lightTemplates))
	copy(templates, lightTemplates)
	if req.PackageType == model.PackageStandard {
		templates = append(templates, standardExtraTemplates...)
	}

	// Create tasks and documents for each template
	for _, tmpl := range templates {
		task := &model.Task{
			IdeaID:             idea.ID,
			Type:               tmpl.Type,
			Title:              tmpl.Title,
			Description:        tmpl.Description,
			AcceptanceCriteria: tmpl.AcceptanceCriteria,
			Dependencies:       tmpl.Dependencies,
			TokenLimitHint:     tmpl.TokenLimitHint,
			Status:             model.TaskStatusOpen,
		}

		task, err := s.store.CreateTask(ctx, task)
		if err != nil {
			return nil, fmt.Errorf("failed to create task %s: %w", tmpl.Type, err)
		}

		// Create an empty document for the task
		_, err = s.store.CreateDocument(ctx, task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create document for task %s: %w", tmpl.Type, err)
		}
	}

	return idea, nil
}

// ListIdeasResponse wraps a list of ideas with total count for pagination.
type ListIdeasResponse struct {
	Ideas []*model.Idea `json:"ideas"`
	Total int           `json:"total"`
}

// ListIdeas returns ideas filtered by status with pagination.
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

	return &ListIdeasResponse{Ideas: ideas, Total: total}, nil
}

// GetIdea returns an idea by ID.
func (s *Service) GetIdea(ctx context.Context, id int64) (*model.Idea, error) {
	idea, err := s.store.GetIdeaByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get idea: %w", err)
	}
	return idea, nil
}

// IdeaContextEntry represents one task's document output in the idea context.
type IdeaContextEntry struct {
	TaskID   int64          `json:"task_id"`
	TaskType model.TaskType `json:"task_type"`
	Title    string         `json:"title"`
	Status   model.TaskStatus `json:"status"`
	Content  string         `json:"content,omitempty"`
}

// IdeaContextResponse is the aggregated context of all task documents for an idea.
type IdeaContextResponse struct {
	IdeaID  int64              `json:"idea_id"`
	Entries []IdeaContextEntry `json:"entries"`
}

// GetIdeaContext aggregates all task documents into a context payload for agent consumption.
func (s *Service) GetIdeaContext(ctx context.Context, ideaID int64) (*IdeaContextResponse, error) {
	// Verify idea exists
	_, err := s.store.GetIdeaByID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get idea: %w", err)
	}

	tasks, err := s.store.ListTasksByIdeaID(ctx, ideaID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	resp := &IdeaContextResponse{
		IdeaID:  ideaID,
		Entries: make([]IdeaContextEntry, 0, len(tasks)),
	}

	for _, task := range tasks {
		entry := IdeaContextEntry{
			TaskID:   task.ID,
			TaskType: task.Type,
			Title:    task.Title,
			Status:   task.Status,
		}

		// Only include content from approved tasks
		if task.Status == model.TaskStatusApproved {
			doc, err := s.store.GetDocumentByTaskID(ctx, task.ID)
			if err == nil && doc != nil {
				entry.Content = doc.Content
			}
		}

		resp.Entries = append(resp.Entries, entry)
	}

	return resp, nil
}

// checkDependencies verifies that all dependency tasks for a given task are approved.
func (s *Service) checkDependencies(ctx context.Context, task *model.Task) error {
	if task.Dependencies == "" {
		return nil
	}

	deps := strings.Split(task.Dependencies, ",")
	tasks, err := s.store.ListTasksByIdeaID(ctx, task.IdeaID)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	depMap := make(map[string]bool, len(deps))
	for _, d := range deps {
		depMap[strings.TrimSpace(d)] = true
	}

	for _, t := range tasks {
		if depMap[string(t.Type)] {
			if t.Status != model.TaskStatusApproved {
				return fmt.Errorf("dependency task %s (%s) is not yet approved", t.Type, t.Title)
			}
		}
	}

	return nil
}
