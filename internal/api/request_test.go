package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	t.Run("Success", func(t *testing.T) {
		body := strings.NewReader(`{"name":"test"}`)

		r := httptest.NewRequest(http.MethodPost, "/", body)

		var p payload

		if err := DecodeJSON(r, &p); err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if p.Name != "test" {
			t.Errorf("got name %s, want test", p.Name)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := strings.NewReader(`invalid`)

		r := httptest.NewRequest(http.MethodPost, "/", body)

		var p payload

		err := DecodeJSON(r, &p)

		if err == nil {
			t.Fatal("got nil, want error")
		}

		httpErr, ok := errors.AsType[*HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", httpErr.StatusCode)
		}
	})
}
