package pixiv_test

import (
	"context"
	"encoding/json/v2"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	pClient "github.com/ccrsxx/api/internal/clients/pixiv"
	"github.com/ccrsxx/api/internal/features/pixiv"
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
			Client:           mock,
			BackendPublicURL: "https://proxy.example.com",
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?page=1&type=private", nil)

		ctrl.GetBookmarks(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res api.SuccessPaginatedResponse[[]pixiv.Bookmark, utils.OffsetPaginationMeta]

		if err := json.UnmarshalRead(w.Body, &res); err != nil {
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
			Client:           mock,
			BackendPublicURL: "https://proxy.example.com",
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?type=private", nil)

		ctrl.GetAllBookmarks(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res api.SuccessResponse[[]pixiv.Bookmark]

		if err := json.UnmarshalRead(w.Body, &res); err != nil {
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

func TestController_GetImageProxy(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockPixivClient{
			imageStream: io.NopCloser(strings.NewReader("image-data")),
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/test.jpg", nil)

		r.SetPathValue("url", "test.jpg")

		ctrl.GetImageProxy(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		if !strings.Contains(w.Header().Get("Cache-Control"), "max-age=31536000") {
			t.Error("want aggressive cache header")
		}

		if w.Body.String() != "image-data" {
			t.Errorf("got %q, want %q", w.Body.String(), "image-data")
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mock := &mockPixivClient{
			imageErr: errors.New("fail"),
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/test.jpg", nil)

		r.SetPathValue("url", "test.jpg")

		ctrl.GetImageProxy(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		mock := &mockPixivClient{
			imageStream: &test.ErrorBodyCloser{Reader: strings.NewReader("image-data")},
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/test.jpg", nil)

		r.SetPathValue("url", "test.jpg")

		ctrl.GetImageProxy(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mockServer := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("image-data"))
		}))

		resp, err := mockServer.Client().Get(mockServer.URL)

		if err != nil {
			t.Fatalf("failed to create stream: %v", err)
		}

		mock := &mockPixivClient{
			imageStream: resp.Body,
		}

		svc := pixiv.NewService(pixiv.ServiceConfig{
			Client: mock,
		})

		ctrl := pixiv.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/image/test.jpg", nil)

		r.SetPathValue("url", "test.jpg")

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetImageProxy(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
