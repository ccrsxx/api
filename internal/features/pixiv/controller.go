package pixiv

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/pixiv"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	page := 1

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	visibility := pixiv.BookmarkVisibilityPublic

	if r.URL.Query().Get("type") == "private" {
		visibility = pixiv.BookmarkVisibilityPrivate
	}

	bookmarks, meta, err := c.service.GetBookmarks(r.Context(), visibility, page)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessPaginatedResponse(w, http.StatusOK, meta, bookmarks); err != nil {
		slog.Warn("pixiv bookmarks response error", "error", err)
	}
}

func (c *Controller) GetAllBookmarks(w http.ResponseWriter, r *http.Request) {
	visibility := pixiv.BookmarkVisibilityPublic

	if r.URL.Query().Get("type") == "private" {
		visibility = pixiv.BookmarkVisibilityPrivate
	}

	bookmarks, err := c.service.GetAllBookmarks(r.Context(), visibility)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := api.NewSuccessResponse(w, http.StatusOK, bookmarks); err != nil {
		slog.Warn("pixiv all bookmarks response error", "error", err)
	}
}
