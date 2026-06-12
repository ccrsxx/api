package statistics

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

func (c *Controller) GetContentsStatistics(w http.ResponseWriter, r *http.Request) {
	contentType := r.URL.Query().Get("type")

	stats, err := c.service.GetContentsStatistics(r.Context(), contentType)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, stats); err != nil {
		slog.Warn("get content statistics response error", "error", err)
	}
}
