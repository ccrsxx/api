package navidrome

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type Controller struct {
	service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		service: svc,
	}
}

func (c *Controller) GetCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	data, err := c.service.GetCurrentlyPlaying(r.Context())

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		slog.Warn("spotify currently playing response error", "error", err)
	}
}

func (c *Controller) GetCoverArt(w http.ResponseWriter, r *http.Request) {
	coverArtID := r.PathValue("coverArtID")

	coverArtStream, err := c.service.client.GetCoverArtStream(r.Context(), coverArtID)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	defer func() {
		if err := coverArtStream.Close(); err != nil {
			slog.Warn("navidrome cover art close body error", "error", err)
		}
	}()

	w.Header().Set("Content-Type", "image/webp")

	// Cache Policy: Aggressive (1 Year)
	// - public:       Allows CDNs and shared proxies to cache this.
	// - immutable:    Prevents browsers from sending "Is this modified?" (304) checks on refresh.
	// - no-transform: Prevents mobile carriers from compressing/blurring the image.
	// - max-age:      31536000 seconds = 1 Year.
	w.Header().Set("Cache-Control", "public, immutable, no-transform, max-age=31536000")

	if _, err := io.Copy(w, coverArtStream); err != nil {
		slog.Warn("navidrome cover art response write error", "error", err)
	}
}
