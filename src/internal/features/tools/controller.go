package tools

import (
	"fmt"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/utils"
)

type controller struct{}

var Controller = &controller{}

func (c *controller) GetIpAddress(w http.ResponseWriter, r *http.Request) {
	ipAddress := utils.GetIpAddressFromRequest(r)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(ipAddress)); err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("ip response error: %w", err))
		return
	}
}

func (c *controller) GetIpInfo(w http.ResponseWriter, r *http.Request) {
	queryIp := r.URL.Query().Get("ip")
	requestIp := utils.GetIpAddressFromRequest(r)

	ipInfo, err := Service.GetIpInfo(queryIp, requestIp)

	if err != nil {
		api.HandleHttpError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, ipInfo); err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("ip info response error: %w", err))
		return
	}
}

func (c *controller) GetHttpHeaders(w http.ResponseWriter, r *http.Request) {
	headers := utils.GetHttpHeadersFromRequest(r)

	if err := api.NewSuccessResponse(w, http.StatusOK, headers); err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("headers response error: %w", err))
		return
	}
}
