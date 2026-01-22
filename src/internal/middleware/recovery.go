package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/ccrsxx/api-go/src/internal/api"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				err := &api.PanicError{
					Value:   rcv,
					Stack:   string(debug.Stack()),
					Message: "internal server error",
				}

				api.HandleHttpError(w, r, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
