# MusicKitKat Quick Start Guide

This guide will help you get started with the MusicKitKat SDK for Go.

## Installation

```bash
go get github.com/user/musickitkat
```

## Prerequisites

To use the Apple Music API, you need:

1. An Apple Developer account
2. Membership in the Apple Developer Program
3. Access to MusicKit
4. A private key for creating developer tokens

## Authentication

### Developer Token

The first step is to create a developer token using your Apple Developer credentials:

```go
import (
    "github.com/user/musickitkat/auth"
)

// Read private key from file
privateKey, err := os.ReadFile("path/to/private/key.p8")
if err != nil {
    // Handle error
}

// Create developer token
developerToken, err := auth.NewDeveloperToken(
    "your-team-id",
    "your-key-id",
    privateKey,
    "your-music-id"
)
if err != nil {
    // Handle error
}
```

### Client Initialization

Initialize the MusicKitKat client with your developer token:

```go
import (
    "github.com/user/musickitkat"
)

// Initialize client with developer token
client := musickitkat.NewClient(
    musickitkat.WithDeveloperToken(developerToken),
)
```

## Basic Usage

### Searching for Music

```go
import (
    "context"
    "fmt"
)

// Search for songs
ctx := context.Background()
results, err := client.Search.Search(ctx, "The Beatles", []string{string(musickitkat.SearchTypesSongs)}, nil)
if err != nil {
    // Handle error
}

// Print results
for _, song := range results.Results.Songs.Data {
    fmt.Printf("Song: %s by %s\n", song.Attributes.Name, song.Attributes.ArtistName)
}
```

### Getting Album Details

```go
// Get album by ID
album, err := client.Catalog.GetAlbum(ctx, "album-id")
if err != nil {
    // Handle error
}

fmt.Printf("Album: %s by %s\n", album.Attributes.Name, album.Attributes.ArtistName)
```

### User Library Access

To access a user's library, you need a user token:

```go
// Set user token on client
client.UserToken = "user-token"

// Get user's playlists
playlists, err := client.Library.GetUserPlaylists(ctx)
if err != nil {
    // Handle error
}

for _, playlist := range playlists {
    fmt.Printf("Playlist: %s\n", playlist.Attributes.Name)
}
```

## Next Steps

- See the [examples](../examples) directory for more examples
- Refer to the [API documentation](./README.md) for detailed information on available services and methods