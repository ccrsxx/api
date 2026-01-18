package home

import (
	"github.com/ccrsxx/api-go/src/internal/api"
)

func LoadRoutes(r *api.CustomRouter) {
	r.HandleFunc("GET /{$}", ping)
}
