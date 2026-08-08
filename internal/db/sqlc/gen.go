package sqlc

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewQueries(ctx context.Context, databaseString string) (*pgxpool.Pool, *Queries) {
	poolConfig, err := pgxpool.ParseConfig(databaseString)

	if err != nil {
		panic(fmt.Errorf("db config parse error: %w", err))
	}

	// Emit a span per query so slow SQL shows up in the trace waterfall
	// alongside the HTTP and upstream spans. Tracing being disabled makes this
	// a no-op via the global no-op tracer provider.
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithTrimSQLInSpanName(),
	)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		panic(fmt.Errorf("db create error: %w", err))
	}

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("db ping error: %w", err))
	}

	queries := New(pool)

	return pool, queries
}
