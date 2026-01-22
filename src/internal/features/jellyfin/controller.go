package jellyfin

import (
	"fmt"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
)

type controller struct{}

var Controller = &controller{}

func (c *controller) GetCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	data, err := Service.GetCurrentlyPlaying(r.Context())

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("jellyfin currently playing response error: %w", err))
		return
	}
}
