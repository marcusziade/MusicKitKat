package models

// Artist represents an artist in the Apple Music API.
type Artist struct {
	// Resource information
	Resource

	// Attributes of the artist
	Attributes ArtistAttributes `json:"attributes,omitempty"`

	// Relationships of the artist
	Relationships ArtistRelationships `json:"relationships,omitempty"`
}

// ArtistAttributes represents the attributes of an artist.
type ArtistAttributes struct {
	// The artist artwork.
	Artwork Artwork `json:"artwork,omitempty"`

	// The editorial notes.
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`

	// The genre names.
	GenreNames []string `json:"genreNames"`

	// The name of the artist.
	Name string `json:"name"`

	// The URL.
	URL string `json:"url"`
}

// ArtistRelationships represents the relationships of an artist.
type ArtistRelationships struct {
	// The albums relationship.
	Albums Relationship `json:"albums,omitempty"`

	// The genres relationship.
	Genres Relationship `json:"genres,omitempty"`

	// The music videos relationship.
	MusicVideos Relationship `json:"music-videos,omitempty"`

	// The playlists relationship.
	Playlists Relationship `json:"playlists,omitempty"`

	// The station relationship.
	Station Relationship `json:"station,omitempty"`
}

// ArtistsResponse represents a response containing artists.
type ArtistsResponse struct {
	// The artists data.
	Data []Artist `json:"data"`

	// The response errors.
	Errors []interface{} `json:"errors,omitempty"`

	// The response meta.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// The next URL.
	Next string `json:"next,omitempty"`
}

// GetArtworkURL returns the URL for the artist artwork with the specified dimensions.
func (a *Artist) GetArtworkURL(width, height int) string {
	return a.Attributes.Artwork.URL
}