package pixiv

import (
	"testing"
	"time"

	pClient "github.com/ccrsxx/api/internal/clients/pixiv"
)

func Test_parseArtworkToBookmark(t *testing.T) {
	mockPixivImageURL := "https://proxy.example.com"

	t.Run("Full Data", func(t *testing.T) {
		artwork := pClient.Artwork{
			ID:             "12345",
			URL:            "https://i.pximg.net/c/250x250_80_a2/img-master/img/2024/01/01/00/00/00/12345_p0_square1200.jpg",
			Title:          "Test Artwork",
			UserID:         "67890",
			UserName:       "TestArtist",
			Pages:          1,
			AIType:         pClient.AIGenerated,
			IsBookmarkable: true,
			Width:          2000,
			Height:         1000,
			Tags:           []string{"tag1", "tag2"},
			CreatedAt:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		got, err := parseArtworkToBookmark(artwork, mockPixivImageURL)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got.ID != "12345" {
			t.Errorf("got %s, want 12345", got.ID)
		}

		if got.Title != "Test Artwork" {
			t.Errorf("got %s, want Test Artwork", got.Title)
		}

		if got.ArtistID != "67890" {
			t.Errorf("got %s, want 67890", got.ArtistID)
		}

		if got.ArtistName != "TestArtist" {
			t.Errorf("got %s, want TestArtist", got.ArtistName)
		}

		if got.PixivURL != "https://pixiv.net/artworks/12345" {
			t.Errorf("got %s, want pixiv url", got.PixivURL)
		}

		if !got.AiGenerated {
			t.Error("want aiGenerated true")
		}

		if len(got.Tags) != 2 {
			t.Errorf("got %d tags, want 2", len(got.Tags))
		}

		// Width > Height, so width gets clamped to 1200, height scales proportionally
		if got.Width != 1200 {
			t.Errorf("got width %d, want 1200", got.Width)
		}

		if got.Height != 600 {
			t.Errorf("got height %d, want 600", got.Height)
		}
	})

	t.Run("Not Bookmarkable", func(t *testing.T) {
		artwork := pClient.Artwork{
			IsBookmarkable: false,
		}

		_, err := parseArtworkToBookmark(artwork, mockPixivImageURL)

		if err == nil {
			t.Error("want error for non-bookmarkable artwork")
		}
	})

	t.Run("Custom Thumb URL Normalization", func(t *testing.T) {
		artwork := pClient.Artwork{
			ID:             "111",
			URL:            "https://i.pximg.net/c/250x250_80_a2/custom-thumb/img/2024/01/01/00/00/00/111_p0_custom1200.jpg",
			UserID:         "222",
			IsBookmarkable: true,
			Width:          800,
			Height:         600,
		}

		got, err := parseArtworkToBookmark(artwork, mockPixivImageURL)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		// Should contain img-master and master1200, not custom-thumb/custom1200
		wantPrefix := mockPixivImageURL + "/i.pximg.net/c/1200x1200_90_webp/img-master/"

		// Only check prefix since the rest of the URL is normalized and can vary
		if got.ImageURL[:len(wantPrefix)] != wantPrefix {
			t.Errorf("got %s, want prefix %s", got.ImageURL, wantPrefix)
		}
	})

	t.Run("Not AI Generated", func(t *testing.T) {
		artwork := pClient.Artwork{
			ID:             "111",
			URL:            "https://i.pximg.net/img-master/test.jpg",
			UserID:         "222",
			IsBookmarkable: true,
			AIType:         pClient.NotAIGenerated,
			Width:          800,
			Height:         600,
		}

		got, err := parseArtworkToBookmark(artwork, mockPixivImageURL)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got.AiGenerated {
			t.Error("want aiGenerated false")
		}
	})
}

func Test_calculateMaster1200Dimensions(t *testing.T) {
	t.Run("No Downscale Needed", func(t *testing.T) {
		w, h := calculateMaster1200Dimensions(800, 600)

		if w != 800 || h != 600 {
			t.Errorf("got %dx%d, want 800x600", w, h)
		}
	})

	t.Run("Width Is Longest", func(t *testing.T) {
		w, h := calculateMaster1200Dimensions(2400, 1200)

		if w != 1200 {
			t.Errorf("got width %d, want 1200", w)
		}

		if h != 600 {
			t.Errorf("got height %d, want 600", h)
		}
	})

	t.Run("Height Is Longest", func(t *testing.T) {
		w, h := calculateMaster1200Dimensions(800, 1600)

		if w != 600 {
			t.Errorf("got width %d, want 600", w)
		}

		if h != 1200 {
			t.Errorf("got height %d, want 1200", h)
		}
	})

	t.Run("Exact 1200", func(t *testing.T) {
		w, h := calculateMaster1200Dimensions(1200, 1200)

		if w != 1200 || h != 1200 {
			t.Errorf("got %dx%d, want 1200x1200", w, h)
		}
	})
}

func Test_getValidPixivImageURL(t *testing.T) {
	tests := []struct {
		name     string
		imageURL string
		want     string
		wantErr  bool
	}{
		{
			name:     "Canonical Master URL",
			imageURL: "i.pximg.net/img-master/img/2024/01/01/00/00/00/12345_p0_master1200.jpg",
			want:     "https://i.pximg.net/img-master/img/2024/01/01/00/00/00/12345_p0_master1200.jpg",
		},
		{
			name:     "Protocol Relative, Uppercase Host And Extension (Normalized)",
			imageURL: "//I.PXIMG.NET/a.PNG",
			want:     "https://i.pximg.net/a.PNG",
		},
		{
			name:     "Query And Fragment Dropped By Rebuild",
			imageURL: "i.pximg.net/a.jpg?w=200#anchor",
			want:     "https://i.pximg.net/a.jpg",
		},

		// One case per guard. Rejects short-circuit, so each needs its own input.
		{
			name:     "Empty Input",
			imageURL: "",
			wantErr:  true,
		},
		{
			name:     "Explicit Scheme (Mixed Case)",
			imageURL: "HtTpS://i.pximg.net/a.jpg",
			wantErr:  true,
		},
		{
			name:     "Control Character (Unparseable)",
			imageURL: "i.pximg.net/a.jpg\rHost: evil.com",
			wantErr:  true,
		},
		{
			name:     "Suffix Confusion",
			imageURL: "i.pximg.net.evil.com/a.jpg",
			wantErr:  true,
		},
		{
			name:     "Valid Host As Userinfo",
			imageURL: "i.pximg.net@evil.com/a.jpg", // real host is evil.com
			wantErr:  true,
		},
		{
			name:     "Custom Port",
			imageURL: "i.pximg.net:8080/a.jpg", // Hostname() strips the port
			wantErr:  true,
		},
		{
			name:     "Encoded Question Mark Fakes Extension",
			imageURL: "i.pximg.net/evil.svg%3f.jpg", // parsed.Path decodes to ".jpg"
			wantErr:  true,
		},
		{
			name:     "Traversal Onto Disallowed Extension",
			imageURL: "i.pximg.net/a.jpg/../evil.svg",
			wantErr:  true,
		},
		{
			name:     "Disallowed Extension",
			imageURL: "i.pximg.net/evil.svg",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValidPixivImageURL(tt.imageURL)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got %q, want error", got)
				}

				if got != "" {
					t.Errorf("got %q alongside error, want empty string", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("unwanted error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
