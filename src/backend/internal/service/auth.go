package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/claway/server/internal/model"
)

// openClawTokenResponse is the response from OpenClaw token exchange.
type openClawTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// openClawUserResponse is the response from OpenClaw user profile endpoint.
type openClawUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthCallbackResponse is the response returned after successful OAuth callback.
type AuthCallbackResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// HandleOpenClawCallback exchanges the authorization code for an access token,
// fetches the user profile, creates or finds the user, and issues a session JWT.
func (s *Service) HandleOpenClawCallback(ctx context.Context, code string) (*AuthCallbackResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	// Exchange code for access token
	tokenURL := fmt.Sprintf("%s/oauth/token", s.cfg.OpenClawBaseURL)
	resp, err := http.PostForm(tokenURL, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {s.cfg.OpenClawClientID},
		"client_secret": {s.cfg.OpenClawClientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp openClawTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Fetch user profile
	profileURL := fmt.Sprintf("%s/api/v1/me", s.cfg.OpenClawBaseURL)
	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, profileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile request: %w", err)
	}
	profileReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}
	defer profileResp.Body.Close()

	if profileResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(profileResp.Body)
		return nil, fmt.Errorf("profile fetch failed (status %d): %s", profileResp.StatusCode, string(body))
	}

	var profile openClawUserResponse
	if err := json.NewDecoder(profileResp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile response: %w", err)
	}

	// Find or create user
	user, err := s.store.GetUserByOpenClawID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to look up user: %w", err)
	}

	if user == nil {
		user, err = s.store.CreateUser(ctx, profile.ID, profile.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Issue JWT
	token, err := s.issueJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to issue JWT: %w", err)
	}

	return &AuthCallbackResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetMe returns the current user's profile.
func (s *Service) GetMe(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// issueJWT creates a signed JWT token for the given user ID.
func (s *Service) issueJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signed, nil
}
