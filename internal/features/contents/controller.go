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
	contentType := r.PathValue("type")

	data, err := c.service.GetContentData(r.Context(), contentType)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, data); err != nil {
		slog.Warn("get content data response error", "error", err)
	}
}

type upsertContentRequest struct {
	Slug string `json:"slug"`
}

func (c *Controller) UpsertContent(w http.ResponseWriter, r *http.Request) {
	contentType := r.PathValue("type")

	var body upsertContentRequest

	if err := api.DecodeJSON(r, &body); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	content, err := c.service.UpsertContent(r.Context(), body.Slug, contentType)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusCreated, content); err != nil {
		slog.Warn("upsert content response error", "error", err)
	}
}
