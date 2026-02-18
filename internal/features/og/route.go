package og

import (
	"net/http"
)

func LoadRoutes(router *http.ServeMux) {
	router.HandleFunc("/og", Controller.getOg)
}
