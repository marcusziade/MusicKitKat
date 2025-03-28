package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/marcusziade/musickitkat/client"
	"github.com/marcusziade/musickitkat/models"
)

// CatalogService provides access to the catalog endpoints of the Apple Music API.
type CatalogService struct {
	BaseService
	storefront string
}

// NewCatalogService creates a new CatalogService with the provided client.
func NewCatalogService(client *client.Client) *CatalogService {
	return &CatalogService{
		BaseService: *NewBaseService(client),
		storefront:  "us", // Default storefront
	}
}

// SetStorefront sets the default storefront for the catalog service.
func (s *CatalogService) SetStorefront(storefront string) {
	s.storefront = storefront
}

// GetSong gets a song by ID.
func (s *CatalogService) GetSong(ctx context.Context, id string) (*models.Song, error) {
	path := fmt.Sprintf("catalog/%s/songs/%s", s.storefront, id)

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

// GetSongs gets multiple songs by IDs.
func (s *CatalogService) GetSongs(ctx context.Context, ids []string) ([]models.Song, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one ID is required")
	}

	queryParams := url.Values{}
	queryParams.Set("ids", commaSeparated(ids))

	path := s.buildPath(fmt.Sprintf("catalog/%s/songs", s.storefront), queryParams)

	var response models.SongsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetAlbum gets an album by ID.
func (s *CatalogService) GetAlbum(ctx context.Context, id string) (*models.Album, error) {
	path := fmt.Sprintf("catalog/%s/albums/%s", s.storefront, id)

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

// GetAlbums gets multiple albums by IDs.
func (s *CatalogService) GetAlbums(ctx context.Context, ids []string) ([]models.Album, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one ID is required")
	}

	queryParams := url.Values{}
	queryParams.Set("ids", commaSeparated(ids))

	path := s.buildPath(fmt.Sprintf("catalog/%s/albums", s.storefront), queryParams)

	var response models.AlbumsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetArtist gets an artist by ID.
func (s *CatalogService) GetArtist(ctx context.Context, id string) (*models.Artist, error) {
	path := fmt.Sprintf("catalog/%s/artists/%s", s.storefront, id)

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

// GetArtists gets multiple artists by IDs.
func (s *CatalogService) GetArtists(ctx context.Context, ids []string) ([]models.Artist, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one ID is required")
	}

	queryParams := url.Values{}
	queryParams.Set("ids", commaSeparated(ids))

	path := s.buildPath(fmt.Sprintf("catalog/%s/artists", s.storefront), queryParams)

	var response models.ArtistsResponse
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetPlaylist gets a playlist by ID.
func (s *CatalogService) GetPlaylist(ctx context.Context, id string) (*models.Playlist, error) {
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

// GetPlaylists gets multiple playlists by IDs.
func (s *CatalogService) GetPlaylists(ctx context.Context, ids []string) ([]models.Playlist, error) {
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

// GetSongPreviewURL gets the preview URL for a song by ID.
func (s *CatalogService) GetSongPreviewURL(ctx context.Context, id string) (string, error) {
	song, err := s.GetSong(ctx, id)
	if err != nil {
		return "", err
	}

	previewURL := song.GetPreviewURL()
	if previewURL == "" {
		return "", fmt.Errorf("no preview available for song: %s", id)
	}

	return previewURL, nil
}

// joinWithDelimiter joins string slices with the specified delimiter.
func joinWithDelimiter(items []string, delimiter string) string {
	if len(items) == 0 {
		return ""
	}

	result := items[0]
	for i := 1; i < len(items); i++ {
		result += delimiter + items[i]
	}

	return result
}
