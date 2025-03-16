package models

import "time"

// Playlist represents a playlist in the Apple Music API.
type Playlist struct {
	// Resource information
	Resource

	// Attributes of the playlist
	Attributes PlaylistAttributes `json:"attributes,omitempty"`

	// Relationships of the playlist
	Relationships PlaylistRelationships `json:"relationships,omitempty"`
}

// PlaylistAttributes represents the attributes of a playlist.
type PlaylistAttributes struct {
	// The artwork.
	Artwork Artwork `json:"artwork,omitempty"`

	// The curator name.
	CuratorName string `json:"curatorName,omitempty"`

	// The description.
	Description EditorialNotes `json:"description,omitempty"`

	// Whether the playlist is a featured playlist.
	IsFeatured bool `json:"isFeatured,omitempty"`

	// The last modified date.
	LastModifiedDate string `json:"lastModifiedDate,omitempty"`

	// The name of the playlist.
	Name string `json:"name"`

	// The play parameters.
	PlayParams PlayParameters `json:"playParams,omitempty"`

	// The playlist type.
	PlaylistType string `json:"playlistType"`

	// The URL.
	URL string `json:"url"`

	// The track count.
	TrackCount int `json:"trackCount"`
}

// PlaylistRelationships represents the relationships of a playlist.
type PlaylistRelationships struct {
	// The curator relationship.
	Curator Relationship `json:"curator,omitempty"`

	// The tracks relationship.
	Tracks Relationship `json:"tracks,omitempty"`

	// The featured artists relationship.
	FeaturedArtists Relationship `json:"featured-artists,omitempty"`
}

// PlaylistsResponse represents a response containing playlists.
type PlaylistsResponse struct {
	// The playlists data.
	Data []Playlist `json:"data"`

	// The response errors.
	Errors []interface{} `json:"errors,omitempty"`

	// The response meta.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// The next URL.
	Next string `json:"next,omitempty"`
}

// GetArtworkURL returns the URL for the playlist artwork with the specified dimensions.
func (p *Playlist) GetArtworkURL(width, height int) string {
	return p.Attributes.Artwork.URL
}

// FormatLastModifiedDate formats the last modified date as a time.Time.
func (p *Playlist) FormatLastModifiedDate() (time.Time, error) {
	return time.Parse(time.RFC3339, p.Attributes.LastModifiedDate)
}