package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/features/metrics"
	"github.com/ccrsxx/api/internal/test"
)

func TestController_GetMetrics(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		w := httptest.NewRecorder()

		ctrl := metrics.NewController()

		ctrl.GetMetrics(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}

		if !strings.Contains(w.Header().Get("Content-Type"), "text/plain") {
			t.Errorf("got %s, want Content-Type containing text/plain", w.Header().Get("Content-Type"))
		}

		body := w.Body.String()

		if !strings.Contains(body, "go_goroutines") {
			t.Error("want body to contain 'go_goroutines'")
		}

		if !strings.Contains(body, "process_cpu_seconds_total") {
			t.Error("want body to contain 'process_cpu_seconds_total'")
		}
	})

	t.Run("Response Write Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		ctrl := metrics.NewController()

		ctrl.GetMetrics(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
