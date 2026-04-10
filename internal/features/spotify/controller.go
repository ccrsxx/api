package spotify

import (
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type controller struct {
	service *service
}

func NewController(svc *service) *controller {
	return &controller{
		service: svc,
	}
}

func (c *controller) getCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	data, err := c.service.GetCurrentlyPlaying(r.Context())

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		slog.Warn("spotify currently playing response error", "error", err)
	}
}
