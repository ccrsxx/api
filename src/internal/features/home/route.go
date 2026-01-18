package home

import (
	"github.com/ccrsxx/api-go/src/internal/api"
)

func LoadRoutes(router *api.CustomRouter) {
	router.HandleFunc("GET /{$}", ping)
}
