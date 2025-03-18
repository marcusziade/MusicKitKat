package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/user/musickitkat"
	"github.com/user/musickitkat/auth"
)

func main() {
	// Get the developer token components from environment variables
	teamID := os.Getenv("APPLE_TEAM_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKey := os.Getenv("APPLE_PRIVATE_KEY")

	// Create a new developer token with a 6-month expiration
	developerToken, err := auth.NewDeveloperToken(
		teamID,
		keyID,
		[]byte(privateKey),
		time.Hour*24*180,
	)
	if err != nil {
		log.Fatalf("Failed to create developer token: %v", err)
	}

	// Create a new client with the developer token
	client := musickitkat.NewClient(
		musickitkat.WithDeveloperToken(developerToken),
		musickitkat.WithLogLevel(musickitkat.LogLevelInfo),
	)

	// Set up context for the API requests
	ctx := context.Background()

	// Example 1: Get a song by ID and get the preview URL directly from the song object
	songID := "900032829" // Replace with a real song ID
	song, err := client.Catalog.GetSong(ctx, songID)
	if err != nil {
		log.Fatalf("Failed to get song: %v", err)
	}

	previewURL := song.GetPreviewURL()
	fmt.Printf("Song: %s by %s\n", song.Attributes.Name, song.Attributes.ArtistName)
	fmt.Printf("Preview URL: %s\n\n", previewURL)

	// Example 2: Use the helper method to directly get the preview URL for a song
	anotherSongID := "203709340" // Replace with a real song ID
	directPreviewURL, err := client.Catalog.GetSongPreviewURL(ctx, anotherSongID)
	if err != nil {
		log.Fatalf("Failed to get song preview URL: %v", err)
	}

	fmt.Printf("Direct Preview URL for song ID %s: %s\n", anotherSongID, directPreviewURL)

	// Example of how to handle a song without previews
	fmt.Println("\nHandling songs without previews:")
	noPreviewSongID := "invalid-id" // This will likely fail, just for demo
	noPreviewURL, err := client.Catalog.GetSongPreviewURL(ctx, noPreviewSongID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Preview URL: %s\n", noPreviewURL)
	}
}