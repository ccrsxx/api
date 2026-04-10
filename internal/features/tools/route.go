package tools

import (
	"net/http"
)

type Config struct {
	ToolsController *Controller
	SharedGetIpInfo http.Handler
}

func LoadRoutes(router *http.ServeMux, config Config) {
	mux := http.NewServeMux()

	mux.Handle("GET /ipinfo", config.SharedGetIpInfo)

	mux.HandleFunc("GET /ip", config.ToolsController.GetIpAddress)

	mux.HandleFunc("GET /headers", config.ToolsController.GetHttpHeaders)

	router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
