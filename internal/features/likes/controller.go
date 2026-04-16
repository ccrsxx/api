package likes

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

func (c *Controller) GetLikeStatus(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	ipAddress := utils.GetIPAddressFromRequest(r)

	status, err := c.service.GetLikeStatus(r.Context(), slug, ipAddress)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, status); err != nil {
		slog.Warn("get like status response error", "error", err)
	}
}

func (c *Controller) IncrementLike(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	ipAddress := utils.GetIPAddressFromRequest(r)

	status, err := c.service.IncrementLike(r.Context(), slug, ipAddress)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusCreated, status); err != nil {
		slog.Warn("increment like response error", "error", err)
	}
}
