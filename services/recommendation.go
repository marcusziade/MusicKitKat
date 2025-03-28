package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/marcusziade/musickitkat/client"
)

// RecommendationService provides access to recommendation endpoints of the Apple Music API.
type RecommendationService struct {
	BaseService
	storefront string
}

// NewRecommendationService creates a new RecommendationService with the provided client.
func NewRecommendationService(client *client.Client) *RecommendationService {
	return &RecommendationService{
		BaseService: *NewBaseService(client),
		storefront:  "us", // Default storefront
	}
}

// SetStorefront sets the default storefront for the recommendation service.
func (s *RecommendationService) SetStorefront(storefront string) {
	s.storefront = storefront
}

// GetRecommendations gets recommendations for the user.
func (s *RecommendationService) GetRecommendations(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath(fmt.Sprintf("me/recommendations"), queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetRecommendation gets a recommendation by ID.
func (s *RecommendationService) GetRecommendation(ctx context.Context, id string) (interface{}, error) {
	path := fmt.Sprintf("me/recommendations/%s", id)

	var response struct {
		Data interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetFeaturedPlaylists gets featured playlists.
func (s *RecommendationService) GetFeaturedPlaylists(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath(fmt.Sprintf("catalog/%s/playlists/featured", s.storefront), queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetPersonalRecommendations gets personal recommendations for the user.
func (s *RecommendationService) GetPersonalRecommendations(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath(fmt.Sprintf("me/recommendations/personal"), queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetCuratedPlaylists gets curated playlists.
func (s *RecommendationService) GetCuratedPlaylists(ctx context.Context, limit int) (interface{}, error) {
	queryParams := url.Values{}
	s.setLimit(limit, queryParams)

	path := s.buildPath(fmt.Sprintf("catalog/%s/playlists/curated", s.storefront), queryParams)

	var response struct {
		Data []interface{} `json:"data"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
