package pixiv

import (
	"errors"
	"net/url"
	"path"
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

// getValidPixivImageURL takes a scheme-less pixiv image reference (e.g.
// "i.pximg.net/img-master/img/2024/01/01/00/00/00/12345678_p0_master1200.jpg")
// and returns a canonical HTTPS URL guaranteed to point at i.pximg.net.
//
// The result is rebuilt from validated parts only -- never echoed back from the
// caller's input -- so anything not explicitly checked below (query, fragment,
// port, credentials) is dropped by construction.
func getValidPixivImageURL(imageURL string) (string, error) {
	const host = "i.pximg.net"

	// Caller contract is scheme-less input. Lowered so "HtTpS://" can't slip by.
	if lower := strings.ToLower(imageURL); strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") {
		return "", errors.New("pixiv image url must not include a scheme")
	}

	// Trim *all* leading slashes: stops "//evil.com/a.jpg" from surviving as an
	// authority, and "/i.pximg.net/a.jpg" from parsing into an empty host.
	parsed, err := url.Parse("https://" + strings.TrimLeft(imageURL, "/"))

	if err != nil {
		return "", errors.New("pixiv image invalid url")
	}

	switch {
	// "foo:bar@i.pximg.net/a.jpg" passes a host check but leaks Basic auth.
	case parsed.User != nil:
		return "", errors.New("pixiv image url must not include credentials")

	// Hostname() strips the port, so "i.pximg.net:8080" would compare equal.
	case parsed.Port() != "":
		return "", errors.New("pixiv image url must not include a port")

	// EqualFold: url.Parse lowercases the scheme but never the host.
	case !strings.EqualFold(parsed.Hostname(), host):
		return "", errors.New("pixiv image invalid host")
	}

	// Validate the escaped path -- the bytes actually sent on the wire.
	// parsed.Path is decoded, so "evil.svg%3f.jpg" would fake a ".jpg" there.
	escaped := parsed.EscapedPath()

	switch {
	// Real pximg paths never need encoding, so ban it outright instead of
	// reasoning about what %2e / %2f / %00 decode to. Also catches raw spaces.
	case strings.Contains(escaped, "%"):
		return "", errors.New("pixiv image url must not be percent-encoded")

	// Reject non-canonical paths ("/a/../b.jpg", or an empty path) rather than
	// rewriting them -- a rewrite means fetching something we never inspected.
	case escaped != path.Clean(escaped):
		return "", errors.New("pixiv image invalid path")
	}

	switch strings.ToLower(path.Ext(escaped)) {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		// allowed
	default:
		return "", errors.New("pixiv image invalid extension")
	}

	// Host is the constant, not parsed.Host, so nothing upstream can redirect it.
	out := url.URL{Scheme: "https", Host: host, Path: parsed.Path}

	return out.String(), nil
}
