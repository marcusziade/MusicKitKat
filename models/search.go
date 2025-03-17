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
	MusicVideos MusicVideosResponse `json:"music-videos,omitempty"`

	// The station results.
	Stations StationsResponse `json:"stations,omitempty"`

	// The top results.
	TopResults TopResultsResponse `json:"top,omitempty"`
	
	// The curator results.
	Curators CuratorsResponse `json:"curators,omitempty"`
	
	// The radio station results.
	RadioStations RadioStationsResponse `json:"radio-stations,omitempty"`
	
	// The apple curators results.
	AppleCurators AppleCuratorsResponse `json:"apple-curators,omitempty"`
	
	// The record label results.
	RecordLabels RecordLabelsResponse `json:"record-labels,omitempty"`
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
	
	// Fields to include for each resource type.
	Include []string `json:"include,omitempty"`
	
	// Fields to exclude for each resource type.
	Exclude []string `json:"exclude,omitempty"`
	
	// Relationships to return for each resource.
	// For example: artists, genres, stations.
	// Multiple relationship types can be comma-separated.
	Extend []string `json:"extend,omitempty"`
}

// DefaultSearchLimit is the default limit for search results.
const DefaultSearchLimit = 25

// DefaultSearchOffset is the default offset for search results.
const DefaultSearchOffset = 0

// MusicVideosResponse represents a music videos response.
type MusicVideosResponse struct {
	Data []MusicVideo `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// MusicVideo represents a music video.
type MusicVideo struct {
	Resource
	Attributes    MusicVideoAttributes    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// MusicVideoAttributes represents attributes of a music video.
type MusicVideoAttributes struct {
	ArtistName string `json:"artistName,omitempty"`
	Artwork    Artwork `json:"artwork,omitempty"`
	ContentRating string `json:"contentRating,omitempty"`
	DurationInMillis int64 `json:"durationInMillis,omitempty"`
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`
	GenreNames []string `json:"genreNames,omitempty"`
	ISRC string `json:"isrc,omitempty"`
	Name string `json:"name,omitempty"`
	PlayParams PlayParameters `json:"playParams,omitempty"`
	PreviewURL string `json:"previewUrl,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	TrackNumber int `json:"trackNumber,omitempty"`
	URL string `json:"url,omitempty"`
	VideoSubType string `json:"videoSubType,omitempty"`
}

// StationsResponse represents a stations response.
type StationsResponse struct {
	Data []Station `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// Station represents a station.
type Station struct {
	Resource
	Attributes    StationAttributes    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// StationAttributes represents attributes of a station.
type StationAttributes struct {
	Artwork     Artwork        `json:"artwork,omitempty"`
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`
	IsLive      bool           `json:"isLive,omitempty"`
	Name        string         `json:"name,omitempty"`
	PlayParams  PlayParameters `json:"playParams,omitempty"`
	URL         string         `json:"url,omitempty"`
}

// TopResultsResponse represents top results.
type TopResultsResponse struct {
	Data []Resource `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// CuratorsResponse represents a curators response.
type CuratorsResponse struct {
	Data []Curator `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// Curator represents a curator.
type Curator struct {
	Resource
	Attributes    CuratorAttributes    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// CuratorAttributes represents attributes of a curator.
type CuratorAttributes struct {
	Artwork       Artwork        `json:"artwork,omitempty"`
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
}

// RadioStationsResponse represents a radio stations response.
type RadioStationsResponse struct {
	Data []RadioStation `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// RadioStation represents a radio station.
type RadioStation struct {
	Resource
	Attributes    RadioStationAttributes    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// RadioStationAttributes represents attributes of a radio station.
type RadioStationAttributes struct {
	Artwork       Artwork        `json:"artwork,omitempty"`
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`
	IsLive        bool           `json:"isLive,omitempty"`
	Name          string         `json:"name,omitempty"`
	PlayParams    PlayParameters `json:"playParams,omitempty"`
	URL           string         `json:"url,omitempty"`
}

// AppleCuratorsResponse represents an apple curators response.
type AppleCuratorsResponse struct {
	Data []AppleCurator `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// AppleCurator represents an apple curator.
type AppleCurator struct {
	Resource
	Attributes    AppleCuratorAttributes    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// AppleCuratorAttributes represents attributes of an apple curator.
type AppleCuratorAttributes struct {
	Artwork       Artwork        `json:"artwork,omitempty"`
	EditorialNotes EditorialNotes `json:"editorialNotes,omitempty"`
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
}

// RecordLabelsResponse represents a record labels response.
type RecordLabelsResponse struct {
	Data []RecordLabel `json:"data,omitempty"`
	Href string `json:"href,omitempty"`
	Next string `json:"next,omitempty"`
}

// RecordLabel represents a record label.
type RecordLabel struct {
	Resource
	Attributes    RecordLabelAttributes    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// RecordLabelAttributes represents attributes of a record label.
type RecordLabelAttributes struct {
	Artwork       Artwork        `json:"artwork,omitempty"`
	Description   string         `json:"description,omitempty"`
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
}