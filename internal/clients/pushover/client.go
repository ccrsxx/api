package pushover

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	apiURL     string
	token      string
	userKey    string
	httpClient *http.Client
}

type Config struct {
	APIURL     string
	Token      string
	UserKey    string
	HTTPClient *http.Client
}

const (
	defaultPushoverAPIURL = "https://api.pushover.net/1/messages.json"
)

func NewClient(cfg Config) *Client {
	if cfg.APIURL == "" {
		cfg.APIURL = defaultPushoverAPIURL
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 8 * time.Second}
	}

	return &Client{
		apiURL:     cfg.APIURL,
		token:      cfg.Token,
		userKey:    cfg.UserKey,
		httpClient: cfg.HTTPClient,
	}
}

func (c *Client) SendMessage(ctx context.Context, messageRequest MessageRequest) error {
	messageRequest.Token = c.token
	messageRequest.User = c.userKey

	body, err := json.Marshal(messageRequest)

	if err != nil {
		return fmt.Errorf("pushover send message marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("pushover send message request creation error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("pushover send message request call error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("pushover send message close body error", "error", err)
		}
	}()

	// nolint:nilaway
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("pushover send message request status error: %s", res.Status)
	}

	var data MessageResponse

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return fmt.Errorf("pushover send message decode response error: %w", err)
	}

	return nil
}
