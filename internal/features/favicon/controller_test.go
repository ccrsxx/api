package favicon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestController_getFavicon(t *testing.T) {
	originalIcon := icon

	defer func() {
		icon = originalIcon
	}()

	icon = []byte("fake-icon-data")

	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
		w := httptest.NewRecorder()

		Controller.getFavicon(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "image/x-icon" {
			t.Errorf("got %s, want Content-Type image/x-icon", contentType)
		}

		if w.Body.String() != "fake-icon-data" {
			t.Error("want body to contain icon data")
		}
	})

	t.Run("Response Write Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		Controller.getFavicon(w, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
