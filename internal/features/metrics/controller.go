package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller struct {
	handler http.Handler
}

func NewController() *Controller {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &Controller{
		handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	}
}

func (c *Controller) GetMetrics(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
}
