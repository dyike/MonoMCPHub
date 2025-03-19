package unsplash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	models "github.com/dyike/MonoMCPHub/repo/models/unsplash"
)

const (
	UnsplashAPIEndpoint = "https://api.unsplash.com"
)

// UnsplashClient is a client for the Unsplash API
type UnsplashClient struct {
	apiKey    string
	baseURL   string
	client    *http.Client
	rateLimit struct {
		remaining     int
		limit         int
		resetTime     time.Time
		lastCheckTime time.Time
	}
}

type UnsplashConfig struct {
	AccessKey string
	Timeout   time.Duration
}

// NewUnsplashClient creates a new Unsplash API client
func NewUnsplashClient(cfg *UnsplashConfig) *UnsplashClient {
	return &UnsplashClient{
		apiKey:  cfg.AccessKey,
		baseURL: UnsplashAPIEndpoint,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// SearchPhotos searches for photos on Unsplash
func (c *UnsplashClient) SearchPhotos(params url.Values) (*models.SearchResult, error) {
	endpoint := fmt.Sprintf("%s/search/photos", c.baseURL)

	// Create request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", c.apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.URL.RawQuery = params.Encode()

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Update rate limit information
	c.updateRateLimits(resp)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse response
	var result models.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}

// GetPhoto gets details of a specific photo
func (c *UnsplashClient) GetPhoto(id string) (*models.Photo, error) {
	endpoint := fmt.Sprintf("%s/photos/%s", c.baseURL, id)

	// Create request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", c.apiKey))
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Update rate limit information
	c.updateRateLimits(resp)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse response
	var photo models.Photo
	if err := json.NewDecoder(resp.Body).Decode(&photo); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &photo, nil
}

// GetRandomPhotos gets random photos
func (c *UnsplashClient) GetRandomPhotos(count int) ([]models.Photo, error) {
	endpoint := fmt.Sprintf("%s/photos/random", c.baseURL)

	// Build query params
	params := url.Values{}
	params.Add("count", fmt.Sprintf("%d", count))

	// Create request
	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", c.apiKey))
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Update rate limit information
	c.updateRateLimits(resp)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse response
	var photos []models.Photo
	if err := json.NewDecoder(resp.Body).Decode(&photos); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return photos, nil
}

// Check if we've hit rate limits
func (c *UnsplashClient) IsRateLimited() bool {
	return c.rateLimit.remaining <= 0 && time.Now().Before(c.rateLimit.resetTime)
}

// Get time until rate limit reset
func (c *UnsplashClient) TimeToRateLimitReset() time.Duration {
	if !c.IsRateLimited() {
		return 0
	}
	return time.Until(c.rateLimit.resetTime)
}

// Update rate limit information from response headers
func (c *UnsplashClient) updateRateLimits(resp *http.Response) {
	if limitStr := resp.Header.Get("X-Ratelimit-Limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &c.rateLimit.limit)
	}

	if remainingStr := resp.Header.Get("X-Ratelimit-Remaining"); remainingStr != "" {
		fmt.Sscanf(remainingStr, "%d", &c.rateLimit.remaining)
	}

	if resetStr := resp.Header.Get("X-Ratelimit-Reset"); resetStr != "" {
		var resetUnix int64
		fmt.Sscanf(resetStr, "%d", &resetUnix)
		c.rateLimit.resetTime = time.Unix(resetUnix, 0)
	}

	c.rateLimit.lastCheckTime = time.Now()
}
