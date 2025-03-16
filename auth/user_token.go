package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// UserTokenResponse represents the response from the Apple Music API when requesting a user token.
type UserTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// UserTokenManager manages user tokens for the Apple Music API.
type UserTokenManager struct {
	httpClient     *http.Client
	oauthConfig    *oauth2.Config
	developerToken *DeveloperToken
	tokenCache     TokenCache
}

// TokenCache interface for storing and retrieving user tokens.
type TokenCache interface {
	Get(userID string) (*oauth2.Token, error)
	Save(userID string, token *oauth2.Token) error
}

// MemoryTokenCache implements TokenCache in memory.
type MemoryTokenCache struct {
	tokens map[string]*oauth2.Token
}

// NewMemoryTokenCache creates a new MemoryTokenCache.
func NewMemoryTokenCache() *MemoryTokenCache {
	return &MemoryTokenCache{
		tokens: make(map[string]*oauth2.Token),
	}
}

// Get retrieves a token from the cache.
func (c *MemoryTokenCache) Get(userID string) (*oauth2.Token, error) {
	token, ok := c.tokens[userID]
	if !ok {
		return nil, fmt.Errorf("token not found for user %s", userID)
	}
	return token, nil
}

// Save stores a token in the cache.
func (c *MemoryTokenCache) Save(userID string, token *oauth2.Token) error {
	c.tokens[userID] = token
	return nil
}

// NewUserTokenManager creates a new UserTokenManager.
func NewUserTokenManager(developerToken *DeveloperToken, clientID, redirectURL string, cache TokenCache) *UserTokenManager {
	if cache == nil {
		cache = NewMemoryTokenCache()
	}

	oauthConfig := &oauth2.Config{
		ClientID:    clientID,
		RedirectURL: redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
		Scopes: []string{"musickit"},
	}

	return &UserTokenManager{
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		oauthConfig:    oauthConfig,
		developerToken: developerToken,
		tokenCache:     cache,
	}
}

// GetAuthURL returns the URL to redirect the user to for authorization.
func (m *UserTokenManager) GetAuthURL(state string) string {
	return m.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges an authorization code for a user token.
func (m *UserTokenManager) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return m.oauthConfig.Exchange(ctx, code)
}

// RefreshToken refreshes an expired user token.
func (m *UserTokenManager) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	source := m.oauthConfig.TokenSource(ctx, token)
	newToken, err := source.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return newToken, nil
}

// GetUserToken gets a user token from the cache or refreshes it if expired.
func (m *UserTokenManager) GetUserToken(ctx context.Context, userID string) (*oauth2.Token, error) {
	token, err := m.tokenCache.Get(userID)
	if err != nil {
		return nil, err
	}

	if token.Valid() {
		return token, nil
	}

	newToken, err := m.RefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	err = m.tokenCache.Save(userID, newToken)
	if err != nil {
		return nil, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return newToken, nil
}

// RequestUserToken requests a user token from the Apple Music API.
func (m *UserTokenManager) RequestUserToken(ctx context.Context, musicUserToken string) (*UserTokenResponse, error) {
	data := url.Values{}
	data.Set("music-user-token", musicUserToken)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.music.apple.com/v1/me/tokens",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+m.developerToken.String())

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user token: %s, status code: %d", string(body), resp.StatusCode)
	}

	var tokenResp UserTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

