package docs

import (
	"net/http"
)

func LoadRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /docs", Controller.getDocs)
}
