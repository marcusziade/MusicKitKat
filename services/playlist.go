package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/user/musickitkat/client"
	"github.com/user/musickitkat/models"
)

// PlaylistService provides access to playlist-related endpoints of the Apple Music API.
type PlaylistService struct {
	BaseService
	storefront string
}

// NewPlaylistService creates a new PlaylistService with the provided client.
func NewPlaylistService(client *client.Client) *PlaylistService {
	return &PlaylistService{
		BaseService: *NewBaseService(client),
		storefront:  "us", // Default storefront
	}
}

// SetStorefront sets the default storefront for the playlist service.
func (s *PlaylistService) SetStorefront(storefront string) {
	s.storefront = storefront
}

// GetCatalogPlaylist gets a playlist from the catalog by ID.
func (s *PlaylistService) GetCatalogPlaylist(ctx context.Context, id string) (*models.Playlist, error) {
	path := fmt.Sprintf("catalog/%s/playlists/%s", s.storefront, id)

	var response models.PlaylistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("playlist not found: %s", id)
	}

	return &response.Data[0], nil
}

// GetCatalogPlaylists gets multiple playlists from the catalog by IDs.
func (s *PlaylistService) GetCatalogPlaylists(ctx context.Context, ids []string) ([]models.Playlist, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one ID is required")
	}

	queryParams := url.Values{}
	queryParams.Set("ids", commaSeparated(ids))

	path := s.buildPath(fmt.Sprintf("catalog/%s/playlists", s.storefront), queryParams)

	var response models.PlaylistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetCatalogPlaylistTracks gets the tracks in a playlist from the catalog.
func (s *PlaylistService) GetCatalogPlaylistTracks(ctx context.Context, id string) ([]models.Song, error) {
	path := fmt.Sprintf("catalog/%s/playlists/%s/tracks", s.storefront, id)

	var response models.SongsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetUserPlaylist gets a user's playlist by ID.
func (s *PlaylistService) GetUserPlaylist(ctx context.Context, id string) (*models.Playlist, error) {
	path := fmt.Sprintf("me/library/playlists/%s", id)

	var response models.PlaylistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("playlist not found: %s", id)
	}

	return &response.Data[0], nil
}

// GetUserPlaylists gets all playlists in the user's library.
func (s *PlaylistService) GetUserPlaylists(ctx context.Context) ([]models.Playlist, error) {
	path := "me/library/playlists"

	var response models.PlaylistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetUserPlaylistsWithOptions gets playlists in the user's library with the specified options.
func (s *PlaylistService) GetUserPlaylistsWithOptions(ctx context.Context, options models.QueryParameters) ([]models.Playlist, error) {
	path := "me/library/playlists"

	// Build query parameters from options
	queryParams := s.buildQueryParams(options)
	if len(queryParams) > 0 {
		path = s.buildPath(path, queryParams)
	}

	var response models.PlaylistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetUserPlaylistTracks gets the tracks in a user's playlist.
func (s *PlaylistService) GetUserPlaylistTracks(ctx context.Context, id string) ([]models.Song, error) {
	path := fmt.Sprintf("me/library/playlists/%s/tracks", id)

	var response models.SongsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// CreatePlaylist creates a new playlist in the user's library.
func (s *PlaylistService) CreatePlaylist(ctx context.Context, name, description string, trackIDs []string) (*models.Playlist, error) {
	if name == "" {
		return nil, fmt.Errorf("playlist name is required")
	}

	tracks := make([]map[string]interface{}, len(trackIDs))
	for i, id := range trackIDs {
		tracks[i] = map[string]interface{}{
			"id":   id,
			"type": "songs",
		}
	}

	requestBody := map[string]interface{}{
		"attributes": map[string]interface{}{
			"name":        name,
			"description": description,
		},
		"relationships": map[string]interface{}{
			"tracks": map[string]interface{}{
				"data": tracks,
			},
		},
	}

	path := "me/library/playlists"

	var response models.PlaylistsResponse
	err := s.client.Post(ctx, path, requestBody, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("failed to create playlist")
	}

	return &response.Data[0], nil
}

// AddTracksToPlaylist adds tracks to a user's playlist.
func (s *PlaylistService) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	if len(trackIDs) == 0 {
		return fmt.Errorf("at least one track ID is required")
	}

	tracks := make([]map[string]interface{}, len(trackIDs))
	for i, id := range trackIDs {
		tracks[i] = map[string]interface{}{
			"id":   id,
			"type": "songs",
		}
	}

	requestBody := map[string]interface{}{
		"data": tracks,
	}

	path := fmt.Sprintf("me/library/playlists/%s/tracks", playlistID)

	var response interface{}
	err := s.client.Post(ctx, path, requestBody, &response)
	if err != nil {
		return err
	}

	return nil
}

// RemoveTracksFromPlaylist removes tracks from a user's playlist.
func (s *PlaylistService) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIndices []int) error {
	if len(trackIndices) == 0 {
		return fmt.Errorf("at least one track index is required")
	}

	path := fmt.Sprintf("me/library/playlists/%s/tracks", playlistID)

	var response interface{}
	err := s.client.Delete(ctx, path, &response)
	if err != nil {
		return err
	}

	return nil
}
