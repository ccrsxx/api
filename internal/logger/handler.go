package logger

import (
	"context"
	"log/slog"
)

// ContextHandler is a slog.Handler middleware that enriches every record with
// correlation identifiers pulled from the request context.
//
// This is what makes existing slog calls across every feature package become
// correlated for free: as long as the call uses a *Context variant
// (slog.InfoContext, etc.) the trace_id/request_id/route are appended without
// threading a logger through function signatures.
//
// High-cardinality identifiers (trace_id, request_id) are emitted as regular
// attributes here; the log collector (Alloy) promotes them to Loki structured
// metadata rather than indexed labels.
type ContextHandler struct {
	slog.Handler
}

// NewContextHandler wraps an existing handler.
func NewContextHandler(next slog.Handler) *ContextHandler {
	return &ContextHandler{Handler: next}
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if info, ok := RequestInfoFromContext(ctx); ok {
		if info.TraceID != "" {
			r.AddAttrs(slog.String("trace_id", info.TraceID))
		}

		if info.SpanID != "" {
			r.AddAttrs(slog.String("span_id", info.SpanID))
		}

		if info.RequestID != "" {
			r.AddAttrs(slog.String("request_id", info.RequestID))
		}

		if info.Route != "" {
			r.AddAttrs(slog.String("route", info.Route))
		}
	}

	return h.Handler.Handle(ctx, r)
}

// WithAttrs and WithGroup must return a *ContextHandler so the enrichment is
// preserved across derived loggers (e.g. slog.With(...)).
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithGroup(name)}
}
