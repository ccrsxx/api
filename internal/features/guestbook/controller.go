package guestbook

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

func (c *Controller) GetGuestbook(w http.ResponseWriter, r *http.Request) {
	guestbook, err := c.service.ListGuestbook(r.Context())

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, guestbook); err != nil {
		slog.Warn("get guestbook response error", "error", err)
	}
}

func (c *Controller) CreateGuestbook(w http.ResponseWriter, r *http.Request) {
	var input CreateGuestbookInput

	if err := api.DecodeJSON(r, &input); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	guestbook, err := c.service.CreateGuestbook(r.Context(), input)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusCreated, guestbook); err != nil {
		slog.Warn("create guestbook response error", "error", err)
	}
}

func (c *Controller) DeleteGuestbook(w http.ResponseWriter, r *http.Request) {
	guestbookID := r.PathValue("id")

	if err := c.service.DeleteGuestbook(r.Context(), guestbookID); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
