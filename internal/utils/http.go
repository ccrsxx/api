package utils

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHTTPClient returns an HTTP client instrumented for distributed tracing.
//
// Every outbound call gets its own span and propagates W3C trace context, so a
// slow third party (Spotify, Jellyfin, Navidrome, ...) appears as a labelled
// span in the request's trace waterfall instead of unexplained latency.
//
// When tracing is disabled the global tracer provider is a no-op, so this
// behaves like a plain http.Client.
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}
