package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorObject struct {
	ID      string   `json:"id"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type ErrorResponse struct {
	Error ErrorObject `json:"error"`
}

func newResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(v)

	if err != nil {
		id := uuid.New().String()

		fallback := ErrorResponse{
			Error: ErrorObject{
				ID:      id,
				Message: "An internal server error occurred",
				Details: []string{},
			},
		}

		fallbackData, _ := json.Marshal(fallback)

		w.WriteHeader(http.StatusInternalServerError)

		errMarshal := fmt.Errorf("failed to marshal response: %w", err)

		if _, writeErr := w.Write(fallbackData); writeErr != nil {
			errFallback := fmt.Errorf("failed to marshal fallback response: %w", writeErr)
			return errors.Join(errMarshal, errFallback)
		}

		return errMarshal
	}

	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func NewSuccessResponse[T any](w http.ResponseWriter, statusCode int, data T) error {
	var response any = data

	val := reflect.ValueOf(data)
	kind := val.Kind()

	if kind == reflect.Ptr && !val.IsNil() {
		kind = val.Elem().Kind()
	}

	if kind == reflect.Struct {
		response = &SuccessResponse[T]{
			Data: data,
		}
	}

	return newResponse(w, statusCode, response)
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, message string, details []string, id string) error {
	if details == nil {
		details = []string{}
	}

	return newResponse(w, statusCode, ErrorResponse{
		Error: ErrorObject{
			ID:      id,
			Message: message,
			Details: details,
		},
	})
}
