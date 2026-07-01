package navidrome_test

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
	"github.com/ccrsxx/api/internal/clients/navidrome"
	navidromeFeature "github.com/ccrsxx/api/internal/features/navidrome"
	"github.com/ccrsxx/api/internal/test"
)

type mockNavidromeClient struct {
	nowPlayingResult []navidrome.NowPlayingEntry
	nowPlayingErr    error

	coverArtResult io.ReadCloser
	coverArtErr    error
}

func (m *mockNavidromeClient) GetNowPlaying(ctx context.Context) ([]navidrome.NowPlayingEntry, error) {
	return m.nowPlayingResult, m.nowPlayingErr
}

func (m *mockNavidromeClient) GetCoverArtStream(ctx context.Context, coverArtID string) (io.ReadCloser, error) {
	return m.coverArtResult, m.coverArtErr
}

func TestController_GetCurrentlyPlaying(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockNavidromeClient{}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctrl.GetCurrentlyPlaying(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res api.SuccessResponse[struct {
			IsPlaying bool `json:"isPlaying"`
		}]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Data.IsPlaying {
			t.Error("want isPlaying false")
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mock := &mockNavidromeClient{
			nowPlayingErr: errors.New("fail"),
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctrl.GetCurrentlyPlaying(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mock := &mockNavidromeClient{}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetCurrentlyPlaying(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_GetCoverArt(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockNavidromeClient{
			coverArtResult: io.NopCloser(strings.NewReader("image-data")),
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/cover-art/ca-123", nil)

		r.SetPathValue("coverArtID", "ca-123")

		ctrl.GetCoverArt(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		if w.Header().Get("Content-Type") != "image/webp" {
			t.Errorf("got %s, want image/webp", w.Header().Get("Content-Type"))
		}

		if !strings.Contains(w.Header().Get("Cache-Control"), "max-age=31536000") {
			t.Error("want aggressive cache header")
		}

		if w.Body.String() != "image-data" {
			t.Errorf("got %q, want %q", w.Body.String(), "image-data")
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mock := &mockNavidromeClient{
			coverArtErr: errors.New("fail"),
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/cover-art/ca-123", nil)

		r.SetPathValue("coverArtID", "ca-123")

		ctrl.GetCoverArt(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		mock := &mockNavidromeClient{
			coverArtResult: &test.ErrorBodyCloser{Reader: strings.NewReader("image-data")},
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/cover-art/ca-123", nil)

		r.SetPathValue("coverArtID", "ca-123")

		ctrl.GetCoverArt(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("image-data"))
		}))

		defer mockServer.Close()

		resp, err := mockServer.Client().Get(mockServer.URL)

		if err != nil {
			t.Fatalf("failed to create stream: %v", err)
		}

		mock := &mockNavidromeClient{
			coverArtResult: resp.Body,
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		ctrl := navidromeFeature.NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/cover-art/ca-123", nil)

		r.SetPathValue("coverArtID", "ca-123")

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetCoverArt(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
