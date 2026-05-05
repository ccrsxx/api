package auth

import (
	"net/http"

	"github.com/ccrsxx/api/internal/api"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(svc *Service) *Middleware {
	return &Middleware{
		service: svc,
	}
}

func (m *Middleware) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")

		_, err := m.service.GetAuthorizationFromBearerToken(r.Context(), headerToken)

		if err != nil {
			api.HandleHTTPError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) IsAuthorizedFromQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryToken := r.URL.Query().Get("token")

		_, err := m.service.GetAuthorizationFromQuery(r.Context(), queryToken)

		if err != nil {
			api.HandleHTTPError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) IsAuthorizedFromOauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := m.service.ValidateOauthToken(ctx, r)

		if err != nil {
			api.HandleHTTPError(w, r, err)
			return
		}

		ctxWithUser := SetUserContext(ctx, user)

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

func (m *Middleware) IsAdminFromOauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.service.IsAdminFromOauth(r.Context())

		if err != nil {
			api.HandleHTTPError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
