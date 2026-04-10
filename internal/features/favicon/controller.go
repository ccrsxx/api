package favicon

import (
	"log/slog"
	"net/http"
)

type Controller struct {
	icon []byte
}

func NewController(icon []byte) *Controller {
	return &Controller{
		icon: icon,
	}
}

func (c *Controller) getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(c.icon); err != nil {
		slog.Warn("favicon response error", "error", err)
	}
}
