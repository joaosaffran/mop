package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps the standard http.Client with common configurations
type Client struct {
	client  *http.Client
	headers map[string]string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithHeader adds a default header to all requests
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithBearerToken sets the Authorization header with a Bearer token
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.headers["Authorization"] = "Bearer " + token
	}
}

// NewClient creates a new HTTP client with the given options
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}

	// Set default headers
	c.headers["Content-Type"] = "application/json"

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Body    []byte
	Headers map[string]string
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Do performs an HTTP request and returns the response
func (c *Client) Do(req Request) (*Response, error) {
	var bodyReader io.Reader
	if req.Body != nil {
		bodyReader = bytes.NewBuffer(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for key, value := range c.headers {
		httpReq.Header.Set(key, value)
	}

	// Set request-specific headers (override defaults)
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}, nil
}

// Post performs a POST request
func (c *Client) Post(url string, body []byte) (*Response, error) {
	return c.Do(Request{
		Method: http.MethodPost,
		URL:    url,
		Body:   body,
	})
}

// Get performs a GET request
func (c *Client) Get(url string) (*Response, error) {
	return c.Do(Request{
		Method: http.MethodGet,
		URL:    url,
	})
}
