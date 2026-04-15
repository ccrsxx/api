package sqlc

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewQueries(ctx context.Context, databaseString string) (*pgxpool.Pool, *Queries) {
	pool, err := pgxpool.New(ctx, databaseString)

	if err != nil {
		slog.Error("db creation error", "error", err)
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("db ping error", "error", err)
		panic(err)
	}

	queries := New(pool)

	return pool, queries
}
