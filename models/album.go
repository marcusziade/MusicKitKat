package models

import "time"

// Album represents an album in the Apple Music API.
type Album struct {
	// Resource information
	Resource

	// Attributes of the album
	Attributes AlbumAttributes `json:"attributes,omitempty"`

	// Relationships of the album
	Relationships AlbumRelationships `json:"relationships,omitempty"`
}

// AlbumAttributes represents the attributes of an album.
type AlbumAttributes struct {
	// The artist name.
	ArtistName string `json:"artistName"`

	// The album artwork.
	Artwork Artwork `json:"artwork"`

	// The content rating.
	ContentRating string `json:"contentRating,omitempty"`

	// The copyright text.
	Copyright string `json:"copyright,omitempty"`

	// The editorial notes.
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`

	// The genre names.
	GenreNames []string `json:"genreNames"`

	// Whether the album is complete.
	IsComplete bool `json:"isComplete"`

	// Whether the album is a compilation.
	IsCompilation bool `json:"isCompilation"`

	// Whether the album is a single.
	IsSingle bool `json:"isSingle"`

	// The name of the album.
	Name string `json:"name"`

	// The play parameters.
	PlayParams PlayParameters `json:"playParams,omitempty"`

	// The record label.
	RecordLabel string `json:"recordLabel,omitempty"`

	// The release date.
	ReleaseDate string `json:"releaseDate"`

	// The track count.
	TrackCount int `json:"trackCount"`

	// The UPC code.
	UPC string `json:"upc,omitempty"`

	// The URL.
	URL string `json:"url"`
}

// AlbumRelationships represents the relationships of an album.
type AlbumRelationships struct {
	// The artists relationship.
	Artists Relationship `json:"artists,omitempty"`

	// The genres relationship.
	Genres Relationship `json:"genres,omitempty"`

	// The tracks relationship.
	Tracks Relationship `json:"tracks,omitempty"`

	// The record labels relationship.
	RecordLabels Relationship `json:"record-labels,omitempty"`
}

// AlbumsResponse represents a response containing albums.
type AlbumsResponse struct {
	// The albums data.
	Data []Album `json:"data"`

	// The response errors.
	Errors []interface{} `json:"errors,omitempty"`

	// The response meta.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// The next URL.
	Next string `json:"next,omitempty"`
}

// GetArtworkURL returns the URL for the album artwork with the specified dimensions.
func (a *Album) GetArtworkURL(width, height int) string {
	return a.Attributes.Artwork.URL
}

// FormatReleaseDate formats the release date as a time.Time.
func (a *Album) FormatReleaseDate() (time.Time, error) {
	return time.Parse("2006-01-02", a.Attributes.ReleaseDate)
}

