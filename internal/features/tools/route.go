package tools

import (
	"net/http"
)

type Config struct {
	Router                    *http.ServeMux
	ToolsController           *Controller
	SharedGetIpInfoController http.Handler
}

func LoadRoutes(config Config) {
	mux := http.NewServeMux()

	mux.Handle("GET /ipinfo", config.SharedGetIpInfoController)

	mux.HandleFunc("GET /ip", config.ToolsController.GetIpAddress)

	mux.HandleFunc("GET /headers", config.ToolsController.GetHttpHeaders)

	config.Router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
