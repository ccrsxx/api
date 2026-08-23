package contacts_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/contacts"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	svc := contacts.NewService(contacts.ServiceConfig{})

	ctx := t.Context()

	authMw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

	contacts.LoadRoutes(contacts.Config{
		Router:         mux,
		Service:        svc,
		AppContext:     ctx,
		AuthMiddleware: authMw,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/contacts",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
