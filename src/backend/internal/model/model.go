package model

import (
	"time"
)

type PackageType string

const (
	PackageLight    PackageType = "light"
	PackageStandard PackageType = "standard"
)

type IdeaStatus string

const (
	IdeaStatusDraft      IdeaStatus = "draft"
	IdeaStatusActive     IdeaStatus = "active"
	IdeaStatusCompleted  IdeaStatus = "completed"
	IdeaStatusCancelled  IdeaStatus = "cancelled"
)

type TaskStatus string

const (
	TaskStatusOpen     TaskStatus = "open"
	TaskStatusClaimed  TaskStatus = "claimed"
	TaskStatusSubmitted TaskStatus = "submitted"
	TaskStatusApproved TaskStatus = "approved"
	TaskStatusRejected TaskStatus = "rejected"
)

type TaskType string

const (
	TaskTypeD1 TaskType = "D1" // 竞品分析
	TaskTypeD2 TaskType = "D2" // 用户画像
	TaskTypeD3 TaskType = "D3" // 产品需求文档 (PRD)
	TaskTypeD4 TaskType = "D4" // 技术可行性评估
)

type User struct {
	ID             int64          `json:"id"`
	OpenClawID     string         `json:"openclaw_id,omitempty"`
	Username       string         `json:"username"`
	DisplayName    string         `json:"display_name"`
	AvatarURL      string         `json:"avatar_url"`
	AgentAPIKey    NullString `json:"-"`
	CreditsBalance float64        `json:"credits_balance"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type OAuthAccount struct {
	ID               int64        `json:"id"`
	UserID           int64        `json:"user_id"`
	Provider         string       `json:"provider"`
	ProviderUserID   string       `json:"provider_user_id"`
	ProviderUsername  string       `json:"provider_username"`
	ProviderEmail    string       `json:"provider_email"`
	AccessToken      string       `json:"-"`
	RefreshToken     string       `json:"-"`
	TokenExpiresAt   NullTime `json:"-"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

type Idea struct {
	ID                 int64       `json:"id"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	TargetUserHint     string      `json:"target_user_hint"`
	ProblemDefinition  string      `json:"problem_definition"`
	InitiatorID        int64       `json:"initiator_id"`
	InitiatorCutPercent float64    `json:"initiator_cut_percent"`
	PackageType        PackageType `json:"package_type"`
	Status             IdeaStatus  `json:"status"`
	CreatedAt          time.Time   `json:"created_at"`
	Deadline           NullTime `json:"deadline"`
}

type Task struct {
	ID                 int64          `json:"id"`
	IdeaID             int64          `json:"idea_id"`
	Type               TaskType       `json:"type"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	AcceptanceCriteria string         `json:"acceptance_criteria"`
	Dependencies       string         `json:"dependencies"`
	TokenLimitHint     int            `json:"token_limit_hint"`
	Status             TaskStatus     `json:"status"`
	ClaimedBy          NullInt64  `json:"claimed_by"`
	ClaimedAt          NullTime   `json:"claimed_at"`
	SubmittedAt        NullTime   `json:"submitted_at"`
	ApprovedAt         NullTime   `json:"approved_at"`
	OutputContent      NullString `json:"output_content"`
	OutputNote         NullString `json:"output_note"`
	QualityScore       NullFloat64 `json:"quality_score"`
	RejectReason       NullString `json:"reject_reason"`
	CostUSDAccumulated float64        `json:"cost_usd_accumulated"`
}

type Document struct {
	ID             int64     `json:"id"`
	TaskID         int64     `json:"task_id"`
	Content        string    `json:"content"`
	CurrentVersion int       `json:"current_version"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DocumentVersion struct {
	ID               int64          `json:"id"`
	DocumentID       int64          `json:"document_id"`
	Version          int            `json:"version"`
	Content          string         `json:"content"`
	DiffFromPrevious NullString `json:"diff_from_previous"`
	CreatedAt        time.Time      `json:"created_at"`
	CreatedBy        int64          `json:"created_by"`
}

type TokenUsageLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TaskID    int64     `json:"task_id"`
	Model     string    `json:"model"`
	TokensIn  int       `json:"tokens_in"`
	TokensOut int       `json:"tokens_out"`
	CostUSD   float64   `json:"cost_usd"`
	Timestamp time.Time `json:"timestamp"`
}

type Contribution struct {
	ID            int64   `json:"id"`
	IdeaID        int64   `json:"idea_id"`
	TaskID        int64   `json:"task_id"`
	UserID        int64   `json:"user_id"`
	CostUSD       float64 `json:"cost_usd"`
	QualityScore  float64 `json:"quality_score"`
	WeightedScore float64 `json:"weighted_score"`
	WeightPercent float64 `json:"weight_percent"`
}

type CreditTransaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   int64     `json:"reference_id"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type PRD struct {
	ID           int64        `json:"id"`
	IdeaID       int64        `json:"idea_id"`
	Content      string       `json:"content"`
	PublishedAt  NullTime `json:"published_at"`
	PriceCredits float64      `json:"price_credits"`
	ReadCount    int          `json:"read_count"`
}
