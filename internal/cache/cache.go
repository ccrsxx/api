package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type ProviderCache string

const (
	ProviderMemory     ProviderCache = "memory"
	ProviderCloudflare ProviderCache = "cloudflare-kv"
)

type cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

var cacheManager = &CacheManager{
	memory: newMemoryCache(defaultCleanupInterval),
}

type CacheManager struct {
	memory cache
	// TODO: Add Cloudflare KV later
}

func StaticTTL[T any](d time.Duration) func(T) time.Duration {
	return func(_ T) time.Duration {
		return d
	}
}

func GetCachedData[T any](
	ctx context.Context,
	key string,
	provider ProviderCache,
	fetcher func() (T, error),
	ttlFunc func(T) time.Duration,
) (T, error) {
	var empty T

	var client cache

	switch provider {
	case ProviderCloudflare:
		// TODO: Add Cloudflare KV later
		client = nil
	case ProviderMemory:
		client = cacheManager.memory
	}

	// Fallback to memory cache if no client is specified
	if client == nil {
		client = cacheManager.memory
	}

	slog.Debug("get cached data", "key", key, "provider", provider)

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
