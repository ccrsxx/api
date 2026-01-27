package og

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/config"
)

type controller struct {
}

var Controller = &controller{}

func (c *controller) getOg(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	imageStream, err := Service.getOg(r.Context(), q.Encode())

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	defer func() {
		if err := imageStream.Close(); err != nil {
			slog.Error("failed to close image stream", "error", err)
		}
	}()

	w.Header().Set("Content-Type", "image/png")

	if config.Config().IsProduction {
		// Cache Policy: Aggressive (1 Year)
		// - public:       Allows CDNs and shared proxies to cache this.
		// - immutable:    Prevents browsers from sending "Is this modified?" (304) checks on refresh.
		// - no-transform: Prevents mobile carriers from compressing/blurring the image.
		// - max-age:      31536000 seconds = 1 Year.
		w.Header().Set("Cache-Control", "public, immutable, no-transform, max-age=31536000")
	}

	if _, err := io.Copy(w, imageStream); err != nil {
		slog.Warn("og response write error", "error", err)
	}
}
