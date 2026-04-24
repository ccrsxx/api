package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
)

func TestDecodeJSON(t *testing.T) {
	type payload struct {
		Name string `json:"name" validate:"required,min=3"`
	}

	t.Run("Success", func(t *testing.T) {
		body := strings.NewReader(`{"name":"test"}`)

		r := httptest.NewRequest(http.MethodPost, "/", body)

		var p payload

		if err := api.DecodeJSON(r, &p); err != nil {

			t.Fatalf("unwanted error: %v", err)
		}

		if p.Name != "test" {
			t.Errorf("got name %s, want test", p.Name)
		}
	})

	t.Run("Invalid JSON Format", func(t *testing.T) {
		body := strings.NewReader(`invalid-json-string`)

		r := httptest.NewRequest(http.MethodPost, "/", body)

		var p payload

		err := api.DecodeJSON(r, &p)

		if err == nil {
			t.Fatal("got nil, want error")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want *api.HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", httpErr.StatusCode)
		}
	})

	t.Run("Validation Failed", func(t *testing.T) {
		body := strings.NewReader(`{"name":"ab"}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		var p payload

		err := api.DecodeJSON(r, &p)

		if err == nil {
			t.Fatal("got nil, want error")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want *api.HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", httpErr.StatusCode)
		}
	})
}
