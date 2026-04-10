package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

// Cache defines the standard behavior for all cache providers.
type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

func StaticTTL[T any](d time.Duration) func(T) time.Duration {
	return func(_ T) time.Duration {
		return d
	}
}

// GetOrFetch tries to get data from the injected Cache client.
// If it misses, it runs the fetcher and sets the data in the background.
func GetOrFetch[T any](
	ctx context.Context,
	client Cache,
	key string,
	fetcher func() (T, error),
	ttlFunc func(T) time.Duration,
) (T, error) {
	var empty T

	if client == nil {
		return fetcher()
	}

	slog.Debug("get cached data", "key", key)

	val, err := client.Get(ctx, key)

	if err == nil {
		if casted, ok := val.(T); ok {
			slog.Debug("cache hit", "key", key)
			return casted, nil
		}
	}

	data, err := fetcher()

	if err != nil {
		return empty, err
	}

	slog.Debug("cache miss", "key", key)

	ttl := ttlFunc(data)

	go func() {
		err := client.Set(context.Background(), key, data, ttl)

		slog.Debug("set cache", "key", key)

		if err != nil {
			slog.Warn("set cache error", "key", key, "error", err)
		}
	}()

	return data, nil
}
