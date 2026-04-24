package og

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	ogURL      string
	HTTPClient *http.Client
}

type ServiceConfig struct {
	OgURL      string
	HTTPClient *http.Client
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 8 * time.Second,
		}
	}

	return &Service{
		ogURL:      cfg.OgURL,
		HTTPClient: cfg.HTTPClient,
	}
}

func (s *Service) GetOg(ctx context.Context, query string) (io.ReadCloser, error) {
	ogURL := s.ogURL

	targetURL := fmt.Sprintf("%s?%s", ogURL, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)

	if err != nil {
		return nil, fmt.Errorf("og request creation error: %w", err)
	}

	res, err := s.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("og request call error: %w", err)
	}

	// nolint:nilaway
	if res.StatusCode != http.StatusOK {
		if err := res.Body.Close(); err != nil {
			slog.Warn("og response close body error", "error", err)
		}

		return nil, fmt.Errorf("og request status error: %d", res.StatusCode)
	}

	return res.Body, nil
}
