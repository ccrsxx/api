package pixiv

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultBaseURL    = "https://pixiv.net"
	MaxBookmarksLimit = 100
)

type Config struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

type Client struct {
	token      string
	userID     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	httpClient := cfg.HTTPClient

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 8 * time.Second}
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	userID := ""

	if tokenParts := strings.Split(cfg.Token, "_"); len(tokenParts) > 0 {
		userID = tokenParts[0]
	}

	return &Client{
		token:      cfg.Token,
		userID:     userID,
		baseURL:    cfg.BaseURL,
		httpClient: httpClient,
	}
}

func (c *Client) GetBookmarks(ctx context.Context, visibility BookmarkVisibility, page int) ([]Artwork, int, error) {
	offset := (page - 1) * MaxBookmarksLimit

	url := fmt.Sprintf("%s/ajax/user/%s/illusts/bookmarks?tag=&offset=%d&limit=%d&rest=%s", c.baseURL, c.userID, offset, MaxBookmarksLimit, visibility)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, 0, fmt.Errorf("pixiv bookmarks create request error: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", c.token))

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, 0, fmt.Errorf("pixiv bookmarks request error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("pixiv bookmarks close body error", "error", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("pixiv bookmarks response status error: %d", res.StatusCode)
	}

	var response Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, 0, fmt.Errorf("pixiv bookmarks decode response error: %w", err)
	}

	if response.Error {
		return nil, 0, fmt.Errorf("pixiv bookmarks api returned error: %s", response.Message)
	}

	artworks := make([]Artwork, 0, len(response.Body.Works))

	for _, rawArtwork := range response.Body.Works {
		var artwork Artwork

		if err := json.Unmarshal(rawArtwork, &artwork); err != nil {
			slog.Warn("pixiv bookmarks skip invalid artwork parse", "error", err)
			continue
		}

		artworks = append(artworks, artwork)
	}

	return artworks, response.Body.Total, nil
}
