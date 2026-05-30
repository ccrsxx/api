package pixiv

import (
	"context"
	"errors"
	"testing"

	"github.com/ccrsxx/api/internal/clients/pixiv"
)

type mockPixivClient struct {
	artworks []pixiv.Artwork
	total    int
	err      error
}

func (m *mockPixivClient) GetBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility, page int) ([]pixiv.Artwork, int, error) {
	return m.artworks, m.total, m.err
}

func TestService_GetBookmarks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			artworks: []pixiv.Artwork{
				{
					ID:             "123",
					URL:            "https://i.pximg.net/img-master/test.jpg",
					Title:          "Test",
					UserID:         "456",
					UserName:       "Artist",
					IsBookmarkable: true,
					Width:          800,
					Height:         600,
				},
			},
			total: 1,
		}

		svc := NewService(ServiceConfig{
			Client:        mock,
			PixivImageURL: "https://proxy.example.com",
		})

		bookmarks, meta, err := svc.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if len(bookmarks) != 1 {
			t.Errorf("got %d bookmarks, want 1", len(bookmarks))
		}

		if meta.Page != 1 {
			t.Errorf("got page %d, want 1", meta.Page)
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mock := &mockPixivClient{
			err: errors.New("network fail"),
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		_, _, err := svc.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want error")
		}
	})

	t.Run("Invalid Artwork Skipped", func(t *testing.T) {
		mock := &mockPixivClient{
			artworks: []pixiv.Artwork{
				{
					ID:             "1",
					URL:            "https://i.pximg.net/img-master/test.jpg",
					IsBookmarkable: true,
					UserID:         "2",
					Width:          800,
					Height:         600,
				},
				{
					ID:             "2",
					IsBookmarkable: false, // Will be skipped by parseArtworkToBookmark
				},
			},
			total: 2,
		}

		svc := NewService(ServiceConfig{
			Client:        mock,
			PixivImageURL: "https://proxy.example.com",
		})

		bookmarks, _, err := svc.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if len(bookmarks) != 1 {
			t.Errorf("got %d bookmarks, want 1 (non-bookmarkable should be skipped)", len(bookmarks))
		}
	})
}

func TestService_GetAllBookmarks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			artworks: []pixiv.Artwork{
				{
					ID:             "1",
					URL:            "https://i.pximg.net/img-master/test.jpg",
					UserID:         "2",
					IsBookmarkable: true,
					Width:          800,
					Height:         600,
				},
			},
			total: 1,
		}

		svc := NewService(ServiceConfig{
			Client:        mock,
			PixivImageURL: "https://proxy.example.com",
		})

		bookmarks, err := svc.GetAllBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if len(bookmarks) != 1 {
			t.Errorf("got %d bookmarks, want 1", len(bookmarks))
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mock := &mockPixivClient{
			err: errors.New("network fail"),
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		_, err := svc.GetAllBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic)

		if err == nil {
			t.Error("want error")
		}
	})
}
