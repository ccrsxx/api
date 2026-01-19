package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/utils"
)

func GetIpAddress(w http.ResponseWriter, r *http.Request) error {
	ipAddress := utils.GetIpAddressFromRequest(r)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(ipAddress)); err != nil {
		return fmt.Errorf("ip response error: %w", err)
	}

	return nil
}

func GetIpInfo(w http.ResponseWriter, r *http.Request) error {
	queryIp := r.URL.Query().Get("ip")

	if queryIp != "" {
		if net.ParseIP(queryIp) == nil {
			return api.NewHttpError(http.StatusBadRequest, "Invalid IP address", nil)
		}
	}

	ipAddress := queryIp

	if ipAddress == "" {
		ipAddress = utils.GetIpAddressFromRequest(r)
	}

	parsedIp := net.ParseIP(ipAddress)

	// Should never happen since request always contains IP address
	// But just in case, we handle it, in case library panic on nil input in future
	if parsedIp == nil {
		return api.NewHttpError(http.StatusBadRequest, "Invalid IP address", nil)
	}

	ipInfo, err := utils.IPInfo().GetIPInfo(parsedIp)

	if err != nil {
		return fmt.Errorf("get ip info error: %w", err)
	}

	return api.NewSuccessResponse(w, http.StatusOK, ipInfo)
}

func GetHttpHeaders(w http.ResponseWriter, r *http.Request) error {
	headers := utils.GetHttpHeadersFromRequest(r)

	return api.NewSuccessResponse(w, http.StatusOK, headers)
}
