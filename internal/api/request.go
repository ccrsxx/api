package api

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON[T any](r *http.Request, v *T) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &HTTPError{
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
