package tools

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

func (c *Controller) GetIPAddress(w http.ResponseWriter, r *http.Request) {
	ipAddress := utils.GetIPAddressFromRequest(r)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(ipAddress)); err != nil {
		slog.Warn("ip response error", "error", err)
	}
}

func (c *Controller) GetIPInfo(w http.ResponseWriter, r *http.Request) {
	queryIP := r.URL.Query().Get("ip")
	requestIP := utils.GetIPAddressFromRequest(r)

	ipInfo, err := c.service.GetIPInfo(queryIP, requestIP)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, ipInfo); err != nil {
		slog.Warn("ip info response error", "error", err)
	}
}

func (c *Controller) GetHTTPHeaders(w http.ResponseWriter, r *http.Request) {
	headers := utils.GetHTTPHeadersFromRequest(r)

	if err := api.NewSuccessRawResponse(w, http.StatusOK, headers); err != nil {
		slog.Warn("headers response error", "error", err)
	}
}
