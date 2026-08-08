package logger

import "context"

type contextKey struct{}

var requestInfoKey = contextKey{}

// RequestInfo carries per-request correlation identifiers through the request
// context so they can be attached to every log line and, later, propagated to
// traces.
type RequestInfo struct {
	RequestID string
	TraceID   string
	SpanID    string
	Route     string
}

// WithRequestInfo returns a copy of ctx that carries the given RequestInfo.
func WithRequestInfo(ctx context.Context, info RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfoKey, info)
}

// RequestInfoFromContext extracts RequestInfo from ctx, if present.
func RequestInfoFromContext(ctx context.Context) (RequestInfo, bool) {
	info, ok := ctx.Value(requestInfoKey).(RequestInfo)
	return info, ok
}
