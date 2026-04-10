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

func (c *Controller) GetIpAddress(w http.ResponseWriter, r *http.Request) {
	ipAddress := utils.GetIpAddressFromRequest(r)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(ipAddress)); err != nil {
		slog.Warn("ip response error", "error", err)
	}
}

func (c *Controller) GetIpInfo(w http.ResponseWriter, r *http.Request) {
	queryIp := r.URL.Query().Get("ip")
	requestIp := utils.GetIpAddressFromRequest(r)

	ipInfo, err := c.service.getIpInfo(queryIp, requestIp)

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, ipInfo); err != nil {
		slog.Warn("ip info response error", "error", err)
	}
}

func (c *Controller) GetHttpHeaders(w http.ResponseWriter, r *http.Request) {
	headers := utils.GetHttpHeadersFromRequest(r)

	if err := api.NewSuccessResponse(w, http.StatusOK, headers); err != nil {
		slog.Warn("headers response error", "error", err)
	}
}
