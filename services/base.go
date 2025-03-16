// Package services provides service implementations for the Apple Music API.
package services

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/user/musickitkat/client"
	"github.com/user/musickitkat/models"
)

// BaseService is the base service for all API services.
type BaseService struct {
	client *client.Client
}

// NewBaseService creates a new BaseService with the provided client.
func NewBaseService(client *client.Client) *BaseService {
	return &BaseService{
		client: client,
	}
}

// buildPath builds a path with query parameters.
func (s *BaseService) buildPath(path string, queryParams url.Values) string {
	if len(queryParams) == 0 {
		return path
	}

	return fmt.Sprintf("%s?%s", path, queryParams.Encode())
}

// buildQueryParams builds query parameters from a QueryParameters struct.
func (s *BaseService) buildQueryParams(params models.QueryParameters) url.Values {
	queryParams := url.Values{}

	if params.Limit > 0 {
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Offset > 0 {
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}

	if len(params.Include) > 0 {
		queryParams.Set("include", strings.Join(params.Include, ","))
	}

	if len(params.Exclude) > 0 {
		queryParams.Set("exclude", strings.Join(params.Exclude, ","))
	}

	if params.LanguageTag != "" {
		queryParams.Set("l", params.LanguageTag)
	}

	if params.Storefront != "" {
		queryParams.Set("storefront", params.Storefront)
	}

	return queryParams
}

// setLimit sets the limit query parameter.
func (s *BaseService) setLimit(limit int, queryParams url.Values) {
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
}

// setOffset sets the offset query parameter.
func (s *BaseService) setOffset(offset int, queryParams url.Values) {
	if offset > 0 {
		queryParams.Set("offset", strconv.Itoa(offset))
	}
}

// setInclude sets the include query parameter.
func (s *BaseService) setInclude(include []string, queryParams url.Values) {
	if len(include) > 0 {
		queryParams.Set("include", strings.Join(include, ","))
	}
}

// setTypes sets the types query parameter.
func (s *BaseService) setTypes(types []string, queryParams url.Values) {
	if len(types) > 0 {
		queryParams.Set("types", strings.Join(types, ","))
	}
}

// setTerm sets the term query parameter.
func (s *BaseService) setTerm(term string, queryParams url.Values) {
	if term != "" {
		queryParams.Set("term", term)
	}
}

