package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// NewTxQuerier exposes the internal tx-scoped querier factory for testing.
func (s *Service) NewTxQuerier(tx pgx.Tx) querier {
	return s.newTx(tx)
}

// SetSignToken overrides the package-level signToken function for testing.
// It returns a restore function that should be deferred.
func SetSignToken(fn func(token *jwt.Token, key []byte) (string, error)) func() {
	original := signToken
	signToken = fn

	return func() {
		signToken = original
	}
}

// SetParseToken overrides the package-level parseToken function for testing.
// It returns a restore function that should be deferred.
func SetParseToken(fn func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, opts ...jwt.ParserOption) (*jwt.Token, error)) func() {
	original := parseToken
	parseToken = fn

	return func() {
		parseToken = original
	}
}
