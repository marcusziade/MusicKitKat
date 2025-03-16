package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/user/musickitkat/client"
	"github.com/user/musickitkat/models"
)

// SearchService provides search functionality for the Apple Music API.
type SearchService struct {
	BaseService
	storefront string
}

// NewSearchService creates a new SearchService with the provided client.
func NewSearchService(client *client.Client) *SearchService {
	return &SearchService{
		BaseService: *NewBaseService(client),
		storefront:  "us", // Default storefront
	}
}

// SetStorefront sets the default storefront for the search service.
func (s *SearchService) SetStorefront(storefront string) {
	s.storefront = storefront
}

// Search searches for resources in the catalog.
func (s *SearchService) Search(ctx context.Context, term string, types []string, options *models.SearchOptions) (*models.SearchResults, error) {
	if term == "" {
		return nil, fmt.Errorf("search term is required")
	}

	queryParams := url.Values{}
	queryParams.Set("term", term)

	if len(types) > 0 {
		queryParams.Set("types", commaSeparated(types))
	}

	if options != nil {
		if options.Limit > 0 {
			queryParams.Set("limit", fmt.Sprintf("%d", options.Limit))
		}

		if options.Offset > 0 {
			queryParams.Set("offset", fmt.Sprintf("%d", options.Offset))
		}

		if options.Storefront != "" {
			s.storefront = options.Storefront
		}

		if options.LanguageTag != "" {
			queryParams.Set("l", options.LanguageTag)
		}
	}

	path := s.buildPath(fmt.Sprintf("catalog/%s/search", s.storefront), queryParams)

	var response models.SearchResults
	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchHints gets search term hints for the provided term.
func (s *SearchService) SearchHints(ctx context.Context, term string) ([]string, error) {
	if term == "" {
		return nil, fmt.Errorf("search term is required")
	}

	queryParams := url.Values{}
	queryParams.Set("term", term)

	path := s.buildPath(fmt.Sprintf("catalog/%s/search/hints", s.storefront), queryParams)

	var response struct {
		Results struct {
			Terms []string `json:"terms"`
		} `json:"results"`
	}

	err := s.client.Get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	return response.Results.Terms, nil
}

