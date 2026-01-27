package home

import (
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/src/internal/api"
	"github.com/ccrsxx/api/src/internal/utils"
)

type controller struct{}

var Controller = &controller{}

func (c *controller) ping(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url"`
	}

	err := api.NewSuccessResponse(w, http.StatusOK, response{
		Message:          "Welcome to the API! The server is up and running.",
		DocumentationURL: utils.GetPublicUrlFromRequest(r) + "/docs",
	})

	if err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("home response error: %w", err))
		return
	}
}
