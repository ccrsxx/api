package favicon

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed favicon.ico
var icon []byte

func getFavicon(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(icon); err != nil {
		return fmt.Errorf("favicon response error: %w", err)
	}

	return nil
}
