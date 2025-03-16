# MusicKitKat - Apple Music SDK for Go

A comprehensive Go SDK for Apple Music that allows developers to integrate Apple Music functionality into their Go applications with minimal friction.

## Features

- Authentication with Apple Music API
- Catalog services (songs, albums, artists, music videos)
- Library services (user's library management)
- Playlist services (creation, modification, deletion)
- Search functionality with filtering options
- Recommendations and featured content
- Station and radio endpoints
- Streaming support

## Installation

```bash
go get github.com/user/musickitkat
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/user/musickitkat"
	"github.com/user/musickitkat/auth"
)

func main() {
	// Initialize the client with developer token
	developerToken, err := auth.NewDeveloperToken(
		"your-team-id",
		"your-key-id",
		[]byte("your-private-key"),
		"your-music-id",
	)
	if err != nil {
		log.Fatalf("Failed to create developer token: %v", err)
	}
	
	client := musickitkat.NewClient(
		musickitkat.WithDeveloperToken(developerToken),
	)
	
	// Search for songs
	ctx := context.Background()
	results, err := client.Catalog.Search(ctx, "The Beatles", musickitkat.SearchTypesSongs, nil)
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}
	
	// Print results
	for _, song := range results.Songs.Data {
		fmt.Printf("Song: %s by %s\n", song.Attributes.Name, song.Attributes.ArtistName)
	}
}
```

## Documentation

For detailed documentation, please visit [GoDoc](https://godoc.org/github.com/user/musickitkat).

## License

MIT