package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/utils"
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

	ipAddress := utils.GetIpAddressFromRequest(r)

	var httpErr *HttpError

	if errors.As(err, &httpErr) {
		slog.Error("http error handled",
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
			logErrorResponse(errorId, err)
		}

		return
	}

	slog.Error("http error unhandled",
		"error", err,
		"error_id", errorId,
		"path", r.URL.Path,
		"method", r.Method,
		"ip_address", ipAddress,
	)

	if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorId); err != nil {
		logErrorResponse(errorId, err)
	}
}

func logErrorResponse(errorId string, err error) {
	slog.Error("send error response failed", "error_id", errorId, "error", err)
}
