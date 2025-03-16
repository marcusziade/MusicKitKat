package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/user/musickitkat/client"
	"github.com/user/musickitkat/models"
)

// LibraryService provides access to the user's library endpoints of the Apple Music API.
type LibraryService struct {
	BaseService
}

// NewLibraryService creates a new LibraryService with the provided client.
func NewLibraryService(client *client.Client) *LibraryService {
	return &LibraryService{
		BaseService: *NewBaseService(client),
	}
}

// GetLibrarySongs gets songs from the user's library.
func (s *LibraryService) GetLibrarySongs(ctx context.Context, limit, offset int) ([]models.Song, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)
	s.setOffset(offset, queryParams)

	path := s.buildPath("me/library/songs", queryParams)

	var response models.SongsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetLibrarySong gets a song from the user's library by ID.
func (s *LibraryService) GetLibrarySong(ctx context.Context, id string) (*models.Song, error) {
	path := fmt.Sprintf("me/library/songs/%s", id)

	var response models.SongsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("song not found: %s", id)
	}

	return &response.Data[0], nil
}

// GetLibraryAlbums gets albums from the user's library.
func (s *LibraryService) GetLibraryAlbums(ctx context.Context, limit, offset int) ([]models.Album, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)
	s.setOffset(offset, queryParams)

	path := s.buildPath("me/library/albums", queryParams)

	var response models.AlbumsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetLibraryAlbum gets an album from the user's library by ID.
func (s *LibraryService) GetLibraryAlbum(ctx context.Context, id string) (*models.Album, error) {
	path := fmt.Sprintf("me/library/albums/%s", id)

	var response models.AlbumsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("album not found: %s", id)
	}

	return &response.Data[0], nil
}

// GetLibraryArtists gets artists from the user's library.
func (s *LibraryService) GetLibraryArtists(ctx context.Context, limit, offset int) ([]models.Artist, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)
	s.setOffset(offset, queryParams)

	path := s.buildPath("me/library/artists", queryParams)

	var response models.ArtistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetLibraryArtist gets an artist from the user's library by ID.
func (s *LibraryService) GetLibraryArtist(ctx context.Context, id string) (*models.Artist, error) {
	path := fmt.Sprintf("me/library/artists/%s", id)

	var response models.ArtistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("artist not found: %s", id)
	}

	return &response.Data[0], nil
}

// GetRecentlyAdded gets resources recently added to the user's library.
func (s *LibraryService) GetRecentlyAdded(ctx context.Context, limit, offset int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)
	s.setOffset(offset, queryParams)

	path := s.buildPath("me/library/recently-added", queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetHeavyRotation gets resources in the user's heavy rotation.
func (s *LibraryService) GetHeavyRotation(ctx context.Context, limit, offset int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)
	s.setOffset(offset, queryParams)

	path := s.buildPath("me/library/heavy-rotation", queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// AddToLibrary adds resources to the user's library.
func (s *LibraryService) AddToLibrary(ctx context.Context, ids []string, resourceType string) error {
	if len(ids) == 0 {
		return fmt.Errorf("at least one ID is required")
	}

	if resourceType == "" {
		return fmt.Errorf("resource type is required")
	}

	resources := make([]map[string]interface{}, len(ids))
	for i, id := range ids {
		resources[i] = map[string]interface{}{
			"id":   id,
			"type": resourceType,
		}
	}

	requestBody := map[string]interface{}{
		"data": resources,
	}

	path := "me/library"

	var response interface{}
	err := s.client.Post(ctx, path, requestBody, &response)
	if err != nil {
		return err
	}

	return nil
}

