package tools

import (
	"net/http"
)

type Config struct {
	Router                    *http.ServeMux
	ToolsController           *Controller
	SharedGetIpInfoController http.Handler
}

func LoadRoutes(cfg Config) {
	mux := http.NewServeMux()

	mux.Handle("GET /ipinfo", cfg.SharedGetIpInfoController)

	mux.HandleFunc("GET /ip", cfg.ToolsController.GetIpAddress)

	mux.HandleFunc("GET /headers", cfg.ToolsController.GetHttpHeaders)

	cfg.Router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
