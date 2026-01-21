package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/clients/ipinfo"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

type service struct{}

var Service = &service{}

func (s *service) GetIpInfo(queryIp string, requestIp string) (*ipinfoLib.Core, error) {
	if queryIp != "" && net.ParseIP(queryIp) == nil {
		return nil, api.NewHttpError(http.StatusBadRequest, "Invalid IP address", nil)
	}

	ipAddress := queryIp

	if ipAddress == "" {
		ipAddress = requestIp
	}

	parsedIp := net.ParseIP(ipAddress)

	if parsedIp == nil {
		return nil, api.NewHttpError(http.StatusBadRequest, "Invalid IP address", nil)
	}

	info, err := ipinfo.Client().GetIPInfo(parsedIp)

	if err != nil {
		return nil, fmt.Errorf("get ip info error: %w", err)
	}

	return info, nil
}
