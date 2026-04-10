package jellyfin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	URL      string
	ApiKey   string
	ImageURL string
	Username string
}

type Client struct {
	url        string
	apiKey     string
	imageURL   string
	username   string
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		url:        cfg.URL,
		apiKey:     cfg.ApiKey,
		imageURL:   cfg.ImageURL,
		username:   cfg.Username,
		httpClient: &http.Client{Timeout: 8 * time.Second},
	}
}

func (c *Client) GetSessions(ctx context.Context) ([]SessionInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+"/Sessions", nil)

	if err != nil {
		return nil, fmt.Errorf("jellyfin currently playing request creation error: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-MediaBrowser-Token", c.apiKey)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("jellyfin currently playing request call error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println("jellyfin currently playing close body error:", err)
		}
	}()

	// nolint:nilaway
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jellyfin currently playing request status error: %d", res.StatusCode)
	}

	var sessions []SessionInfo

	if err := json.NewDecoder(res.Body).Decode(&sessions); err != nil {
		return nil, fmt.Errorf("jellyfin currently playing decode response error: %w", err)
	}

	return sessions, nil
}
