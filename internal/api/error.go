package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/utils"
	"github.com/google/uuid"
)

var isDevelopmentMode bool

func Load(isDevelopment bool) {
	isDevelopmentMode = isDevelopment
}

type PanicError struct {
	Value   any
	Stack   string
	Message string
}

func (e *PanicError) Error() string {
	return e.Message
}

type HTTPError struct {
	Message    string
	Details    []string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func HandleHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	errorID := uuid.New().String()

	ipAddress := utils.GetIPAddressFromRequest(r)

	if panicErr, ok := errors.AsType[*PanicError](err); ok {
		parsedStack := panicErr.Stack

		if isDevelopmentMode {
			parsedStack = "disabled in development mode"

			fmt.Printf("panic stack trace:\n%s\n", panicErr.Stack)
		}

		slog.Error("http panic error",
			"message", panicErr.Message,
			"value", panicErr.Value,
			"stack", parsedStack,
			"error", err,
			"error_id", errorID,
			"path", r.URL.Path,
			"method", r.Method,
			"ip_address", ipAddress,
		)

		if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorID); err != nil {
			logErrorResponse(err, errorID)
		}

		return
	}

	if httpErr, ok := errors.AsType[*HTTPError](err); ok {
		slog.Error("http handled error",
			"message", httpErr.Message,
			"status_code", httpErr.StatusCode,
			"details", httpErr.Details,
			"error", err,
			"error_id", errorID,
			"path", r.URL.Path,
			"method", r.Method,
			"ip_address", ipAddress,
		)

		if err := NewErrorResponse(w, httpErr.StatusCode, httpErr.Message, httpErr.Details, errorID); err != nil {
			logErrorResponse(err, errorID)
		}

		return
	}

	// Any unhandled errors

	slog.Error("http unhandled error",
		"error", err,
		"error_id", errorID,
		"path", r.URL.Path,
		"method", r.Method,
		"ip_address", ipAddress,
	)

	if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorID); err != nil {
		logErrorResponse(err, errorID)
	}
}

func logErrorResponse(err error, errorID string) {
	slog.Error("send error response failed", "error", err, "error_id", errorID)
}
