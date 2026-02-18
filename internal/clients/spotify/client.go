package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ccrsxx/api/internal/cache"
	"github.com/ccrsxx/api/internal/config"
)

const (
	defaultAuthURL = "https://accounts.spotify.com/api/token"
	defaultApiURL  = "https://api.spotify.com/v1/me/player/currently-playing"
)

type Client struct {
	apiURL       string
	authURL      string
	clientID     string
	clientSecret string
	refreshToken string
	httpClient   *http.Client
}

var (
	once     sync.Once
	instance *Client
)

func New(clientID, clientSecret, refreshToken, authURL, apiURL string) *Client {
	return &Client{
		apiURL:       apiURL,
		authURL:      authURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		httpClient:   &http.Client{Timeout: 8 * time.Second},
	}
}

func DefaultClient() *Client {
	once.Do(func() {
		instance = New(
			config.Env().SpotifyClientID,
			config.Env().SpotifyClientSecret,
			config.Env().SpotifyRefreshToken,
			defaultAuthURL,
			defaultApiURL,
		)
	})

	return instance
}

func (c *Client) GetCurrentlyPlaying(ctx context.Context) (*SpotifyCurrentlyPlaying, error) {
	token, err := c.getAccessToken(ctx)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.apiURL, nil)

	if err != nil {
		return nil, fmt.Errorf("spotify currently playing request creation error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("spotify currently playing request call error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("spotify currently playing close body error", "error", err)
		}
	}()

	// 204 No Content means nothing is currently playing
	// nolint:nilaway
	if res.StatusCode == http.StatusNoContent {
		slog.Debug("spotify currently playing no content")

		return nil, nil
	}

	// nolint:nilaway
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify currently playing request status error: %d", res.StatusCode)
	}

	var data SpotifyCurrentlyPlaying

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("spotify currently playing decode error: %w", err)
	}

	if data.Item == nil || data.Item.Type != "track" {
		return nil, fmt.Errorf("spotify currently playing invalid item type: %v", data.Item)
	}

	return &data, nil
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	type tokenResponse struct {
		ExpiresIn   int    `json:"expires_in"`
		AccessToken string `json:"access_token"` // expiry time in seconds
	}

	fetcher := func() (tokenResponse, error) {
		requestBody := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {c.refreshToken},
		}

		req, err := http.NewRequestWithContext(ctx, "POST", c.authURL, strings.NewReader(requestBody.Encode()))

		if err != nil {
			return tokenResponse{}, fmt.Errorf("spotify access token request creation error: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		authString := c.clientID + ":" + c.clientSecret

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

		req.Header.Set("Authorization", "Basic "+encodedAuth)

		res, err := c.httpClient.Do(req)

		if err != nil {
			return tokenResponse{}, fmt.Errorf("spotify access token request call error: %w", err)
		}

		defer func() {
			if err := res.Body.Close(); err != nil {
				slog.Warn("spotify access token close body error", "error", err)
			}
		}()

		if res.StatusCode != http.StatusOK {
			return tokenResponse{}, fmt.Errorf("spotify access token request status error: %s", res.Status)
		}

		type spotifyTokenResponse struct {
			Scope       string `json:"scope"`
			ExpiresIn   int    `json:"expires_in"`
			TokenType   string `json:"token_type"`
			AccessToken string `json:"access_token"`
		}

		var tokenRes spotifyTokenResponse

		if err := json.NewDecoder(res.Body).Decode(&tokenRes); err != nil {
			return tokenResponse{}, fmt.Errorf("spotify access token decode error: %w", err)
		}

		return tokenResponse{
			ExpiresIn:   tokenRes.ExpiresIn,
			AccessToken: tokenRes.AccessToken,
		}, nil
	}

	ttlFunc := func(data tokenResponse) time.Duration {
		// Add 60 seconds buffer to avoid using expired tokens
		bufferExpiryOffset := 60 * time.Second

		expiresIn := time.Duration(data.ExpiresIn) * time.Second

		return expiresIn - bufferExpiryOffset
	}

	data, err := cache.GetCachedData(
		ctx,
		"api:spotify:access_token",
		"memory",
		fetcher,
		ttlFunc,
	)

	if err != nil {
		return "", fmt.Errorf("spotify cache access token error: %w", err)
	}

	return data.AccessToken, nil
}
