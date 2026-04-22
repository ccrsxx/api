package home

import (
	"net"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/tools"
	"github.com/ccrsxx/api/internal/test"
	"github.com/ipinfo/go/v2/ipinfo"
)

type mockIPInfoClient struct{}

func (m *mockIPInfoClient) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	return &ipinfo.Core{IP: ip}, nil
}

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	svc := tools.NewService(tools.ServiceConfig{
		IPInfoClient: &mockIPInfoClient{},
	})

	ctrl := tools.NewController(svc)

	LoadRoutes(Config{
		Router:                    mux,
		ToolsController:           ctrl,
		SharedGetIPInfoController: http.HandlerFunc(ctrl.GetIPInfo),
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/",
			Method: http.MethodGet,
		},
		{
			Path:   "/",
			Host:   "ip.example.com",
			Method: http.MethodGet,
		},
		{
			Path:   "/",
			Host:   "ipinfo.example.com",
			Method: http.MethodGet,
		},
		{
			Path:   "/",
			Host:   "headers.example.com",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}

