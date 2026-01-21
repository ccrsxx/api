package spotify

import (
	"fmt"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/clients/spotify"
)

type controller struct{}

var Controller = &controller{}

func (c *controller) getCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	data, err := spotify.Client().GetNowCurrentlyPlaying(r.Context())

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("spotify currently playing response error: %w", err))
		return
	}
}
