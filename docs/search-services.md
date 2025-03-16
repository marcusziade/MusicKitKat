# Search Services

MusicKitKat provides a comprehensive search interface through the `SearchService`, allowing you to search for various types of resources in the Apple Music catalog.

## Basic Search

The most common usage is a simple search for a specific term:

```go
import (
    "context"
    "fmt"
    "github.com/user/musickitkat"
)

// Initialize client
client := musickitkat.NewClient(
    musickitkat.WithDeveloperToken(developerToken),
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

## Search Types

MusicKitKat supports searching for the following types of resources:

```go
// SearchTypes represents the types of resources that can be searched.
const (
    SearchTypesSongs       SearchTypes = "songs"
    SearchTypesAlbums      SearchTypes = "albums"
    SearchTypesArtists     SearchTypes = "artists"
    SearchTypesPlaylists   SearchTypes = "playlists"
    SearchTypesMusicVideos SearchTypes = "music-videos"
    SearchTypesStations    SearchTypes = "stations"
)
```

You can search for multiple types at once:

```go
// Search for songs and albums
types := []string{
    string(musickitkat.SearchTypesSongs),
    string(musickitkat.SearchTypesAlbums),
}

results, err := client.Search.Search(ctx, "The Beatles", types, nil)
```

## Search Options

You can customize your search with various options:

```go
import (
    "github.com/user/musickitkat/models"
)

// Create search options
options := &models.SearchOptions{
    Limit:       25,         // Number of results per type
    Offset:      0,          // Offset for pagination
    Storefront:  "us",       // Storefront to search in
    LanguageTag: "en-US",    // Language tag
}

results, err := client.Search.Search(ctx, "The Beatles", types, options)
```

## Accessing Search Results

The search results are organized by type:

```go
// Access song results
if len(results.Results.Songs.Data) > 0 {
    fmt.Println("Found songs:")
    for _, song := range results.Results.Songs.Data {
        fmt.Printf("- %s by %s\n", song.Attributes.Name, song.Attributes.ArtistName)
    }
}

// Access album results
if len(results.Results.Albums.Data) > 0 {
    fmt.Println("Found albums:")
    for _, album := range results.Results.Albums.Data {
        fmt.Printf("- %s by %s\n", album.Attributes.Name, album.Attributes.ArtistName)
    }
}

// Access artist results
if len(results.Results.Artists.Data) > 0 {
    fmt.Println("Found artists:")
    for _, artist := range results.Results.Artists.Data {
        fmt.Printf("- %s\n", artist.Attributes.Name)
    }
}

// Access playlist results
if len(results.Results.Playlists.Data) > 0 {
    fmt.Println("Found playlists:")
    for _, playlist := range results.Results.Playlists.Data {
        fmt.Printf("- %s\n", playlist.Attributes.Name)
    }
}
```

## Search Hints

You can also get search term hints for a partial search term:

```go
// Get search hints
hints, err := client.Search.SearchHints(ctx, "beat")
if err != nil {
    // Handle error
}

fmt.Println("Search hints:")
for _, hint := range hints {
    fmt.Printf("- %s\n", hint)
}
```

## Pagination

For search results with many matches, you can paginate through the results:

```go
// First page
options := &models.SearchOptions{
    Limit:  25,
    Offset: 0,
}
results, err := client.Search.Search(ctx, "pop", types, options)

// Next page
options.Offset += options.Limit
moreResults, err := client.Search.Search(ctx, "pop", types, options)
```

## Storefront Customization

You can customize the storefront for all searches:

```go
// Set default storefront
client.Search.SetStorefront("jp") // Japan storefront

// Make search in the default storefront
results, err := client.Search.Search(ctx, "BABYMETAL", types, nil)
```

Or specify a storefront just for a specific search:

```go
// Search in a specific storefront
options := &models.SearchOptions{
    Storefront: "fr", // France storefront
}
results, err := client.Search.Search(ctx, "Daft Punk", types, options)
```