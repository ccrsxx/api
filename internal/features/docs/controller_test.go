package docs_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/features/docs"
	"github.com/ccrsxx/api/internal/test"
)

func TestController_GetDocs(t *testing.T) {
	validJSON := []byte(`{"openapi":"3.0.0","info":{"title":"Test","version":"1.0"}}`)

	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/docs", nil)
		w := httptest.NewRecorder()

		ctrl := docs.NewController(validJSON)

		ctrl.GetDocs(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "text/html" {
			t.Errorf("got %s, want Content-Type text/html", contentType)
		}

		if w.Body.Len() == 0 {
			t.Error("got empty, want body to contain HTML")
		}
	})

	t.Run("Render Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/docs", nil)
		w := httptest.NewRecorder()

		ctrl := docs.NewController(nil)

		ctrl.GetDocs(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want status 500", w.Code)
		}
	})

	t.Run("Response Write Error", func(t *testing.T) {
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
		r := httptest.NewRequest(http.MethodGet, "/docs", nil)

		ctrl := docs.NewController(validJSON)

		ctrl.GetDocs(w, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
