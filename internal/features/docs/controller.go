package docs

import (
	"fmt"
	"log/slog"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/utils"
)

type Controller struct {
	openapiSpec []byte
}

func NewController(openapiSpec []byte) *Controller {
	return &Controller{
		openapiSpec: openapiSpec,
	}
}

func (c *Controller) getDocs(w http.ResponseWriter, r *http.Request) {
	serverOverride := scalargo.ServerOverride{
		URL:         utils.GetPublicUrlFromRequest(r),
		Description: "Production server",
	}

	html, err := scalargo.NewV2(
		scalargo.WithTheme(scalargo.ThemeDefault),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithServers(serverOverride),
		scalargo.WithSpecBytes(c.openapiSpec),
	)

	if err != nil {
		api.HandleHttpError(w, r, fmt.Errorf("docs render error: %w", err))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(html)); err != nil {
		slog.Warn("docs response error", "error", err)
	}
}
