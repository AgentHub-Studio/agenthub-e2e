package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP test client with bearer token auth.
type Client struct {
	base  string
	token string
	http  *http.Client
}

// NewClient creates a Client targeting baseURL with the given bearer token.
func NewClient(baseURL, token string) *Client {
	return &Client{
		base:  baseURL,
		token: token,
		http:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Get performs GET {path} and decodes the JSON response into dest.
func (c *Client) Get(path string, dest any) (int, error) {
	return c.do(http.MethodGet, path, nil, dest)
}

// Post performs POST {path} with body and decodes the JSON response into dest.
func (c *Client) Post(path string, body, dest any) (int, error) {
	return c.do(http.MethodPost, path, body, dest)
}

// Put performs PUT {path} with body and decodes the JSON response into dest.
func (c *Client) Put(path string, body, dest any) (int, error) {
	return c.do(http.MethodPut, path, body, dest)
}

// Delete performs DELETE {path}.
func (c *Client) Delete(path string) (int, error) {
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) do(method, path string, body, dest any) (int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.base+path, reqBody)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("execute request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if dest != nil && resp.StatusCode < 400 {
		if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
			return resp.StatusCode, fmt.Errorf("decode response: %w", err)
		}
	}
	return resp.StatusCode, nil
}

// Page mirrors the Go API pagination response shape.
type Page[T any] struct {
	Content       []T   `json:"content"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
	PageNumber    int   `json:"pageNumber"`
	PageSize      int   `json:"pageSize"`
}
