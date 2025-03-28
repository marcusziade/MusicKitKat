package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/marcusziade/musickitkat"
	"github.com/marcusziade/musickitkat/auth"
	"github.com/marcusziade/musickitkat/models"
)

func main() {
	// Determine if we're in "auth flow" or "playlist demo" mode
	if len(os.Args) > 1 && os.Args[1] == "playlists" {
		demoUserPlaylists()
		return
	}

	// Default to auth flow
	runAuthServer()
}

// runAuthServer starts a local web server to handle the OAuth flow
func runAuthServer() {
	// Read credentials from environment variables
	teamID := os.Getenv("APPLE_TEAM_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	musicID := os.Getenv("APPLE_MUSIC_ID")
	clientID := os.Getenv("APPLE_CLIENT_ID")
	redirectURL := os.Getenv("APPLE_REDIRECT_URL")

	missingVars := []string{}
	if teamID == "" {
		missingVars = append(missingVars, "APPLE_TEAM_ID")
	}
	if keyID == "" {
		missingVars = append(missingVars, "APPLE_KEY_ID")
	}
	if privateKeyPath == "" {
		missingVars = append(missingVars, "APPLE_PRIVATE_KEY_PATH")
	}
	if musicID == "" {
		missingVars = append(missingVars, "APPLE_MUSIC_ID")
	}
	if clientID == "" {
		missingVars = append(missingVars, "APPLE_CLIENT_ID")
	}
	if redirectURL == "" {
		missingVars = append(missingVars, "APPLE_REDIRECT_URL")
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %s\n\nRefer to docs/authentication.md for details on setting up Apple Music API credentials.",
			strings.Join(missingVars, ", "))
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

		// Display success message with token details and instructions
		fmt.Fprintf(w, `
            <h1>Authentication Successful</h1>
            <p>Your Apple Music User Token:</p>
            <pre>%s</pre>
            <p>Set this as the APPLE_USER_TOKEN environment variable to access user-specific Apple Music features.</p>
            <p>To see your playlists, run: <code>go run main.go playlists</code> after setting the token.</p>
        `, token.AccessToken)
	})

	// Start the HTTP server
	fmt.Println("Starting Apple Music OAuth server on http://localhost:8080")
	fmt.Println("Visit http://localhost:8080 in your browser to authorize your Apple Music account")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// demoUserPlaylists demonstrates fetching and displaying user playlists
func demoUserPlaylists() {
	// Read credentials from environment variables
	teamID := os.Getenv("APPLE_TEAM_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	musicID := os.Getenv("APPLE_MUSIC_ID")
	userToken := os.Getenv("APPLE_USER_TOKEN")

	// Check for required variables
	missingVars := []string{}
	if teamID == "" {
		missingVars = append(missingVars, "APPLE_TEAM_ID")
	}
	if keyID == "" {
		missingVars = append(missingVars, "APPLE_KEY_ID")
	}
	if privateKeyPath == "" {
		missingVars = append(missingVars, "APPLE_PRIVATE_KEY_PATH")
	}
	if musicID == "" {
		missingVars = append(missingVars, "APPLE_MUSIC_ID")
	}
	if userToken == "" {
		missingVars = append(missingVars, "APPLE_USER_TOKEN")
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %s\n\nRefer to docs/authentication.md for details on setting up Apple Music API credentials.",
			strings.Join(missingVars, ", "))
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

	// Initialize client with both developer and user tokens
	client := musickitkat.NewClient(
		musickitkat.WithDeveloperToken(developerToken),
		musickitkat.WithUserToken(userToken),
		musickitkat.WithLogLevel(musickitkat.LogLevelInfo),
	)

	// Create context
	ctx := context.Background()

	// Get user's playlists with pagination and relationship options
	options := models.QueryParameters{
		Limit:   25,
		Offset:  0,
		Include: []string{"tracks"},
	}

	fmt.Println("Fetching your Apple Music playlists...")
	playlists, err := client.Playlists.GetUserPlaylistsWithOptions(ctx, options)
	if err != nil {
		log.Fatalf("Failed to get user playlists: %v", err)
	}

	// Print playlists
	fmt.Printf("\nFound %d playlists in your Apple Music library:\n\n", len(playlists))
	for i, playlist := range playlists {
		fmt.Printf("%d. %s (%d tracks)\n", i+1, playlist.Attributes.Name, playlist.Attributes.TrackCount)

		// If we have relationships data, print a few tracks
		if playlist.Relationships.Tracks.Data != nil && len(playlist.Relationships.Tracks.Data) > 0 {
			fmt.Println("   Sample tracks:")
			// Show up to 3 tracks as a preview
			maxTracks := 3
			if len(playlist.Relationships.Tracks.Data) < maxTracks {
				maxTracks = len(playlist.Relationships.Tracks.Data)
			}

			// Get the actual songs for this playlist
			if maxTracks > 0 {
				songs, err := client.Playlists.GetUserPlaylistTracks(ctx, playlist.ID)
				if err == nil && len(songs) > 0 {
					for j := 0; j < maxTracks; j++ {
						if j < len(songs) {
							fmt.Printf("   - %s by %s\n", songs[j].Attributes.Name, songs[j].Attributes.ArtistName)
						}
					}
				}
			}
		}
		fmt.Println()
	}
}

