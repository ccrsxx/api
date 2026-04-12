package tools

import (
	"net"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/test"
	"github.com/ipinfo/go/v2/ipinfo"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	svc := NewService(ServiceConfig{
		Fetcher: func(ip net.IP) (*ipinfo.Core, error) {
			return &ipinfo.Core{IP: ip}, nil
		},
	})

	ctrl := NewController(svc)

	config := Config{
		ToolsController:           ctrl,
		SharedGetIPInfoController: http.HandlerFunc(ctrl.GetIPInfo),
	}

	LoadRoutes(Config{
		Router:                    mux,
		ToolsController:           config.ToolsController,
		SharedGetIPInfoController: config.SharedGetIPInfoController,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/tools/ip",
			Method: http.MethodGet,
		},
		{
			Path:   "/tools/headers",
			Method: http.MethodGet,
		},
		{
			Path:   "/tools/ipinfo",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
