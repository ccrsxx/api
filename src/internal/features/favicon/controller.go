package favicon

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
)

var Controller = &controller{}

type controller struct{}

//go:embed favicon.ico
var icon []byte

func (c *controller) getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(icon); err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("favicon response error: %w", err))
		return
	}
}
