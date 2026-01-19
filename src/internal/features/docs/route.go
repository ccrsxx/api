package docs

import (
	"github.com/ccrsxx/api-go/src/internal/api"
)

func LoadRoutes(router *api.CustomRouter) {
	router.HandleFunc("GET /docs", getDocs)
}
