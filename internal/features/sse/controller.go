package sse

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/utils"
)

type Controller struct {
	ctx     context.Context
	service *Service
}

func NewController(ctx context.Context, svc *Service) *Controller {
	return &Controller{
		ctx:     ctx,
		service: svc,
	}
}

func (c *Controller) getCurrentPlayingSSE(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	ctx := r.Context()

	clientChan := make(chan string, 4)

	ipAddress := utils.GetIpAddressFromRequest(r)
	userAgent := r.UserAgent()

	c.service.AddClient(ctx, clientChan, ipAddress, userAgent)

	defer c.service.RemoveClient(ctx, clientChan)

	appShutdown := c.ctx.Done()

	clientDisconnected := r.Context().Done()

	for {
		select {
		case <-appShutdown:
			slog.Info("server is shutting down, closing sse connection", "ip", ipAddress)
			return
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
