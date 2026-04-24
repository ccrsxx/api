package og

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type Controller struct {
	service      *Service
	isProduction bool
}

type ControllerConfig struct {
	IsProduction bool
}

func NewController(service *Service, cfg ControllerConfig) *Controller {
	return &Controller{
		service:      service,
		isProduction: cfg.IsProduction,
	}
}

func (c *Controller) GetOg(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	imageStream, err := c.service.GetOg(r.Context(), q.Encode())

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	defer func() {
		if err := imageStream.Close(); err != nil {
			slog.Error("failed to close image stream", "error", err)
		}
	}()

	w.Header().Set("Content-Type", "image/png")

	if c.isProduction {
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
