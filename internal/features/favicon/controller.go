package favicon

import (
	_ "embed"
	"log/slog"
	"net/http"
)

var Controller = &controller{}

type controller struct{}

//go:embed favicon.ico
var icon []byte

func (c *controller) getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(icon); err != nil {
		slog.Warn("favicon response error", "error", err)
	}
}
