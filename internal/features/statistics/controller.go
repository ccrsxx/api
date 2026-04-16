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

func (c *Controller) GetContentStatistics(w http.ResponseWriter, r *http.Request) {
	contentType := r.PathValue("type")

	stats, err := c.service.GetContentStatistics(r.Context(), contentType)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, stats); err != nil {
		slog.Warn("get content statistics response error", "error", err)
	}
}
