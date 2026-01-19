package tools

import (
	"net"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/utils"
)

func GetIpAddress(w http.ResponseWriter, r *http.Request) error {
	ipAddress := utils.GetIpAddressFromRequest(r)

	return api.NewSuccessResponse(w, http.StatusOK, ipAddress)
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

	ipInfo, err := utils.IPInfo().GetIPInfo(net.ParseIP(ipAddress))

	if err != nil {
		return err
	}

	return api.NewSuccessResponse(w, http.StatusOK, ipInfo)
}

func GetHttpHeaders(w http.ResponseWriter, r *http.Request) error {
	headers := utils.GetHttpHeadersFromRequest(r)

	return api.NewSuccessResponse(w, http.StatusOK, headers)
}
