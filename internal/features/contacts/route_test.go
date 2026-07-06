package contacts_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/contacts"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	svc := contacts.NewService(contacts.ServiceConfig{})

	contacts.LoadRoutes(contacts.Config{
		Router:  mux,
		Service: svc,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/contacts",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
