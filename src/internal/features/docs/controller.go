package docs

import (
	_ "embed"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/utils"
)

//go:embed openapi.json
var openapiSpec []byte

func getDocs(w http.ResponseWriter, r *http.Request) error {
	serverOverride := scalargo.ServerOverride{
		URL:         utils.GetPublicUrlFromRequest(r),
		Description: "Production server",
	}

	html, err := scalargo.NewV2(
		scalargo.WithTheme(scalargo.ThemeDefault),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithServers(serverOverride),
		scalargo.WithSpecBytes(openapiSpec),
	)

	if err != nil {
		return api.NewHttpError(http.StatusInternalServerError, "docs render failure")
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(html)); err != nil {
		return api.NewHttpError(http.StatusInternalServerError, "docs send failure")
	}

	return nil
}
