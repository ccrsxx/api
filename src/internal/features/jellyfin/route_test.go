package jellyfin

import (
	"context"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/src/internal/clients/jellyfin"
	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	originalKey := config.Env().SecretKey

	defer func() {
		config.Env().SecretKey = originalKey
	}()

	config.Env().SecretKey = "test-secret"

	originalFetcher := Service.fetcher

	defer func() {
		Service.fetcher = originalFetcher
	}()

	Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
		return nil, nil
	}

	mux := http.NewServeMux()
	LoadRoutes(mux)

	tests := []test.RouteTestCase{
		{
			Path:       "/jellyfin/currently-playing",
			Method:     http.MethodGet,
			Headers:    http.Header{"Authorization": []string{"Bearer test-secret"}},
			StatusCode: http.StatusOK,
		},
		{
			Path:       "/jellyfin/currently-playing",
			Method:     http.MethodGet,
			StatusCode: http.StatusUnauthorized, // No auth header should return 401
		},
	}

	test.AssertRoutes(t, mux, tests)
}
