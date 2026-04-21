package api

import (
	"encoding/json"
	"net/http"

	"github.com/ccrsxx/api/internal/utils"
)

func DecodeJSON[T any](r *http.Request, v *T) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &HTTPError{
			Message:    "Invalid body",
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := utils.Validate.Struct(v); err != nil {
		_, details := utils.FormatValidationError(err)
		return &HTTPError{
			Message:    "Invalid body",
			Details:    details,
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
