package auth

import (
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type middleware struct{}

var Middleware = &middleware{}

func (m *middleware) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")

		_, err := Service.getAuthorizationFromBearerToken(r.Context(), headerToken)

		if err != nil {
			api.HandleHttpError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) IsAuthorizedFromQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryToken := r.URL.Query().Get("token")

		_, err := Service.getAuthorizationFromQuery(r.Context(), queryToken)

		if err != nil {
			api.HandleHttpError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
