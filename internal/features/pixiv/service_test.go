package pixiv

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/pixiv"
)

type mockPixivClient struct {
	bookmarkArtworks []pixiv.Artwork
	bookmarksTotal   int
	bookmarksErr     error

	imageStream io.ReadCloser
	imageErr    error
}

func (m *mockPixivClient) GetBookmarks(ctx context.Context, visibility pixiv.BookmarkVisibility, page int) ([]pixiv.Artwork, int, error) {
	return m.bookmarkArtworks, m.bookmarksTotal, m.bookmarksErr
}

func (m *mockPixivClient) GetImageStream(ctx context.Context, url string) (io.ReadCloser, error) {
	return m.imageStream, m.imageErr
}

func TestService_GetBookmarks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			bookmarkArtworks: []pixiv.Artwork{
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
			bookmarksTotal: 1,
		}

		svc := NewService(ServiceConfig{
			Client:           mock,
			BackendPublicURL: "https://proxy.example.com",
		})

		bookmarks, meta, err := svc.GetBookmarks(t.Context(), pixiv.BookmarkVisibilityPublic, 1)

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
			bookmarksErr: errors.New("network fail"),
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		_, _, err := svc.GetBookmarks(t.Context(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want error")
		}
	})

	t.Run("Invalid Artwork Skipped", func(t *testing.T) {
		mock := &mockPixivClient{
			bookmarkArtworks: []pixiv.Artwork{
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
			bookmarksTotal: 2,
		}

		svc := NewService(ServiceConfig{
			Client:           mock,
			BackendPublicURL: "https://proxy.example.com",
		})

		bookmarks, _, err := svc.GetBookmarks(t.Context(), pixiv.BookmarkVisibilityPublic, 1)

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
			bookmarkArtworks: []pixiv.Artwork{
				{
					ID:             "1",
					URL:            "https://i.pximg.net/img-master/test.jpg",
					UserID:         "2",
					IsBookmarkable: true,
					Width:          800,
					Height:         600,
				},
			},
			bookmarksTotal: 1,
		}

		svc := NewService(ServiceConfig{
			Client:           mock,
			BackendPublicURL: "https://proxy.example.com",
		})

		bookmarks, err := svc.GetAllBookmarks(t.Context(), pixiv.BookmarkVisibilityPublic)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if len(bookmarks) != 1 {
			t.Errorf("got %d bookmarks, want 1", len(bookmarks))
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mock := &mockPixivClient{
			bookmarksErr: errors.New("network fail"),
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		_, err := svc.GetAllBookmarks(t.Context(), pixiv.BookmarkVisibilityPublic)

		if err == nil {
			t.Error("want error")
		}
	})
}

func TestService_GetImageProxy(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			imageStream: io.NopCloser(strings.NewReader("image-data")),
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		stream, err := svc.GetImageProxy(t.Context(), "test.jpg")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		defer func() {
			if err := stream.Close(); err != nil {
				t.Errorf("failed to close stream: %v", err)
			}
		}()

		data, _ := io.ReadAll(stream)

		if string(data) != "image-data" {
			t.Errorf("got %s, want image-data", string(data))
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		mock := &mockPixivClient{
			imageErr: pixiv.ErrPixivInvalidURL,
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		_, err := svc.GetImageProxy(t.Context(), "test.jpg")

		if err == nil {
			t.Fatal("want error for invalid url")
		}

		var httpErr *api.HTTPError
		if !errors.As(err, &httpErr) {
			t.Fatalf("want *api.HTTPError, got %T", err)
		}

		if httpErr.StatusCode != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", httpErr.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mock := &mockPixivClient{
			imageErr: errors.New("network fail"),
		}

		svc := NewService(ServiceConfig{
			Client: mock,
		})

		_, err := svc.GetImageProxy(t.Context(), "test.jpg")

		if err == nil {
			t.Error("want error")
		}
	})
}
