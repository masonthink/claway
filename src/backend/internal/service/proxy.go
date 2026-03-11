package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/claway/server/internal/model"
)

// Model pricing per million tokens (USD).
var modelPricing = map[string]struct {
	InputPerM  float64
	OutputPerM float64
}{
	"claude-opus-4":     {InputPerM: 15.0, OutputPerM: 75.0},
	"claude-sonnet-4-5": {InputPerM: 3.0, OutputPerM: 15.0},
	"gpt-4o":            {InputPerM: 2.5, OutputPerM: 10.0},
	"gpt-4o-mini":       {InputPerM: 0.15, OutputPerM: 0.6},
}

// ChatRequest represents an OpenAI-compatible chat completion request.
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	// Pass through any additional fields
	Stream      *bool    `json:"stream,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
}

// ChatMessage is a single message in the chat request.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the upstream LLM response (simplified).
type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   *ChatUsage   `json:"usage"`
}

// ChatChoice is a single choice in the chat response.
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatUsage contains token usage information from the upstream response.
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ProxyChat forwards a chat request to the upstream LLM, records token usage,
// and returns the response. The task must be claimed by the current user.
func (s *Service) ProxyChat(ctx context.Context, userID, taskID int64, reqBody []byte) ([]byte, error) {
	// Verify task belongs to the user and is in progress
	task, err := s.store.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if !task.ClaimedBy.Valid || task.ClaimedBy.Int64 != userID {
		return nil, fmt.Errorf("task is not claimed by you")
	}

	if task.Status != model.TaskStatusClaimed {
		return nil, fmt.Errorf("task is not in progress (status: %s)", task.Status)
	}

	// Parse the request to extract model info
	var chatReq ChatRequest
	if err := json.Unmarshal(reqBody, &chatReq); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if chatReq.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	// Forward to upstream LLM
	upstreamURL := fmt.Sprintf("%s/v1/chat/completions", s.cfg.UpstreamLLMBaseURL)

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create upstream request: %w", err)
	}
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+s.cfg.UpstreamLLMAPIKey)

	upstreamResp, err := http.DefaultClient.Do(upstreamReq)
	if err != nil {
		return nil, fmt.Errorf("upstream LLM request failed: %w", err)
	}
	defer upstreamResp.Body.Close()

	respBody, err := io.ReadAll(upstreamResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read upstream response: %w", err)
	}

	// If upstream returned an error, pass it through
	if upstreamResp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("upstream returned status %d", upstreamResp.StatusCode)
	}

	// Parse response to extract usage
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		// Still return the response even if we can't parse usage
		return respBody, nil
	}

	// Record token usage if available
	if chatResp.Usage != nil {
		costUSD := calculateCost(chatReq.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)

		usageLog := &model.TokenUsageLog{
			UserID:    userID,
			TaskID:    taskID,
			Model:     chatReq.Model,
			TokensIn:  chatResp.Usage.PromptTokens,
			TokensOut: chatResp.Usage.CompletionTokens,
			CostUSD:   costUSD,
			Timestamp: time.Now(),
		}

		// Best-effort logging: don't fail the request if logging fails
		_ = s.store.CreateTokenUsageLog(ctx, usageLog)
		_ = s.store.AccumulateTaskCost(ctx, taskID, costUSD)
	}

	return respBody, nil
}

// calculateCost computes the USD cost based on model pricing.
func calculateCost(modelName string, tokensIn, tokensOut int) float64 {
	pricing, ok := modelPricing[modelName]
	if !ok {
		// Unknown model: use a conservative estimate
		pricing = modelPricing["gpt-4o"]
	}

	inputCost := float64(tokensIn) / 1_000_000 * pricing.InputPerM
	outputCost := float64(tokensOut) / 1_000_000 * pricing.OutputPerM

	return inputCost + outputCost
}
