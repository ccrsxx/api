package pixiv

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

const pixivImageHostname = "i.pximg.net"

// getValidPixivImageURL converts a scheme-less pixiv image reference into a
// secure, canonical HTTPS URL pointing at i.pximg.net.
func getValidPixivImageURL(imageURL string) (string, error) {
	// Enforce scheme-less input contract.
	lowerURL := strings.ToLower(imageURL)

	if strings.HasPrefix(lowerURL, "http://") || strings.HasPrefix(lowerURL, "https://") {
		return "", errors.New("pixiv image url must not include a scheme")
	}

	// Parse with a forced scheme and strip all leading slashes.
	parsed, err := url.Parse("https://" + strings.TrimLeft(imageURL, "/"))

	if err != nil {
		return "", errors.New("pixiv image invalid url")
	}

	// Enforce strict host matching (no credentials, no ports).
	if parsed.User != nil {
		return "", errors.New("pixiv image url must not include credentials")
	}

	if parsed.Port() != "" {
		return "", errors.New("pixiv image url must not include a port")
	}

	if !strings.EqualFold(parsed.Hostname(), pixivImageHostname) {
		return "", errors.New("pixiv image invalid host")
	}

	// Validate the escaped path to prevent encoding or traversal bypasses.
	escaped := parsed.EscapedPath()

	if strings.Contains(escaped, "%") {
		return "", errors.New("pixiv image url must not be percent-encoded")
	}

	if escaped != path.Clean(escaped) {
		return "", errors.New("pixiv image invalid path")
	}

	// Enforce allowed extensions.
	switch strings.ToLower(path.Ext(escaped)) {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		// allowed
	default:
		return "", errors.New("pixiv image invalid extension")
	}

	// Rebuild from scratch using the hardcoded host and validated path.
	out := url.URL{Scheme: "https", Host: pixivImageHostname, Path: parsed.Path}

	return out.String(), nil
}
