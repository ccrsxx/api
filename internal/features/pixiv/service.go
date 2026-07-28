package pixiv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/pixiv"
	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/utils"
)

type pixivClient interface {
	GetBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility, page int) ([]pixiv.Artwork, int, error)
	GetImageStream(ctx context.Context, url string) (io.ReadCloser, error)
}

type Service struct {
	client           pixivClient
	backendPublicURL string
}

type ServiceConfig struct {
	Client           pixivClient
	BackendPublicURL string
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		client:           cfg.Client,
		backendPublicURL: cfg.BackendPublicURL,
	}
}

func (s *Service) GetBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility, page int) ([]model.Bookmark, utils.OffsetPaginationMeta, error) {
	artworks, total, err := s.client.GetBookmarks(ctx, visibility, page)

	if err != nil {
		return nil, utils.OffsetPaginationMeta{}, fmt.Errorf("pixiv bookmarks error: %w", err)
	}

	bookmarks := make([]model.Bookmark, 0, len(artworks))

	for _, artwork := range artworks {
		bookmark, err := parseArtworkToBookmark(artwork, s.backendPublicURL)

		if err != nil {
			slog.Warn("pixiv bookmarks skip invalid artwork parse", "error", err)
			continue
		}

		bookmarks = append(bookmarks, bookmark)
	}

	paginationMeta := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
		Page:        page,
		Limit:       pixiv.MaxBookmarksLimit,
		RecordCount: total,
	})

	return bookmarks, paginationMeta.Meta, nil
}

func (s *Service) GetAllBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility) ([]model.Bookmark, error) {
	var allBookmarks []model.Bookmark

	for page := 1; ; page++ {
		bookmarks, meta, err := s.GetBookmarks(ctx, visibility, page)

		if err != nil {
			return nil, fmt.Errorf("pixiv all bookmarks error: %w", err)
		}

		allBookmarks = append(allBookmarks, bookmarks...)

		if meta.Page >= meta.PageCount {
			break
		}
	}

	return allBookmarks, nil
}

func (s *Service) GetImage(ctx context.Context, imageURL string) (io.ReadCloser, error) {
	imageStream, err := s.client.GetImageStream(ctx, imageURL)

	if errors.Is(err, pixiv.ErrPixivInvalidURL) {
		return nil, &api.HTTPError{
			Message:    "Image url is invalid",
			StatusCode: http.StatusBadRequest,
		}
	}

	if err != nil {
		return nil, fmt.Errorf("pixiv image error: %w", err)
	}

	return imageStream, nil
}
