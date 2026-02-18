package jellyfin

import (
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type controller struct{}

var Controller = &controller{}

func (c *controller) getCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	data, err := Service.GetCurrentlyPlaying(r.Context())

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		slog.Warn("jellyfin currently playing response error", "error", err)
	}
}
