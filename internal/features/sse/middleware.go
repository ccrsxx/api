package sse

import (
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/utils"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(svc *Service) *Middleware {
	return &Middleware{
		service: svc,
	}
}

func (m *Middleware) IsConnectionAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := utils.GetIPAddressFromRequest(r)

		if err := m.service.IsConnectionAllowed(ipAddress); err != nil {
			api.HandleHTTPError(w, r, err)
			return
		}

		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		next.ServeHTTP(w, r)
	})
}
