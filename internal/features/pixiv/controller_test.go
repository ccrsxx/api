package pixiv_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	pClient "github.com/ccrsxx/api/internal/clients/pixiv"
	"github.com/ccrsxx/api/internal/features/pixiv"
	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/test"
	"github.com/ccrsxx/api/internal/utils"
)

type mockPixivClient struct {
	artworks    []pClient.Artwork
	total       int
	err         error
	imageStream io.ReadCloser
	imageErr    error
}

func (m *mockPixivClient) GetBookmarks(ctx context.Context, visibility pClient.BookmarkVisibility, page int) ([]pClient.Artwork, int, error) {
	return m.artworks, m.total, m.err
}

func (m *mockPixivClient) GetImageStream(ctx context.Context, url string) (io.ReadCloser, error) {
	return m.imageStream, m.imageErr
}

func TestController_GetBookmarks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			artworks: []pClient.Artwork{
				{
					ID:             "123",
					URL:            "https://i.pximg.net/img-master/test.jpg",
					UserID:         "456",
					IsBookmarkable: true,
					Width:          800,
					Height:         600,
				},
			},
			total: 1,
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client:        mock,
			PixivImageURL: "https://proxy.example.com",
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?page=1&type=private", nil)

		ctrl.GetBookmarks(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res api.SuccessPaginatedResponse[[]model.Bookmark, utils.OffsetPaginationMeta]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res.Data) != 1 {
			t.Errorf("got %d bookmarks, want 1", len(res.Data))
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mock := &mockPixivClient{
			err: errors.New("fail"),
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctrl.GetBookmarks(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mock := &mockPixivClient{}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetBookmarks(errWriter, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_GetAllBookmarks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			artworks: []pClient.Artwork{
				{
					ID:             "123",
					URL:            "https://i.pximg.net/img-master/test.jpg",
					UserID:         "456",
					IsBookmarkable: true,
					Width:          800,
					Height:         600,
				},
			},
			total: 1,
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client:        mock,
			PixivImageURL: "https://proxy.example.com",
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?type=private", nil)

		ctrl.GetAllBookmarks(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res api.SuccessResponse[[]model.Bookmark]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res.Data) != 1 {
			t.Errorf("got %d bookmarks, want 1", len(res.Data))
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mock := &mockPixivClient{
			err: errors.New("fail"),
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctrl.GetAllBookmarks(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mock := &mockPixivClient{}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetAllBookmarks(errWriter, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_GetImage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStream := io.NopCloser(strings.NewReader("image_data"))
		mock := &mockPixivClient{
			imageStream: mockStream,
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})
		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/i.pximg.net%2Ftest.jpg", nil)
		r.SetPathValue("url", "i.pximg.net/test.jpg")

		ctrl.GetImage(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		expectedCache := "public, immutable, no-transform, max-age=31536000"
		if w.Header().Get("Cache-Control") != expectedCache {
			t.Errorf("got Cache-Control %s, want %s", w.Header().Get("Cache-Control"), expectedCache)
		}

		body, _ := io.ReadAll(w.Body)
		if string(body) != "image_data" {
			t.Errorf("got body %s, want image_data", string(body))
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mock := &mockPixivClient{
			imageErr: errors.New("network error"),
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})
		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/i.pximg.net%2Ftest.jpg", nil)
		r.SetPathValue("url", "i.pximg.net/test.jpg")

		ctrl.GetImage(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mockStream := io.NopCloser(struct{ io.Reader }{strings.NewReader("image_data")})
		mock := &mockPixivClient{
			imageStream: mockStream,
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})
		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}
		r := httptest.NewRequest(http.MethodGet, "/image/i.pximg.net%2Ftest.jpg", nil)
		r.SetPathValue("url", "i.pximg.net/test.jpg")

		ctrl.GetImage(errWriter, r)

		// It sets OK, then tries to copy to body.
		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		mockStream := &test.ErrorBodyCloser{Reader: strings.NewReader("image_data")}
		mock := &mockPixivClient{
			imageStream: mockStream,
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})
		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/i.pximg.net%2Ftest.jpg", nil)
		r.SetPathValue("url", "i.pximg.net/test.jpg")

		ctrl.GetImage(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}
	})
}
