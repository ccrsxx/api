package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Registry is the application's Prometheus registry.
//
// A dedicated registry (instead of prometheus.DefaultRegisterer) keeps the
// exposed metrics explicit and avoids surprises from libraries that register
// into the global default.
var Registry = prometheus.NewRegistry()

// HTTP server metrics.
//
// Labels are deliberately low cardinality: route is the matched mux pattern
// (never the raw path), so /contents/{slug} is one series instead of one per
// slug.
var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests handled.",
		},
		[]string{"route", "method", "status_code"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method"},
	)
)

// Application metrics.
//
// These surface behaviour that is currently invisible: how many SSE clients are
// connected, whether the cache is actually helping, and which upstream
// (Spotify/Jellyfin/Navidrome) is slow. The upstream histogram matters most
// because the SSE service silently falls back to a default payload on error,
// making an outage indistinguishable from "nothing is playing".
var (
	SSEActiveClients = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sse_active_clients",
			Help: "Number of currently connected SSE clients.",
		},
	)

	CacheOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Cache lookups by result (hit or miss).",
		},
		[]string{"result"},
	)

	UpstreamRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "upstream_request_duration_seconds",
			Help:    "Upstream provider request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"platform", "result"},
	)
)

func init() {
	Registry.MustRegister(
		// Go runtime: go_goroutines, heap, GC. go_goroutines is the key leak
		// signal for this app because every SSE client parks a goroutine and
		// cache writes are fire-and-forget.
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		HTTPRequestsTotal,
		HTTPRequestDuration,

		SSEActiveClients,
		CacheOperationsTotal,
		UpstreamRequestDuration,
	)
}
