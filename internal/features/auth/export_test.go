package auth

import (
	"github.com/ccrsxx/api/internal/test"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// TestMockQuerier embeds test.MockQuerier and adds WithTx
// to satisfy auth's unexported querier interface.
type TestMockQuerier struct {
	test.MockQuerier
}

func (m *TestMockQuerier) WithTx(_ pgx.Tx) querier {
	return m
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
