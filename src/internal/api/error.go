package api

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type HttpError struct {
	StatusCode int
	Message    string
	Details    []string
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(statusCode int, message string, details []string) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

func HandleHttpError(w http.ResponseWriter, r *http.Request, err error) {
	errorId := uuid.New().String()

	if err, ok := err.(*HttpError); ok {
		slog.Error("api error handled",
			"error_id", errorId,
			"message", err.Message,
			"status_code", err.StatusCode,
			"method", r.Method,
			"path", r.URL.Path,
		)

		if err := NewErrorResponse(w, err.StatusCode, err.Message, err.Details, errorId); err != nil {
			logErrorResponse(errorId, err)
		}

		return
	}

	slog.Error("api error unhandled",
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
