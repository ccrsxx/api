package pixiv

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ccrsxx/api/internal/clients/pixiv"
	"github.com/ccrsxx/api/internal/utils"
)

type pixivClient interface {
	GetBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility, page int) ([]pixiv.Artwork, int, error)
}

type Service struct {
	client        pixivClient
	pixivImageURL string
}

type ServiceConfig struct {
	Client        pixivClient
	PixivImageURL string
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		client:        cfg.Client,
		pixivImageURL: cfg.PixivImageURL,
	}
}

type Bookmark struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	ArtistID    string   `json:"artistId"`
	ArtistName  string   `json:"artistName"`
	ImageURL    string   `json:"imageUrl"`
	PixivURL    string   `json:"pixivUrl"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Tags        []string `json:"tags"`
	AiGenerated bool     `json:"aiGenerated"`
}

func (s *Service) GetBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility, page int) ([]Bookmark, utils.OffsetPaginationMeta, error) {
	artworks, total, err := s.client.GetBookmarks(ctx, visibility, page)

	if err != nil {
		return nil, utils.OffsetPaginationMeta{}, fmt.Errorf("pixiv bookmarks error: %w", err)
	}

	bookmarks := make([]Bookmark, 0, len(artworks))

	for _, artwork := range artworks {
		bookmark, err := parseArtworkToBookmark(artwork, s.pixivImageURL)

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

func (s *Service) GetAllBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility) ([]Bookmark, error) {
	var allBookmarks []Bookmark

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
