package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/user/musickitkat/auth"
)

func main() {
	// Read credentials from environment variables
	teamID := os.Getenv("APPLE_TEAM_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	musicID := os.Getenv("APPLE_MUSIC_ID")
	clientID := os.Getenv("APPLE_CLIENT_ID")
	redirectURL := os.Getenv("APPLE_REDIRECT_URL")

	if teamID == "" || keyID == "" || privateKeyPath == "" || musicID == "" || clientID == "" || redirectURL == "" {
		log.Fatal("Missing required environment variables.")
	}

	// Read private key
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	// Create developer token
	developerToken, err := auth.NewDeveloperToken(teamID, keyID, privateKey, musicID)
	if err != nil {
		log.Fatalf("Failed to create developer token: %v", err)
	}

	// Create token manager
	tokenManager := auth.NewUserTokenManager(developerToken, clientID, redirectURL, nil)

	// Set up HTTP server for OAuth flow
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Generate a state value to prevent CSRF
		state := "random-state-value"

		// Get the authentication URL
		authURL := tokenManager.GetAuthURL(state)

		// Display the authentication link
		fmt.Fprintf(w, `
            <h1>Apple Music Authentication</h1>
            <p>Click the link below to authenticate with Apple Music:</p>
            <a href="%s">Authenticate</a>
        `, authURL)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization code from the query parameters
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		// Validate state to prevent CSRF
		if state != "random-state-value" {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Exchange the authorization code for a token
		token, err := tokenManager.ExchangeCode(context.Background(), code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to exchange code: %v", err), http.StatusInternalServerError)
			return
		}

		// Display success message with token details
		fmt.Fprintf(w, `
            <h1>Authentication Successful</h1>
            <p>Access Token: %s</p>
            <p>Token Type: %s</p>
            <p>Refresh Token: %s</p>
            <p>Expiry: %s</p>
        `, token.AccessToken, token.TokenType, token.RefreshToken, token.Expiry)
	})

	// Start the HTTP server
	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}