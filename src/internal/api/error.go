package api

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(statusCode int, message string) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func HandleHttpError(w http.ResponseWriter, r *http.Request, err error) {
	errorId := uuid.New().String()

	if apiErr, ok := err.(*HttpError); ok {
		slog.Error("handled error",
			"error_id", errorId,
			"message", apiErr.Message,
			"status_code", apiErr.StatusCode,
			"method", r.Method,
			"path", r.URL.Path,
		)

		if err := NewErrorResponse(w, apiErr.StatusCode, apiErr.Message, nil, errorId); err != nil {
			logErrorResponse(errorId, err)
		}

		return
	}

	slog.Error("unhandled error",
		"error_id", errorId,
		"error", err,
		"method", r.Method,
		"path", r.URL.Path,
	)

	if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorId); err != nil {
		logErrorResponse(errorId, err)
	}
}

func logErrorResponse(errorId string, err error) {
	slog.Error("send error response failed", "error_id", errorId, "error", err)
}
