package sse

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/utils"
)

type controller struct {
	service *service
}

func NewController(svc *service) *controller {
	return &controller{
		service: svc,
	}
}

func (c *controller) getCurrentPlayingSSE(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	ctx := r.Context()

	clientChan := make(chan string, 4)

	ipAddress := utils.GetIpAddressFromRequest(r)
	userAgent := r.UserAgent()

	c.service.AddClient(ctx, clientChan, ipAddress, userAgent)

	defer c.service.RemoveClient(ctx, clientChan)

	clientDisconnected := r.Context().Done()

	for {
		select {
		case <-clientDisconnected:
			slog.Info("sse client disconnected", "ip", ipAddress)
			return
		case msg, ok := <-clientChan:
			if !ok {
				slog.Info("sse client channel closed", "ip", ipAddress)
				return
			}

			if _, err := fmt.Fprint(w, msg); err != nil {
				slog.Warn("sse write error", "error", err)
				return
			}

			if err := rc.Flush(); err != nil {
				slog.Warn("sse flush error", "error", err)
			}
		}
	}
}
