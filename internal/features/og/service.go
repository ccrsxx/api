package og

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/config"
)

type service struct {
	ogUrl      string
	httpClient *http.Client
}

var Service = &service{
	ogUrl:      "http://10.0.0.60:4444/og",
	httpClient: &http.Client{Timeout: 8 * time.Second},
}

func (s *service) getOg(ctx context.Context, query string) (io.ReadCloser, error) {
	ogUrl := s.ogUrl

	if config.Config().IsDevelopment {
		ogUrl = "http://localhost:4444/og"
	}

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
