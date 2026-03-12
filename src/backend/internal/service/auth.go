package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
)

// --- OAuth state management (in-memory, good enough for single-instance MVP) ---

type oauthState struct {
	CodeVerifier string
	CLIPort      string // empty for web flow
	CreatedAt    time.Time
}

var (
	oauthStates   = make(map[string]*oauthState)
	oauthStatesMu sync.Mutex
)

func saveOAuthState(state, codeVerifier, cliPort string) {
	oauthStatesMu.Lock()
	defer oauthStatesMu.Unlock()
	oauthStates[state] = &oauthState{
		CodeVerifier: codeVerifier,
		CLIPort:      cliPort,
		CreatedAt:    time.Now(),
	}
}

func popOAuthState(state string) (*oauthState, bool) {
	oauthStatesMu.Lock()
	defer oauthStatesMu.Unlock()
	s, ok := oauthStates[state]
	if ok {
		delete(oauthStates, state)
	}
	return s, ok
}

// --- X (Twitter) OAuth 2.0 ---

type xTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type xUserResponse struct {
	Data struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Username        string `json:"username"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

// GetXAuthURL generates the X OAuth authorization URL with PKCE.
func (s *Service) GetXAuthURL(cliPort string) (string, error) {
	if s.cfg.XClientID == "" {
		return "", fmt.Errorf("X_CLIENT_ID not configured")
	}

	// Generate PKCE code verifier and challenge
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return "", fmt.Errorf("generate code verifier: %w", err)
	}
	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
	challengeHash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(challengeHash[:])

	// Generate state
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	saveOAuthState(state, codeVerifier, cliPort)

	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {s.cfg.XClientID},
		"redirect_uri":          {s.cfg.XRedirectURI},
		"scope":                 {"tweet.read users.read offline.access"},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	return "https://twitter.com/i/oauth2/authorize?" + params.Encode(), nil
}

// HandleXCallback processes the X OAuth callback.
func (s *Service) HandleXCallback(ctx context.Context, code, state string) (redirectURL string, err error) {
	saved, ok := popOAuthState(state)
	if !ok {
		return "", fmt.Errorf("invalid or expired OAuth state")
	}

	// Exchange code for token
	tokenResp, err := s.exchangeXCode(code, saved.CodeVerifier)
	if err != nil {
		return "", err
	}

	// Fetch user profile
	profile, err := s.fetchXProfile(tokenResp.AccessToken)
	if err != nil {
		return "", err
	}

	// Find or create user
	user, err := s.findOrCreateOAuthUser(ctx, "x", profile.Data.ID, profile.Data.Username, profile.Data.Name, profile.Data.ProfileImageURL, tokenResp)
	if err != nil {
		return "", err
	}

	// Issue JWT
	jwtToken, err := s.issueJWT(user.ID)
	if err != nil {
		return "", err
	}

	// Redirect based on flow type
	if saved.CLIPort != "" {
		// CLI flow: redirect to localhost
		return fmt.Sprintf("http://127.0.0.1:%s/callback?token=%s", saved.CLIPort, url.QueryEscape(jwtToken)), nil
	}

	// Web flow: redirect to frontend
	return fmt.Sprintf("%s/auth/callback?token=%s", s.cfg.FrontendURL, url.QueryEscape(jwtToken)), nil
}

func (s *Service) exchangeXCode(code, codeVerifier string) (*xTokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {s.cfg.XRedirectURI},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequest("POST", "https://api.x.com/2/oauth2/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// X requires Basic Auth (client_id:client_secret) for confidential clients
	req.SetBasicAuth(s.cfg.XClientID, s.cfg.XClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp xTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &tokenResp, nil
}

func (s *Service) fetchXProfile(accessToken string) (*xUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.x.com/2/users/me?user.fields=profile_image_url", nil)
	if err != nil {
		return nil, fmt.Errorf("create profile request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("profile fetch failed (status %d): %s", resp.StatusCode, string(body))
	}

	var profile xUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	return &profile, nil
}

func (s *Service) findOrCreateOAuthUser(ctx context.Context, provider, providerUserID, username, displayName, avatarURL string, tokenResp *xTokenResponse) (*model.User, error) {
	// Check if OAuth account already exists
	oauthAccount, err := s.store.GetOAuthAccount(ctx, provider, providerUserID)
	if err == nil {
		// Existing account — update tokens and return user
		var expiresAt interface{}
		if tokenResp.ExpiresIn > 0 {
			t := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
			expiresAt = t
		}
		_ = s.store.UpdateOAuthTokens(ctx, oauthAccount.ID, tokenResp.AccessToken, tokenResp.RefreshToken, expiresAt)

		user, err := s.store.GetUserByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("get user for oauth account: %w", err)
		}
		return user, nil
	}

	if err != store.ErrNotFound {
		return nil, fmt.Errorf("lookup oauth account: %w", err)
	}

	// New user — create user + oauth account
	user, err := s.store.CreateUserFromOAuth(ctx, username, displayName, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	account := &model.OAuthAccount{
		UserID:           user.ID,
		Provider:         provider,
		ProviderUserID:   providerUserID,
		ProviderUsername:  username,
		ProviderEmail:    "",
		AccessToken:      tokenResp.AccessToken,
		RefreshToken:     tokenResp.RefreshToken,
	}
	if tokenResp.ExpiresIn > 0 {
		account.TokenExpiresAt = sql.NullTime{
			Time:  time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
			Valid: true,
		}
	}
	_, err = s.store.CreateOAuthAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("create oauth account: %w", err)
	}

	return user, nil
}

// --- Legacy OpenClaw OAuth (kept for backward compatibility) ---

type openClawTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type openClawUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthCallbackResponse is the response returned after successful OAuth callback.
type AuthCallbackResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *Service) HandleOpenClawCallback(ctx context.Context, code string) (*AuthCallbackResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

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

	user, err := s.store.GetUserByOpenClawID(ctx, profile.ID)
	if err == store.ErrNotFound {
		user, err = s.store.CreateUser(ctx, profile.ID, profile.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to look up user: %w", err)
	}

	token, err := s.issueJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to issue JWT: %w", err)
	}

	return &AuthCallbackResponse{
		Token: token,
		User:  user,
	}, nil
}

// --- JWT ---

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
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signed, nil
}
