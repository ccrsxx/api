package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

type ipInfoFetcher func(net.IP) (*ipinfoLib.Core, error)

type Service struct {
	fetcher ipInfoFetcher
}

type ServiceConfig struct {
	Fetcher ipInfoFetcher
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		fetcher: cfg.Fetcher,
	}
}

func (s *Service) getIpInfo(queryIp string, requestIp string) (*ipinfoLib.Core, error) {
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
