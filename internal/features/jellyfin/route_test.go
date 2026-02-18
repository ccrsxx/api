package jellyfin

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	LoadRoutes(mux)

	tests := []test.RouteTestCase{

		{
			Path:   "/jellyfin/currently-playing",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
