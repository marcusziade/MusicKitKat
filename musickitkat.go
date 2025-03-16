// Package musickitkat provides a comprehensive Go SDK for Apple Music.
// It allows developers to easily integrate Apple Music functionality into their Go applications.
package musickitkat

import (
	"net/http"
	"time"

	"github.com/user/musickitkat/auth"
	"github.com/user/musickitkat/client"
	"github.com/user/musickitkat/services"
)

// Version is the current version of the MusicKitKat SDK.
const Version = "0.1.0"

// Client is the main entry point for the MusicKitKat SDK.
// It provides access to all Apple Music API services.
type Client struct {
	// Base HTTP client for making API requests
	httpClient *client.Client

	// Authentication tokens
	DeveloperToken string
	UserToken      string

	// Services for interacting with different parts of the Apple Music API
	Catalog         *services.CatalogService
	Library         *services.LibraryService
	Playlists       *services.PlaylistService
	Search          *services.SearchService
	Recommendations *services.RecommendationService
	Radio           *services.RadioService
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient.SetHTTPClient(httpClient)
	}
}

// WithDeveloperToken sets the developer token.
func WithDeveloperToken(token *auth.DeveloperToken) ClientOption {
	return func(c *Client) {
		c.DeveloperToken = token.String()
	}
}

// WithUserToken sets the user token.
func WithUserToken(token string) ClientOption {
	return func(c *Client) {
		c.UserToken = token
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.SetTimeout(timeout)
	}
}

// NewClient creates a new MusicKitKat client with the provided options.
func NewClient(options ...ClientOption) *Client {
	httpClient := client.NewClient()

	c := &Client{
		httpClient: httpClient,
	}

	// Apply all client options
	for _, option := range options {
		option(c)
	}

	// Initialize services
	c.Catalog = services.NewCatalogService(c.httpClient)
	c.Library = services.NewLibraryService(c.httpClient)
	c.Playlists = services.NewPlaylistService(c.httpClient)
	c.Search = services.NewSearchService(c.httpClient)
	c.Recommendations = services.NewRecommendationService(c.httpClient)
	c.Radio = services.NewRadioService(c.httpClient)

	return c
}

// SearchTypes represents the types of resources that can be searched.
type SearchTypes string

const (
	SearchTypesSongs       SearchTypes = "songs"
	SearchTypesAlbums      SearchTypes = "albums"
	SearchTypesArtists     SearchTypes = "artists"
	SearchTypesPlaylists   SearchTypes = "playlists"
	SearchTypesMusicVideos SearchTypes = "music-videos"
	SearchTypesStations    SearchTypes = "stations"
)
