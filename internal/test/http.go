package test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
)

// Mocks a broken network connection (always fails Write).
type ErrorResponseRecorder struct {
	*httptest.ResponseRecorder
}

func (w *ErrorResponseRecorder) Write(b []byte) (int, error) {
	return 0, errors.New("forced write error")
}

// Mocks a broken network connection (always fails Write) for http.ResponseWriter interface.
type ErrorResponseWriter struct {
	http.ResponseWriter
}

func (w *ErrorResponseWriter) Write(b []byte) (int, error) {
	return 0, errors.New("forced write error")
}

// Mocks a ResponseWriter that does NOT implement http.Flusher, to test flush errors.
type NonFlusherResponseWriter struct {
	http.ResponseWriter
}

// Mocks a response body that fails on Close (e.g., network error during close).
type ErrorBodyCloser struct {
	io.Reader
}

func (e *ErrorBodyCloser) Close() error {
	return errors.New("forced close error")
}

// Allows defining a custom RoundTrip function for http.Client, useful for testing without a real server.
type CustomTransport func(req *http.Request) (*http.Response, error)

func (f CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
