package api

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &HTTPError{
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
