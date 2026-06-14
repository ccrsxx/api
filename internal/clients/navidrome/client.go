package navidrome

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/api"
)

type Config struct {
	URL        string
	Username   string
	Password   string
	HTTPClient *http.Client
}

type Client struct {
	url        string
	username   string
	password   string
	authParams string
	httpClient *http.Client
}

const (
	defaultSubsonicAPIVersion = "1.16.1"
)

func NewClient(cfg Config) *Client {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 8 * time.Second}
	}

	authParams := createAuthParams(cfg.Username, cfg.Password)

	return &Client{
		url:        cfg.URL,
		username:   cfg.Username,
		password:   cfg.Password,
		authParams: authParams,
		httpClient: cfg.HTTPClient,
	}
}

func (c *Client) GetNowPlaying(ctx context.Context) ([]NowPlayingEntry, error) {
	url := fmt.Sprintf("%s/rest/getNowPlaying?%s", c.url, c.authParams)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return []NowPlayingEntry{}, fmt.Errorf("navidrome now playing request creation error: %w", err)
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return []NowPlayingEntry{}, fmt.Errorf("navidrome now playing request call error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("navidrome now playing close body error:", "error", err)
		}
	}()

	// nolint:nilaway
	if res.StatusCode != http.StatusOK {
		return []NowPlayingEntry{}, fmt.Errorf("navidrome now playing request status error: %d", res.StatusCode)
	}

	var subsonicJSONWrapper JSONWrapper

	if err := json.NewDecoder(res.Body).Decode(&subsonicJSONWrapper); err != nil {
		return []NowPlayingEntry{}, fmt.Errorf("navidrome now playing decode response error: %w", err)
	}

	if subsonicJSONWrapper.Subsonic.NowPlaying == nil {
		return []NowPlayingEntry{}, errors.New("navidrome now playing entry missing")
	}

	return subsonicJSONWrapper.Subsonic.NowPlaying.Entry, nil
}

func (c *Client) GetCoverArtStream(ctx context.Context, covertArtID string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/rest/getCoverArt?%s&id=%s&size=300", c.url, c.authParams, covertArtID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("navidrome cover art request creation error: %w", err)
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("navidrome cover art request call error: %w", err)
	}

	closeBody := func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("navidrome now playing close body error:", "error", err)
		}
	}

	if res.StatusCode != http.StatusOK {
		closeBody()

		return nil, fmt.Errorf("navidrome now playing request status error: %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")

	// If Content/Type is application/json it's guaranteed to be an error, correct return should be image/webp
	if contentType == "application/json" {
		defer closeBody()

		var subsonicJSONWrapper JSONWrapper

		if err := json.NewDecoder(res.Body).Decode(&subsonicJSONWrapper); err != nil {
			return nil, fmt.Errorf("navidrome cover art decode response error: %w", err)
		}

		// Error code 70 means data not found
		// Ref: https://subsonic.org/pages/api.jsp
		isCoverArtNotFoundError := subsonicJSONWrapper.Subsonic.Error != nil && subsonicJSONWrapper.Subsonic.Error.Code == 70

		if isCoverArtNotFoundError {
			return nil, &api.HTTPError{
				StatusCode: http.StatusNotFound,
				Message:    "Cover art not found",
			}
		}

		return nil, fmt.Errorf("navidrome now playing request status error: %d", res.StatusCode)
	}

	return res.Body, nil
}

// createAuthParams builds the authentication query parameters required by the Subsonic API.
//
// It implements the token-based auth scheme introduced in API v1.13.0, where the password
// is never sent in plaintext. Instead, a random salt is generated per request and combined
// with the password to produce a one-time token:
//
//	token = md5(password + salt)
//
// The returned string contains the full set of required query parameters:
//
//	u  - username
//	t  - token (MD5 hash of password + salt)
//	s  - salt (random hex string, unique per request)
//	v  - Subsonic API version
//	c  - client identifier
//	f  - response format
//
// Ref: https://subsonic.org/pages/api.jsp
func createAuthParams(username, password string) string {
	saltByte := make([]byte, 6) // 6 bytes = 12 hex characters

	// rand.Read never returns an error (panics internally instead)
	_, _ = rand.Read(saltByte)

	salt := hex.EncodeToString(saltByte)

	// Subsonic token formula: md5(password + salt)
	hash := md5.Sum([]byte(password + salt))
	token := hex.EncodeToString(hash[:])

	return fmt.Sprintf("u=%s&t=%s&s=%s&v=%s&c=api&f=json", username, token, salt, defaultSubsonicAPIVersion)
}
