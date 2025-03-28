package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/marcusziade/musickitkat/client"
)

// RadioService provides access to radio endpoints of the Apple Music API.
type RadioService struct {
	BaseService
	storefront string
}

// NewRadioService creates a new RadioService with the provided client.
func NewRadioService(client *client.Client) *RadioService {
	return &RadioService{
		BaseService: *NewBaseService(client),
		storefront:  "us", // Default storefront
	}
}

// SetStorefront sets the default storefront for the radio service.
func (s *RadioService) SetStorefront(storefront string) {
	s.storefront = storefront
}

// GetStations gets all radio stations.
func (s *RadioService) GetStations(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath(fmt.Sprintf("catalog/%s/stations", s.storefront), queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetStation gets a radio station by ID.
func (s *RadioService) GetStation(ctx context.Context, id string) (interface{}, error) {
	path := fmt.Sprintf("catalog/%s/stations/%s", s.storefront, id)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("station not found: %s", id)
	}

	return response.Data[0], nil
}

// GetFeaturedStations gets featured radio stations.
func (s *RadioService) GetFeaturedStations(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath(fmt.Sprintf("catalog/%s/stations/featured", s.storefront), queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetRecentStations gets recently played radio stations.
func (s *RadioService) GetRecentStations(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath("me/recent/stations", queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
