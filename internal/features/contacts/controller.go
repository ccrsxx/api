package contacts

import (
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) CreateContact(w http.ResponseWriter, r *http.Request) {
	var input CreateContactInput

	if err := api.DecodeJSON(r, &input); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	ipAddress := utils.GetIPAddressFromRequest(r)

	if err := c.service.cloudflareClient.VerifyTurnstile(r.Context(), input.Token, ipAddress); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	if err := c.service.CreateContact(r.Context(), input, ipAddress); err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
