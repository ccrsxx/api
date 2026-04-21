package auth

import (
	"context"
	"errors"

	"github.com/ccrsxx/api/internal/db/sqlc"
)

type contextKey string

const userContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) (sqlc.GetUserWithAccountByIDRow, error) {
	user, ok := ctx.Value(userContextKey).(sqlc.GetUserWithAccountByIDRow)

	if !ok {
		return sqlc.GetUserWithAccountByIDRow{}, errors.New("get current user from context error")
	}

	return user, nil
}
