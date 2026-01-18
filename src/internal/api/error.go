package api

import (
	"log"
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

func logErrorResponse(errorId string, err error) {
	log.Printf("Failed to send error response for error id %s: %v", errorId, err)
}

func HandleHttpError(w http.ResponseWriter, r *http.Request, err error) {
	errorId := uuid.New().String()

	if apiErr, ok := err.(*HttpError); ok {
		log.Printf("Handled error with id %s: %s", errorId, apiErr.Message)

		if err := NewErrorResponse(w, apiErr.StatusCode, apiErr.Message, nil, errorId); err != nil {
			logErrorResponse(errorId, err)
		}

		return
	}

	log.Printf("Unhandled error with id %s: %v", errorId, err)

	if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorId); err != nil {
		logErrorResponse(errorId, err)
	}
}
