package home

import (
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/utils"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	err := api.NewSuccessResponse(w, http.StatusOK, model.Ping{
		Message:          "Welcome to the API! The server is up and running.",
		DocumentationURL: utils.GetPublicURLFromRequest(r) + "/docs",
	})

	if err != nil {
		slog.Warn("home response error", "error", err)
	}
}
