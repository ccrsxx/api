package pixiv

import (
	"errors"
	"strings"

	"math"

	pClient "github.com/ccrsxx/api/internal/clients/pixiv"
	"github.com/ccrsxx/api/internal/model"
)

func parseArtworkToBookmark(artwork pClient.Artwork, pixivImageURL string) (model.Bookmark, error) {
	if !artwork.IsBookmarkable {
		return model.Bookmark{}, errors.New("artwork is not bookmarkable")
	}

	imageURL := artwork.URL

	// Remove existing size variant segment (e.g., /c/250x250_80_a2)
	if start := strings.Index(imageURL, "/c/"); start != -1 {
		if end := strings.Index(imageURL[start+3:], "/"); end != -1 {
			imageURL = imageURL[:start] + imageURL[start+3+end:]
		}
	}

	// Normalize path and thumbnail suffix to base img-master
	imageURL = strings.Replace(imageURL, "/custom-thumb/", "/img-master/", 1)
	imageURL = strings.Replace(imageURL, "_custom1200", "_master1200", 1)
	imageURL = strings.Replace(imageURL, "_square1200", "_master1200", 1)

	// Inject WebP transformation block
	imageURL = strings.Replace(imageURL, "/img-master/", "/c/1200x1200_90_webp/img-master/", 1)

	// Rewrite to proxy
	imageURL = strings.Replace(imageURL, "https://", pixivImageURL+"/", 1)

	// Image dimensions
	width, height := calculateMaster1200Dimensions(artwork.Width, artwork.Height)

	// Others
	pixivURL := "https://pixiv.net/artworks/" + string(artwork.ID)
	aiGenerated := artwork.AIType == pClient.AIGenerated

	return model.Bookmark{
		ID:          string(artwork.ID),
		Title:       artwork.Title,
		ImageURL:    imageURL,
		PixivURL:    pixivURL,
		ArtistID:    string(artwork.UserID),
		ArtistName:  artwork.UserName,
		Width:       width,
		Height:      height,
		Tags:        artwork.Tags,
		AiGenerated: aiGenerated,
		CreatedAt:   artwork.CreatedAt,
	}, nil
}

// Scales dimensions down so the longest side is 1200.
func calculateMaster1200Dimensions(originalWidth, originalHeight int) (int, int) {
	const maxDimension = 1200

	// If neither side exceeds 1200, no downscaling is needed.
	if originalWidth <= maxDimension && originalHeight <= maxDimension {
		return originalWidth, originalHeight
	}

	if originalWidth > originalHeight {
		// Width is the longest side, clamp to 1200 and scale height
		ratio := maxDimension / float64(originalWidth)
		newHeight := int(math.Ceil(float64(originalHeight) * ratio))

		return maxDimension, newHeight
	}

	// Height is the longest side (or a perfect square), clamp to 1200 and scale width
	ratio := maxDimension / float64(originalHeight)
	newWidth := int(math.Ceil(float64(originalWidth) * ratio))

	return newWidth, maxDimension
}
