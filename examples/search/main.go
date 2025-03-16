package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/user/musickitkat"
	"github.com/user/musickitkat/auth"
)

func main() {
	// Read credentials from environment variables
	teamID := os.Getenv("APPLE_TEAM_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	musicID := os.Getenv("APPLE_MUSIC_ID")

	if teamID == "" || keyID == "" || privateKeyPath == "" || musicID == "" {
		log.Fatal("Missing required environment variables. Please set APPLE_TEAM_ID, APPLE_KEY_ID, APPLE_PRIVATE_KEY_PATH, and APPLE_MUSIC_ID.")
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

	// Initialize client
	client := musickitkat.NewClient(
		musickitkat.WithDeveloperToken(developerToken),
	)

	// Search for songs
	ctx := context.Background()
	results, err := client.Search.Search(ctx, "The Beatles", []string{string(musickitkat.SearchTypesSongs)}, nil)
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	// Print results
	fmt.Println("Search Results:")
	for _, song := range results.Results.Songs.Data {
		fmt.Printf("- %s by %s (Album: %s)\n", song.Attributes.Name, song.Attributes.ArtistName, song.Attributes.AlbumName)
	}
}

