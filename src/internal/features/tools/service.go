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

func (s *service) getIpInfo(queryIp string, requestIp string) (*ipinfoLib.Core, error) {
	if queryIp != "" && net.ParseIP(queryIp) == nil {
		return nil, &api.HttpError{
			Message: "Invalid IP address", Details: nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	ipAddress := queryIp

	if ipAddress == "" {
		ipAddress = requestIp
	}

	parsedIp := net.ParseIP(ipAddress)

	if parsedIp == nil {
		return nil, &api.HttpError{
			Message: "Invalid IP address", Details: nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	info, err := ipinfo.Client().GetIPInfo(parsedIp)

	if err != nil {
		return nil, fmt.Errorf("get ip info error: %w", err)
	}

	return info, nil
}
