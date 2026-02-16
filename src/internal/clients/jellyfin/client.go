package jellyfin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ccrsxx/api/src/internal/config"
)

type Client struct {
	url        string
	apiKey     string
	imageUrl   string
	username   string
	httpClient *http.Client
}

var (
	once   sync.Once
	client *Client
)

func New(url, apiKey, imageUrl, username string) *Client {
	return &Client{
		url:        url,
		apiKey:     apiKey,
		imageUrl:   imageUrl,
		username:   username,
		httpClient: &http.Client{Timeout: 8 * time.Second},
	}
}

func DefaultClient() *Client {
	once.Do(func() {
		client = New(
			config.Env().JellyfinUrl,
			config.Env().JellyfinApiKey,
			config.Env().JellyfinImageUrl,
			config.Env().JellyfinUsername,
		)
	})

	return client
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jellyfin currently playing request status error: %d", res.StatusCode)
	}

	var sessions []SessionInfo

	if err := json.NewDecoder(res.Body).Decode(&sessions); err != nil {
		return nil, fmt.Errorf("jellyfin currently playing decode response error: %w", err)
	}

	return sessions, nil
}
