package jellyfin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ccrsxx/api-go/src/internal/config"
)

type client struct {
	url        string
	apiKey     string
	imageUrl   string
	username   string
	httpClient *http.Client
}

var (
	once     sync.Once
	instance client
)

func Client() *client {
	once.Do(func() {
		instance = client{
			url:        config.Env().JellyfinUrl,
			apiKey:     config.Env().JellyfinApiKey,
			imageUrl:   config.Env().JellyfinImageUrl,
			username:   config.Env().JellyfinUsername,
			httpClient: &http.Client{Timeout: 1 * time.Second},
		}
	})

	return &instance
}

func (c *client) GetSessions(ctx context.Context) (*[]SessionInfo, error) {
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

	return &sessions, nil
}
