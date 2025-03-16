// Package auth provides authentication functionality for the Apple Music API.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DeveloperToken represents an Apple Music developer token.
type DeveloperToken struct {
	token string
}

// DeveloperTokenConfig contains the necessary information to generate a developer token.
type DeveloperTokenConfig struct {
	TeamID     string
	KeyID      string
	PrivateKey []byte
	MusicID    string
	ExpiresAt  time.Time
}

// DefaultTokenExpiration is the default expiration time for developer tokens (6 months).
const DefaultTokenExpiration = 6 * 30 * 24 * time.Hour

// NewDeveloperToken creates a new developer token with the provided credentials.
func NewDeveloperToken(teamID, keyID string, privateKey []byte, musicID string) (*DeveloperToken, error) {
	return NewDeveloperTokenWithExpiry(teamID, keyID, privateKey, musicID, time.Now().Add(DefaultTokenExpiration))
}

// NewDeveloperTokenWithExpiry creates a new developer token with a custom expiration time.
func NewDeveloperTokenWithExpiry(teamID, keyID string, privateKey []byte, musicID string, expiresAt time.Time) (*DeveloperToken, error) {
	config := DeveloperTokenConfig{
		TeamID:     teamID,
		KeyID:      keyID,
		PrivateKey: privateKey,
		MusicID:    musicID,
		ExpiresAt:  expiresAt,
	}

	return NewDeveloperTokenFromConfig(config)
}

// NewDeveloperTokenFromConfig creates a new developer token from a configuration struct.
func NewDeveloperTokenFromConfig(config DeveloperTokenConfig) (*DeveloperToken, error) {
	key, err := jwt.ParseECPrivateKeyFromPEM(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": config.TeamID,
		"iat": now.Unix(),
		"exp": config.ExpiresAt.Unix(),
		"sub": config.MusicID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = config.KeyID

	signedToken, err := token.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &DeveloperToken{token: signedToken}, nil
}

// String returns the string representation of the developer token.
func (t *DeveloperToken) String() string {
	return t.token
}

// IsExpired checks if the token has expired.
func (t *DeveloperToken) IsExpired() (bool, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(t.token, jwt.MapClaims{})
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid expiration claim")
	}

	return time.Now().Unix() > int64(exp), nil
}

