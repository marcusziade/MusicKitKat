package models

import "time"

// Song represents a song in the Apple Music API.
type Song struct {
	// Resource information
	Resource

	// Attributes of the song
	Attributes SongAttributes `json:"attributes,omitempty"`

	// Relationships of the song
	Relationships SongRelationships `json:"relationships,omitempty"`
}

// SongAttributes represents the attributes of a song.
type SongAttributes struct {
	// The album name.
	AlbumName string `json:"albumName"`

	// The artist name.
	ArtistName string `json:"artistName"`

	// The song artwork.
	Artwork Artwork `json:"artwork"`

	// Whether the song is a composer.
	Composer string `json:"composer,omitempty"`

	// The content rating.
	ContentRating string `json:"contentRating,omitempty"`

	// The disc number.
	DiscNumber int `json:"discNumber"`

	// The duration in milliseconds.
	DurationInMillis int64 `json:"durationInMillis"`

	// The editorial notes.
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`

	// The genre names.
	GenreNames []string `json:"genreNames"`

	// Whether the song has lyrics.
	HasLyrics bool `json:"hasLyrics"`

	// Whether the song is Apple Digital Master.
	IsAppleDigitalMaster bool `json:"isAppleDigitalMaster,omitempty"`

	// The ISRC code.
	ISRC string `json:"isrc,omitempty"`

	// The name of the song.
	Name string `json:"name"`

	// The play parameters.
	PlayParams PlayParameters `json:"playParams,omitempty"`

	// The previews.
	Previews []Preview `json:"previews"`

	// The release date.
	ReleaseDate string `json:"releaseDate"`

	// The track number.
	TrackNumber int `json:"trackNumber"`

	// The URL.
	URL string `json:"url"`
}

// SongRelationships represents the relationships of a song.
type SongRelationships struct {
	// The albums relationship.
	Albums Relationship `json:"albums,omitempty"`

	// The artists relationship.
	Artists Relationship `json:"artists,omitempty"`

	// The genres relationship.
	Genres Relationship `json:"genres,omitempty"`

	// The station relationship.
	Station Relationship `json:"station,omitempty"`

	// The composers relationship.
	Composers Relationship `json:"composers,omitempty"`

	// The music videos relationship.
	MusicVideos Relationship `json:"music-videos,omitempty"`
}

// SongsResponse represents a response containing songs.
type SongsResponse struct {
	// The songs data.
	Data []Song `json:"data"`

	// The response errors.
	Errors []interface{} `json:"errors,omitempty"`

	// The response meta.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// The next URL.
	Next string `json:"next,omitempty"`
}

// GetArtworkURL returns the URL for the song artwork with the specified dimensions.
func (s *Song) GetArtworkURL(width, height int) string {
	return s.Attributes.Artwork.URL
}

// GetPreviewURL returns the URL for the first playable preview of the song.
// Returns an empty string if no playable preview is available.
func (s *Song) GetPreviewURL() string {
	// Check if there are any previews
	if len(s.Attributes.Previews) == 0 {
		// Fallback to PlayParams preview URL if available
		if s.Attributes.PlayParams.PreviewURL != "" {
			return s.Attributes.PlayParams.PreviewURL
		}
		return ""
	}
	
	// Find the first playable preview
	for _, preview := range s.Attributes.Previews {
		if preview.Playable {
			return preview.URL
		}
	}
	
	// If no playable preview found, return the first preview URL
	return s.Attributes.Previews[0].URL
}

// FormatReleaseDate formats the release date as a time.Time.
func (s *Song) FormatReleaseDate() (time.Time, error) {
	return time.Parse("2006-01-02", s.Attributes.ReleaseDate)
}