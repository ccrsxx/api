package contents

import (
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type Controller struct {
	service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		service: svc,
	}
}

func (c *Controller) GetContentData(w http.ResponseWriter, r *http.Request) {
	contentType := r.URL.Query().Get("type")

	data, err := c.service.GetContentData(r.Context(), contentType)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		slog.Warn("get content data response error", "error", err)
	}
}

func (c *Controller) UpsertContent(w http.ResponseWriter, r *http.Request) {
	var input UpsertContentInput

	if err := api.DecodeJSON(r, &input); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	content, err := c.service.UpsertContent(r.Context(), input)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusCreated, content); err != nil {
		slog.Warn("upsert content response error", "error", err)
	}
}
