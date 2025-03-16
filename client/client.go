// Package client provides the HTTP client for making requests to the Apple Music API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/musickitkat/errors"
)

// DefaultBaseURL is the default base URL for the Apple Music API.
const DefaultBaseURL = "https://api.music.apple.com"

// DefaultAPIVersion is the default API version.
const DefaultAPIVersion = "v1"

// DefaultUserAgent is the default User-Agent header value.
const DefaultUserAgent = "MusicKitKat/0.1.0"

// DefaultTimeout is the default request timeout.
const DefaultTimeout = 30 * time.Second

// Client is the HTTP client for making requests to the Apple Music API.
type Client struct {
	// HTTP client
	client *http.Client
	
	// Base URL for API requests
	baseURL string
	
	// API version
	apiVersion string
	
	// User-Agent header value
	userAgent string
	
	// Additional headers
	headers map[string]string
	
	// Developer token
	developerToken string
	
	// User token
	userToken string
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets the base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithAPIVersion sets the API version.
func WithAPIVersion(apiVersion string) ClientOption {
	return func(c *Client) {
		c.apiVersion = apiVersion
	}
}

// WithUserAgent sets the User-Agent header value.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithHeader adds a header to the client.
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// NewClient creates a new Client with the provided options.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		client:     &http.Client{Timeout: DefaultTimeout},
		baseURL:    DefaultBaseURL,
		apiVersion: DefaultAPIVersion,
		userAgent:  DefaultUserAgent,
		headers:    make(map[string]string),
	}
	
	// Apply all client options
	for _, option := range options {
		option(client)
	}
	
	return client
}

// SetHTTPClient sets the HTTP client.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// SetTimeout sets the request timeout.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// SetDeveloperToken sets the developer token.
func (c *Client) SetDeveloperToken(token string) {
	c.developerToken = token
}

// SetUserToken sets the user token.
func (c *Client) SetUserToken(token string) {
	c.userToken = token
}

// buildURL builds the full URL for a request.
func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s/%s/%s", c.baseURL, c.apiVersion, path)
}

// NewRequest creates a new HTTP request.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.buildURL(path)
	
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set default headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	// Set authentication headers
	if c.developerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.developerToken)
	}
	
	if c.userToken != "" {
		req.Header.Set("Music-User-Token", c.userToken)
	}
	
	// Set additional headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	
	return req, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	// Check for API errors
	if resp.StatusCode >= 400 {
		apiErr, err := c.parseErrorResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse error response: %w", err)
		}
		return nil, apiErr
	}
	
	return resp, nil
}

// parseErrorResponse parses an error response from the Apple Music API.
func (c *Client) parseErrorResponse(resp *http.Response) (error, error) {
	var apiErr errors.APIError
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read error response body: %w", err)
	}
	
	// Restore the response body
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	
	err = json.Unmarshal(body, &apiErr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal error response: %w", err)
	}
	
	apiErr.StatusCode = resp.StatusCode
	
	return &apiErr, nil
}

// Get sends a GET request to the Apple Music API.
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return json.NewDecoder(resp.Body).Decode(result)
}

// Post sends a POST request to the Apple Music API.
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	req, err := c.NewRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}
	
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return json.NewDecoder(resp.Body).Decode(result)
}

// Put sends a PUT request to the Apple Music API.
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	req, err := c.NewRequest(ctx, "PUT", path, body)
	if err != nil {
		return err
	}
	
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return json.NewDecoder(resp.Body).Decode(result)
}

// Delete sends a DELETE request to the Apple Music API.
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	req, err := c.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return json.NewDecoder(resp.Body).Decode(result)
}