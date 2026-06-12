package views

import (
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		service: svc,
	}
}

func (c *Controller) GetViewCount(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	views, err := c.service.GetViewCount(r.Context(), slug)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, views); err != nil {
		slog.Warn("get view count response error", "error", err)
	}
}

func (c *Controller) IncrementView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	ipAddress := utils.GetIPAddressFromRequest(r)

	views, err := c.service.IncrementView(r.Context(), slug, ipAddress)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, views); err != nil {
		slog.Warn("increment view response error", "error", err)
	}
}
