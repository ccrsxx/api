package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/api"
)

const (
	defaultCloudflareTurnstileURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

type Client struct {
	apiURL     string
	secretKey  string
	httpClient *http.Client
}

type Config struct {
	APIURL     string
	SecretKey  string
	HTTPClient *http.Client
}

func NewClient(cfg Config) *Client {
	if cfg.APIURL == "" {
		cfg.APIURL = defaultCloudflareTurnstileURL
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 8 * time.Second}
	}

	return &Client{
		apiURL:     defaultCloudflareTurnstileURL,
		secretKey:  cfg.SecretKey,
		httpClient: cfg.HTTPClient,
	}
}

func (c *Client) VerifyTurnstile(ctx context.Context, token string, remoteIP string) error {
	payload := struct {
		Secret   string `json:"secret"`
		Response string `json:"response"`
		RemoteIP string `json:"remoteip"`
	}{
		Secret:   c.secretKey,
		Response: token,
		RemoteIP: remoteIP,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("cloudflare verify turnstile marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("cloudflare verify turnstile creation error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("cloudflare verify turnstile call error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("cloudflare verify turnstile close body error:", "error", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cloudflare verify turnstile unexpected status code: %d", res.StatusCode)
	}

	var response Turnstile

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("cloudflare verify turnstile decode response error: %w", err)
	}

	if !response.Success {
		return &api.HTTPError{
			StatusCode: http.StatusForbidden,
			Message:    "Invalid captcha",
		}
	}

	return nil
}
