package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/utils"
	"github.com/google/uuid"
)

type PanicError struct {
	Value   any
	Stack   string
	Message string
}

func (e *PanicError) Error() string {
	return e.Message
}

type HttpError struct {
	Message    string
	Details    []string
	StatusCode int
}

func (e *HttpError) Error() string {
	return e.Message
}

func HandleHttpError(w http.ResponseWriter, r *http.Request, err error) {
	errorId := uuid.New().String()

	ipAddress := utils.GetIpAddressFromRequest(r)

	if panicErr, ok := errors.AsType[*PanicError](err); ok {
		parsedStack := panicErr.Stack

		if config.Config().IsDevelopment {
			parsedStack = "disabled in development mode"

			fmt.Printf("panic stack trace:\n%s\n", panicErr.Stack)
		}

		slog.Error("http panic error",
			"message", panicErr.Message,
			"value", panicErr.Value,
			"stack", parsedStack,
			"error", err,
			"error_id", errorId,
			"path", r.URL.Path,
			"method", r.Method,
			"ip_address", ipAddress,
		)

		if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorId); err != nil {
			logErrorResponse(err, errorId)
		}

		return
	}

	if httpErr, ok := errors.AsType[*HttpError](err); ok {
		slog.Error("http handled error",
			"message", httpErr.Message,
			"status_code", httpErr.StatusCode,
			"details", httpErr.Details,
			"error", err,
			"error_id", errorId,
			"path", r.URL.Path,
			"method", r.Method,
			"ip_address", ipAddress,
		)

		if err := NewErrorResponse(w, httpErr.StatusCode, httpErr.Message, httpErr.Details, errorId); err != nil {
			logErrorResponse(err, errorId)
		}

		return
	}

	// Any unhandled errors

	slog.Error("http unhandled error",
		"error", err,
		"error_id", errorId,
		"path", r.URL.Path,
		"method", r.Method,
		"ip_address", ipAddress,
	)

	if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorId); err != nil {
		logErrorResponse(err, errorId)
	}
}

func logErrorResponse(err error, errorId string) {
	slog.Error("send error response failed", "error", err, "error_id", errorId)
}
