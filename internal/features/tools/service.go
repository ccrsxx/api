package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

type ipInfoClient interface {
	GetIPInfo(net.IP) (*ipinfoLib.Core, error)
}

type Service struct {
	ipInfoClient ipInfoClient
}

type ServiceConfig struct {
	IPInfoClient ipInfoClient
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		ipInfoClient: cfg.IPInfoClient,
	}
}

func (s *Service) GetIPInfo(queryIP string, requestIP string) (*ipinfoLib.Core, error) {
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

	info, err := s.ipInfoClient.GetIPInfo(parsedIP)

	if err != nil {
		return nil, fmt.Errorf("get ip info error: %w", err)
	}

	return info, nil
}
