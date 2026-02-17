package favicon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/test"
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
			t.Errorf("want status 200, got %d", w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "image/x-icon" {
			t.Errorf("want Content-Type image/x-icon, got %s", contentType)
		}

		if w.Body.String() != "fake-icon-data" {
			t.Error("want body to contain icon data")
		}
	})

	t.Run("Response Write Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		Controller.getFavicon(w, r)
	})
}
