package tools

import (
	"net/http"
)

type Config struct {
	Router                    *http.ServeMux
	ToolsController           *Controller
	SharedGetIPInfoController http.Handler
}

func LoadRoutes(cfg Config) {
	mux := http.NewServeMux()

	mux.Handle("GET /ipinfo", cfg.SharedGetIPInfoController)

	mux.HandleFunc("GET /ip", cfg.ToolsController.GetIPAddress)

	mux.HandleFunc("GET /headers", cfg.ToolsController.GetHTTPHeaders)

	cfg.Router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
