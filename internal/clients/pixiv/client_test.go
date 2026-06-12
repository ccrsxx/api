package pixiv_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/clients/pixiv"
	"github.com/ccrsxx/api/internal/test"
)

const emptyResponse = `{"error":false,"body":{"works":[],"total":0}}`

func TestNewClient(t *testing.T) {
	client := pixiv.NewClient(pixiv.Config{})

	if client == nil {
		t.Fatal("want client to be initialized, got nil")
	}
}

func TestClient_GetBookmarks(t *testing.T) {
	createArtworkByID := func(id string) string {
		return fmt.Sprintf(`{"id":%s,"url":"https://i.pximg.net/img-master/test.jpg","title":"Test","userId":"456","userName":"Artist","profileImageUrl":"","pageCount":1,"xRestrict":0,"sl":2,"aiType":1,"bookmarkData":null,"isBookmarkable":true,"illustType":0,"tags":[],"width":1000,"height":800,"createDate":"2024-01-01T00:00:00+09:00","updateDate":"2024-01-01T00:00:00+09:00"}`, id)
	}

	wrapResponseJSON := func(works ...string) string {
		return fmt.Sprintf(`{"error":false,"body":{"works":[%s],"total":%d}}`, strings.Join(works, ","), len(works))
	}

	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(wrapResponseJSON(createArtworkByID(`"123"`)))); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: mockServer.URL,
		})

		artworks, total, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if len(artworks) != 1 {
			t.Errorf("got %d artworks, want 1", len(artworks))
		}

		if total != 1 {
			t.Errorf("got total %d, want 1", total)
		}
	})

	t.Run("Success With Numeric ID", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(wrapResponseJSON(createArtworkByID(`12345`)))); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: mockServer.URL,
		})

		artworks, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if string(artworks[0].ID) != "12345" {
			t.Errorf("got %s, want 12345", artworks[0].ID)
		}
	})

	t.Run("Invalid FlexibleID", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(wrapResponseJSON(createArtworkByID(`true`)))); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "1_abc",
			BaseURL: mockServer.URL,
		})

		artworks, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("got err: %v, want success with skipped artwork", err)
		}

		if len(artworks) != 0 {
			t.Errorf("got %d artworks, want 0 (invalid id should cause skip)", len(artworks))
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: "http://bad\x7f",
		})

		_, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want error from new request creation")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: "http://invalid.url.local",
		})

		_, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want error from httpClient.Do")
		}
	})

	t.Run("Status Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: mockServer.URL,
		})

		_, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want status error for 401")
		}
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`invalid-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: mockServer.URL,
		})

		_, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want decode error")
		}
	})

	t.Run("API Error Response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`{"error":true,"message":"some api error"}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: mockServer.URL,
		})

		_, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err == nil {
			t.Error("want error from API error response")
		}

		if !strings.Contains(err.Error(), "some api error") {
			t.Errorf("got %v, want api error message", err)
		}
	})

	t.Run("Invalid Artwork Skipped", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			// One valid artwork and one invalid (malformed JSON in works array)
			if _, err := w.Write([]byte(wrapResponseJSON(`"invalid"`, createArtworkByID(`"123"`)))); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: mockServer.URL,
		})

		artworks, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("got err: %v, want success with skipped artwork", err)
		}

		if len(artworks) != 1 {
			t.Errorf("got %d artworks, want 1 (invalid one should be skipped)", len(artworks))
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := pixiv.NewClient(pixiv.Config{
			Token:   "456_abc",
			BaseURL: "http://example.com",
			HTTPClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader(emptyResponse)},
						Header:     make(http.Header),
					}, nil
				}),
			},
		})

		_, _, err := c.GetBookmarks(context.Background(), pixiv.BookmarkVisibilityPublic, 1)

		if err != nil {
			t.Fatalf("got: %v, want GetBookmarks to handle body close error gracefully", err)
		}
	})
}
