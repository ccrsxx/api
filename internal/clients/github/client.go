package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
}

type Config struct {
	APIURL string
}

const (
	defaultGithubUserURL    = "https://api.github.com/user"
	defaultGithubAPIVersion = "2026-03-10"
)

func NewClient(cfg Config) *Client {
	if cfg.APIURL == "" {
		cfg.APIURL = defaultGithubUserURL
	}

	return &Client{
		apiURL:     cfg.APIURL,
		httpClient: &http.Client{Timeout: 8 * time.Second},
	}
}

func (c *Client) GetCurrentUser(ctx context.Context, accessToken string) (User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL, nil)

	if err != nil {
		return User{}, fmt.Errorf("github get user creation error: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	req.Header.Set("Authorization", "Bearer "+accessToken)

	req.Header.Set("X-GitHub-Api-Version", defaultGithubAPIVersion)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return User{}, fmt.Errorf("github get user call error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("github get user close body error:", "error", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("github get user status error: %s", res.Status)
	}

	var githubUser User

	if err := json.NewDecoder(res.Body).Decode(&githubUser); err != nil {
		return User{}, fmt.Errorf("github profile decode response error: %w", err)
	}

	return githubUser, nil
}
