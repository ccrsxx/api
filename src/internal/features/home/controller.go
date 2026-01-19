package home

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/utils"
)

func ping(w http.ResponseWriter, r *http.Request) error {
	type response struct {
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url"`
	}

	return api.NewSuccessResponse(w, http.StatusOK, response{
		Message:          "Welcome to the API! The server is up and running.",
		DocumentationURL: utils.GetPublicUrlFromRequest(r) + "/docs",
	})
}
