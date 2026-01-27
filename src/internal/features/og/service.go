package og

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/config"
)

type service struct {
	httpClient *http.Client
}

var Service = &service{
	httpClient: &http.Client{Timeout: 8 * time.Second},
}

type ogParams struct {
	Title       string `url:"title"`
	Description string
	IconUrl     string
	ImageUrl    string
	ThemeColor  string
}

func (s *service) getOg(ctx context.Context, query string) (io.ReadCloser, error) {
	ogUrl := "http://og:4444/og"

	if config.Config().IsDevelopment {
		ogUrl = "http://localhost:4444/og"
	}

	targetUrl := fmt.Sprintf("%s?%s", ogUrl, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetUrl, nil)

	if err != nil {
		return nil, fmt.Errorf("og request creation error: %w", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("og request call error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("og response close body error", "error", err)
		}

		return nil, fmt.Errorf("og request status error: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
