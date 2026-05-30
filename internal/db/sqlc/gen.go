package sqlc

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewQueries(ctx context.Context, databaseString string) (*pgxpool.Pool, *Queries) {
	pool, err := pgxpool.New(ctx, databaseString)

	if err != nil {
		panic(fmt.Errorf("db create error: %w", err))
	}

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("db ping error : %w", err))
	}

	queries := New(pool)

	return pool, queries
}
