package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ccrsxx/api/src/internal/api"
	"github.com/ccrsxx/api/src/internal/clients/ipinfo"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

type ipInfoFetcher func(net.IP) (*ipinfoLib.Core, error)

type service struct {
	fetcher ipInfoFetcher
}

var Service = &service{
	fetcher: func(ip net.IP) (*ipinfoLib.Core, error) {
		return ipinfo.DefaultClient().GetIPInfo(ip)
	},
}

func (s *service) getIpInfo(queryIp string, requestIp string) (*ipinfoLib.Core, error) {
	if queryIp != "" && net.ParseIP(queryIp) == nil {
		return nil, &api.HttpError{
			Message:    "Invalid IP address",
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
			Message:    "Invalid IP address",
			StatusCode: http.StatusBadRequest,
		}
	}

	info, err := s.fetcher(parsedIp)

	if err != nil {
		return nil, fmt.Errorf("get ip info error: %w", err)
	}

	return info, nil
}
