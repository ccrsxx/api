package server

import (
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/features/metrics"
)

func NewMonitoringServer(cfg config.AppConfig) *http.Server {
	addr := ":" + strconv.Itoa(cfg.MonitoringPort)

	router := http.NewServeMux()

	metrics.LoadRoutes(
		metrics.Config{
			Router: router,
		},
	)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
