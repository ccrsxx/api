package auth

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/config"
)

type middleware struct{}

var Middleware = &middleware{}

func (m *middleware) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")

		token, err := Service.getAuthorizationFromBearerToken(r.Context(), headerToken)

		if err != nil {
			// No need to wrap error, service already return proper error
			api.HandleHttpError(w, r, err)
			return
		}

		isValidSecretKey := config.Env().SecretKey == token

		if !isValidSecretKey {
			// No need to wrap error, service already return proper error
			api.HandleHttpError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (m *middleware) IsAuthorizedFromQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryToken := r.URL.Query().Get("token")

		token, err := Service.getAuthorizationFromQuery(r.Context(), queryToken)

		if err != nil {
			// No need to wrap error, service already return proper error
			api.HandleHttpError(w, r, err)
			return
		}

		isValidSecretKey := config.Env().SecretKey == token

		if !isValidSecretKey {
			// No need to wrap error, service already return proper error
			api.HandleHttpError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
