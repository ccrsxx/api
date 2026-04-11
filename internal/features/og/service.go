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
	ogUrl      string
	httpClient *http.Client
}

type ServiceConfig struct {
	OgUrl      string
	HttpClient *http.Client
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.HttpClient == nil {
		cfg.HttpClient = &http.Client{
			Timeout: 8 * time.Second,
		}
	}

	return &Service{
		ogUrl:      cfg.OgUrl,
		httpClient: cfg.HttpClient,
	}
}

func (s *Service) getOg(ctx context.Context, query string) (io.ReadCloser, error) {
	ogUrl := s.ogUrl

	targetUrl := fmt.Sprintf("%s?%s", ogUrl, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetUrl, nil)

	if err != nil {
		return nil, fmt.Errorf("og request creation error: %w", err)
	}

	res, err := s.httpClient.Do(req)

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
