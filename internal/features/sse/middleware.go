package sse

import (
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/utils"
)

type middleware struct {
	service *service
}

var Middleware = &middleware{}

func NewMiddleware(svc *service) *middleware {
	return &middleware{
		service: svc,
	}
}

func (m *middleware) IsConnectionAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := utils.GetIpAddressFromRequest(r)

		if err := m.service.IsConnectionAllowed(ipAddress); err != nil {
			api.HandleHttpError(w, r, err)
			return
		}

		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		next.ServeHTTP(w, r)
	})
}
