package spotify

import (
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
