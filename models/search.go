package models

// SearchResults represents search results from the Apple Music API.
type SearchResults struct {
	// The response meta.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// The response results.
	Results SearchResultsData `json:"results"`
}

// SearchResultsData represents the data in search results.
type SearchResultsData struct {
	// The song results.
	Songs SongsResponse `json:"songs,omitempty"`

	// The album results.
	Albums AlbumsResponse `json:"albums,omitempty"`

	// The artist results.
	Artists ArtistsResponse `json:"artists,omitempty"`

	// The playlist results.
	Playlists PlaylistsResponse `json:"playlists,omitempty"`

	// The music video results.
	MusicVideos interface{} `json:"music-videos,omitempty"`

	// The station results.
	Stations interface{} `json:"stations,omitempty"`

	// The top results.
	TopResults interface{} `json:"top,omitempty"`
}

// SearchOptions represents options for search requests.
type SearchOptions struct {
	// The limit for each type.
	Limit int `json:"limit,omitempty"`

	// The offset for each type.
	Offset int `json:"offset,omitempty"`

	// The storefront.
	Storefront string `json:"storefront,omitempty"`

	// The language tag.
	LanguageTag string `json:"l,omitempty"`

	// The types to include.
	Types []string `json:"types,omitempty"`
}

// DefaultSearchLimit is the default limit for search results.
const DefaultSearchLimit = 25

// DefaultSearchOffset is the default offset for search results.
const DefaultSearchOffset = 0