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

func (s *Service) getIPInfo(queryIP string, requestIP string) (*ipinfoLib.Core, error) {
	if queryIP != "" && net.ParseIP(queryIP) == nil {
		return nil, &api.HTTPError{
			Message:    "Invalid IP address",
			StatusCode: http.StatusBadRequest,
		}
	}

	ipAddress := queryIP

	if ipAddress == "" {
		ipAddress = requestIP
	}

	parsedIP := net.ParseIP(ipAddress)

	if parsedIP == nil {
		return nil, &api.HTTPError{
			Message:    "Invalid IP address",
			StatusCode: http.StatusBadRequest,
		}
	}

	info, err := s.fetcher(parsedIP)

	if err != nil {
		return nil, fmt.Errorf("get ip info error: %w", err)
	}

	return info, nil
}
