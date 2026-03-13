package model

import (
	"encoding/json"
	"time"
)

// --- Status constants ---

type IdeaStatus string

const (
	IdeaStatusOpen      IdeaStatus = "open"
	IdeaStatusClosed    IdeaStatus = "closed"
	IdeaStatusCancelled IdeaStatus = "cancelled"
)

type ContributionStatus string

const (
	ContributionStatusDraft     ContributionStatus = "draft"
	ContributionStatusSubmitted ContributionStatus = "submitted"
)

// --- Core models ---

// User represents a platform user.
type User struct {
	ID          int64     `json:"id"`
	OpenClawID  string    `json:"openclaw_id,omitempty"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OAuthAccount represents a linked OAuth provider account.
type OAuthAccount struct {
	ID               int64    `json:"id"`
	UserID           int64    `json:"user_id"`
	Provider         string   `json:"provider"`
	ProviderUserID   string   `json:"provider_user_id"`
	ProviderUsername  string   `json:"provider_username"`
	ProviderEmail    string   `json:"provider_email"`
	AccessToken      string   `json:"-"`
	RefreshToken     string   `json:"-"`
	TokenExpiresAt   NullTime `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Idea represents a product idea open for competitive contributions.
type Idea struct {
	ID          int64      `json:"id"`
	InitiatorID int64      `json:"initiator_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	TargetUser  string     `json:"target_user"`
	CoreProblem string     `json:"core_problem"`
	OutOfScope  NullString `json:"out_of_scope"`
	Status      IdeaStatus `json:"status"`
	Deadline    time.Time  `json:"deadline"`
	RevealedAt  NullTime   `json:"revealed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Contribution represents a user's solution document for an idea.
type Contribution struct {
	ID          int64              `json:"id"`
	IdeaID      int64              `json:"idea_id"`
	AuthorID    int64              `json:"author_id"`
	Content     string             `json:"content"`
	DecisionLog json.RawMessage    `json:"decision_log"`
	Status      ContributionStatus `json:"status"`
	ViewCount   int                `json:"view_count"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	SubmittedAt NullTime           `json:"submitted_at"`
}

// Vote represents a user's vote for a contribution.
type Vote struct {
	ID             int64     `json:"id"`
	IdeaID         int64     `json:"idea_id"`
	VoterID        int64     `json:"voter_id"`
	ContributionID int64     `json:"contribution_id"`
	VotedAt        time.Time `json:"voted_at"`
}

// RateLimit tracks daily action counts per user.
type RateLimit struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Action     string    `json:"action"`
	ActionDate time.Time `json:"action_date"`
	Count      int       `json:"count"`
}

// RankedResult represents one entry in the reveal snapshot.
type RankedResult struct {
	ContributionID int64 `json:"contribution_id"`
	VoteCount      int   `json:"vote_count"`
	Rank           int   `json:"rank"`
	IsFeatured     bool  `json:"is_featured"`
}

// RevealSnapshot holds the frozen ranking data at reveal time.
type RevealSnapshot struct {
	ID            int64           `json:"id"`
	IdeaID        int64           `json:"idea_id"`
	RankedResults json.RawMessage `json:"ranked_results"`
	TotalVotes    int             `json:"total_votes"`
	RevealedAt    time.Time       `json:"revealed_at"`
}

// AuthSession represents a pending authentication session for agent-based flows.
// Sessions are stored in-memory (not in the database) because they are short-lived
// (5 minute expiry) and do not need to survive server restarts.
type AuthSession struct {
	ID        string    `json:"id"`
	Token     string    `json:"token,omitempty"` // empty until OAuth completes
	Status    string    `json:"status"`          // "pending" or "completed"
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
