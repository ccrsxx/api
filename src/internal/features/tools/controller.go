package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/clients/ipinfo"
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

	if queryIp != "" && net.ParseIP(queryIp) == nil {
		api.HandleHttpError(w, r, api.NewHttpError(http.StatusBadRequest, "Invalid IP address", nil))
		return
	}

	ipAddress := queryIp

	if ipAddress == "" {
		ipAddress = utils.GetIpAddressFromRequest(r)
	}

	parsedIp := net.ParseIP(ipAddress)

	// Should never happen since request always contains IP address
	// But just in case, we handle it, in case library panic on nil input in future
	if parsedIp == nil {
		api.HandleHttpError(w, r, api.NewHttpError(http.StatusBadRequest, "Invalid IP address", nil))
		return
	}

	ipInfo, err := ipinfo.Client().GetIPInfo(parsedIp)

	if err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("get ip info error: %w", err))
		return
	}

	if err = api.NewSuccessResponse(w, http.StatusOK, ipInfo); err != nil {
		api.HandleHttpError(w, r, err)
		return
	}
}

func (c *controller) GetHttpHeaders(w http.ResponseWriter, r *http.Request) {
	headers := utils.GetHttpHeadersFromRequest(r)

	if err := api.NewSuccessResponse(w, http.StatusOK, headers); err != nil {
		api.HandleHttpError(w, r, err)
		return
	}
}
