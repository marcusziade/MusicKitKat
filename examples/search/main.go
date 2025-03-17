package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/user/musickitkat"
	"github.com/user/musickitkat/auth"
	"github.com/user/musickitkat/models"
)

func main() {
	// Read credentials from environment variables
	teamID := os.Getenv("APPLE_TEAM_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	musicID := os.Getenv("APPLE_MUSIC_ID")

	// Check each variable individually to provide more specific error messages
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

	// Initialize client with debug logging enabled
	client := musickitkat.NewClient(
		musickitkat.WithDeveloperToken(developerToken),
		musickitkat.WithLogLevel(musickitkat.LogLevelInfo),
	)

	// Get search query from command line arguments or use a default
	searchQuery := "The Beatles"
	if len(os.Args) > 1 {
		searchQuery = os.Args[1]
	}
	
	fmt.Printf("Searching Apple Music for: %s\n", searchQuery)
	fmt.Printf("Using developer token with KeyID: %s, TeamID: %s, MusicID: %s\n", 
		keyID, teamID, musicID)

	// Search for multiple resource types
	ctx := context.Background()
	types := []string{
		string(musickitkat.SearchTypesSongs),
		string(musickitkat.SearchTypesAlbums),
		string(musickitkat.SearchTypesArtists),
	}
	
	// Create search options with relationships included
	options := &models.SearchOptions{
		Limit: 5,
		Include: []string{"artists"},
	}
	
	results, err := client.Search.Search(ctx, searchQuery, types, options)
	if err != nil {
		// Enhanced error reporting
		log.Printf("Failed to search Apple Music API: %v", err)
		log.Printf("Debug: Verify your Apple Developer credentials at https://developer.apple.com/account")
		log.Printf("Debug: Ensure your Apple Music private key is valid and accessible")
		log.Fatalf("Debug: If the error persists, check Apple Music API status for service disruptions")
	}

	// Print results
	fmt.Println("Search Results:")
	
	// Print song results
	if len(results.Results.Songs.Data) > 0 {
		fmt.Println("\nSongs:")
		for _, song := range results.Results.Songs.Data {
			fmt.Printf("- %s by %s (Album: %s)\n", 
				song.Attributes.Name, 
				song.Attributes.ArtistName, 
				song.Attributes.AlbumName)
		}
	}
	
	// Print album results
	if len(results.Results.Albums.Data) > 0 {
		fmt.Println("\nAlbums:")
		for _, album := range results.Results.Albums.Data {
			fmt.Printf("- %s by %s\n", 
				album.Attributes.Name, 
				album.Attributes.ArtistName)
		}
	}
	
	// Print artist results
	if len(results.Results.Artists.Data) > 0 {
		fmt.Println("\nArtists:")
		for _, artist := range results.Results.Artists.Data {
			fmt.Printf("- %s\n", artist.Attributes.Name)
		}
	}
}

