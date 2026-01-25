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

	"github.com/ccrsxx/api-go/src/internal/cache"
	"github.com/ccrsxx/api-go/src/internal/config"
)

type client struct {
	clientID     string
	httpClient   *http.Client
	clientSecret string
	refreshToken string
}

var (
	once     sync.Once
	instance client
)

func Client() *client {
	once.Do(func() {
		instance = client{
			clientID:     config.Env().SpotifyClientID,
			httpClient:   &http.Client{Timeout: 8 * time.Second},
			clientSecret: config.Env().SpotifyClientSecret,
			refreshToken: config.Env().SpotifyRefreshToken,
		}
	})

	return &instance
}

func (c *client) GetCurrentlyPlaying(ctx context.Context) (*SpotifyCurrentlyPlaying, error) {
	token, err := c.getAccessToken(ctx)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)

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
	if res.StatusCode == http.StatusNoContent {
		slog.Info("spotify currently playing no content")

		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify currently playing request status error: %s", res.Status)
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

func (c *client) getAccessToken(ctx context.Context) (string, error) {
	type tokenResponse struct {
		ExpiresIn   int    `json:"expires_in"`
		AccessToken string `json:"access_token"` // expiry time in seconds
	}

	fetcher := func() (tokenResponse, error) {
		requestBody := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {c.refreshToken},
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(requestBody.Encode()))

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

		var tokenResp spotifyTokenResponse

		if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
			return tokenResponse{}, fmt.Errorf("spotify access token decode error: %w", err)
		}

		return tokenResponse{
			ExpiresIn:   tokenResp.ExpiresIn,
			AccessToken: tokenResp.AccessToken,
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
