// Package client provides the HTTP client for making requests to the Apple Music API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
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

// LogLevel defines the verbosity of client logging
type LogLevel int

const (
	// LogLevelNone disables logging
	LogLevelNone LogLevel = iota
	// LogLevelError logs only errors
	LogLevelError
	// LogLevelInfo logs request and response info
	LogLevelInfo
	// LogLevelDebug logs detailed request and response content
	LogLevelDebug
)

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

	// Logger instance
	logger *log.Logger

	// Log level
	logLevel LogLevel
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

// WithLogger sets a custom logger.
func WithLogger(logger *log.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithLogLevel sets the logging level.
func WithLogLevel(level LogLevel) ClientOption {
	return func(c *Client) {
		c.logLevel = level
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
		logger:     log.New(io.Discard, "", log.LstdFlags),
		logLevel:   LogLevelNone,
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

// SetLogLevel sets the logging level.
func (c *Client) SetLogLevel(level LogLevel) {
	c.logLevel = level
}

// SetLogger sets a custom logger.
func (c *Client) SetLogger(logger *log.Logger) {
	c.logger = logger
}

// log logs a message at the specified level.
func (c *Client) log(level LogLevel, format string, v ...interface{}) {
	if c.logLevel >= level {
		c.logger.Printf(format, v...)
	}
}

// logRequest logs an HTTP request.
func (c *Client) logRequest(req *http.Request) {
	if c.logLevel >= LogLevelInfo {
		c.logger.Printf("REQUEST: %s %s", req.Method, req.URL.String())
	}

	if c.logLevel >= LogLevelDebug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			c.log(LogLevelError, "Failed to dump request: %v", err)
			return
		}
		c.logger.Printf("REQUEST DUMP:\n%s", dump)
	}
}

// logResponse logs an HTTP response.
func (c *Client) logResponse(resp *http.Response) {
	if c.logLevel >= LogLevelInfo {
		c.logger.Printf("RESPONSE: %d %s", resp.StatusCode, resp.Status)
	}

	if c.logLevel >= LogLevelDebug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			c.log(LogLevelError, "Failed to dump response: %v", err)
			return
		}
		c.logger.Printf("RESPONSE DUMP:\n%s", dump)
	}
}

// buildURL builds the full URL for a request.
func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s/%s/%s", c.baseURL, c.apiVersion, path)
}

// NewRequest creates a new HTTP request.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.buildURL(path)
	c.log(LogLevelInfo, "Creating new request: %s %s", method, url)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			c.log(LogLevelError, "Failed to encode request body: %v", err)
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}

		// Log the request body
		if c.logLevel >= LogLevelDebug {
			rawBody, _ := json.Marshal(body)
			c.logger.Printf("REQUEST BODY: %s", string(rawBody))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		c.log(LogLevelError, "Failed to create request: %v", err)
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

	c.logRequest(req)

	return req, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.log(LogLevelInfo, "Sending request: %s %s", req.Method, req.URL.String())

	resp, err := c.client.Do(req)
	if err != nil {
		c.log(LogLevelError, "Failed to send request: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	c.logResponse(resp)

	// Check for API errors
	if resp.StatusCode >= 400 {
		c.log(LogLevelError, "API returned error status: %d %s", resp.StatusCode, resp.Status)

		// Log response headers which might contain useful info
		c.log(LogLevelDebug, "Error response headers:")
		for key, values := range resp.Header {
			for _, value := range values {
				c.log(LogLevelDebug, "  %s: %s", key, value)
			}
		}

		// Add specific guidance for authentication errors
		if resp.StatusCode == 401 {
			authHeader := resp.Request.Header.Get("Authorization")
			c.log(LogLevelError, "Authentication failed (401 Unauthorized)")
			
			// Check if developer token is present
			if authHeader == "" || authHeader == "Bearer " {
				c.log(LogLevelError, "Developer token is missing. Ensure you've set it with WithDeveloperToken()")
				return nil, fmt.Errorf("API authentication error (status 401): Developer token is missing or invalid. " +
					"Check your APPLE_TEAM_ID, APPLE_KEY_ID, APPLE_MUSIC_ID, and private key")
			}
			
			// Check if User-Token is needed for this endpoint but not provided
			path := resp.Request.URL.Path
			if (strings.Contains(path, "/me/") || strings.Contains(path, "/library/")) && 
			   resp.Request.Header.Get("Music-User-Token") == "" {
				c.log(LogLevelError, "Music-User-Token is required for this endpoint but is missing")
				return nil, fmt.Errorf("API authentication error (status 401): Music-User-Token is required for %s but is missing. " +
					"Use WithUserToken() to set the user token", path)
			}
		}

		apiErr, err := c.parseErrorResponse(resp)
		if err != nil {
			c.log(LogLevelError, "Failed to parse error response: %v", err)
			// Add the status code to the error to make it more informative
			return nil, fmt.Errorf("HTTP %d: failed to parse error response: %w",
				resp.StatusCode, err)
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
		c.log(LogLevelError, "Failed to read error response body: %v", err)
		return nil, fmt.Errorf("failed to read error response body: %w", err)
	}

	// Log the raw error response body
	c.log(LogLevelDebug, "Error response body: %s", string(body))

	// Check if the body is empty or too short to be valid JSON
	if len(body) == 0 {
		c.log(LogLevelError, "Error response body is empty")
		return fmt.Errorf("API error (status %d): empty response body", resp.StatusCode), nil
	}

	if len(bytes.TrimSpace(body)) == 0 {
		c.log(LogLevelError, "Error response body contains only whitespace")
		return fmt.Errorf("API error (status %d): whitespace-only response body", resp.StatusCode), nil
	}

	// Check if the body looks like JSON
	trimmed := bytes.TrimSpace(body)
	isJSON := (len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '['))

	if !isJSON {
		// Try to extract meaningful content if it's HTML or plain text
		contentSample := string(body)
		if len(contentSample) > 100 {
			contentSample = contentSample[:100] + "..."
		}
		c.log(LogLevelError, "Error response is not JSON: %s", contentSample)
		return fmt.Errorf("API error (status %d): non-JSON response: %s",
			resp.StatusCode, contentSample), nil
	}

	// Restore the response body
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// Try to unmarshal as standard API error
	err = json.Unmarshal(body, &apiErr)
	if err != nil {
		c.log(LogLevelError, "Failed to unmarshal error response: %v", err)
		c.log(LogLevelDebug, "Unmarshalling failed for body: %s", string(body))

		// Create a fallback error with the status code and raw body preview
		contentSample := string(body)
		if len(contentSample) > 100 {
			contentSample = contentSample[:100] + "..."
		}

		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, contentSample), nil
	}

	apiErr.StatusCode = resp.StatusCode
	c.log(LogLevelInfo, "Parsed API error: %+v", apiErr)

	return &apiErr, nil
}

// decodeJSONResponse decodes a JSON response into the provided result.
func (c *Client) decodeJSONResponse(resp *http.Response, result interface{}) error {
	// Save the response body for logging if needed
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log(LogLevelError, "Failed to read response body: %v", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the raw response body
	c.log(LogLevelDebug, "Response body: %s", string(body))

	// Restore the response body
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// Try to unmarshal the response
	if err := json.Unmarshal(body, result); err != nil {
		c.log(LogLevelError, "Failed to unmarshal response: %v", err)

		switch {
		case err.Error() == "unexpected end of JSON input":
			c.log(LogLevelError, "JSON is incomplete or empty")
		case err.Error() == "invalid character '\\'' looking for beginning of value":
			c.log(LogLevelError, "Response is not valid JSON, might be plain text or HTML")
		case err.Error() == "invalid character '<' looking for beginning of value":
			c.log(LogLevelError, "Response is likely HTML instead of JSON")
		}

		if err, ok := err.(*json.SyntaxError); ok {
			c.log(LogLevelError, "JSON syntax error at offset %d: %v", err.Offset, err)
			// Print the part of the JSON that caused the error
			if int(err.Offset) < len(body) {
				start := int(err.Offset) - 20
				if start < 0 {
					start = 0
				}
				end := int(err.Offset) + 20
				if end > len(body) {
					end = len(body)
				}
				c.log(LogLevelError, "Error context: ...%s...", string(body[start:end]))
			}
		}

		// Show expected structure of result
		resultType := fmt.Sprintf("%T", result)
		c.log(LogLevelDebug, "Expected to unmarshal into type: %s", resultType)

		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Get sends a GET request to the Apple Music API.
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	c.log(LogLevelInfo, "Making GET request to %s", path)

	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.decodeJSONResponse(resp, result)
}

// Post sends a POST request to the Apple Music API.
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	c.log(LogLevelInfo, "Making POST request to %s", path)

	req, err := c.NewRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.decodeJSONResponse(resp, result)
}

// Put sends a PUT request to the Apple Music API.
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	c.log(LogLevelInfo, "Making PUT request to %s", path)

	req, err := c.NewRequest(ctx, "PUT", path, body)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.decodeJSONResponse(resp, result)
}

// Delete sends a DELETE request to the Apple Music API.
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	c.log(LogLevelInfo, "Making DELETE request to %s", path)

	req, err := c.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.decodeJSONResponse(resp, result)
}

