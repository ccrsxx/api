package cache

import (
	"context"
	"sync"
	"time"
)

const DefaultCleanupInterval = 5 * time.Minute

type MemoryCache struct {
	mu              sync.RWMutex
	items           map[string]item
	cleanupInterval time.Duration
}

type item struct {
	value     any
	expiresAt time.Time
}

func NewMemoryCache(ctx context.Context, cleanupInterval time.Duration) *MemoryCache {
	if cleanupInterval <= 0 {
		cleanupInterval = DefaultCleanupInterval
	}

	store := &MemoryCache{
		items:           map[string]item{},
		cleanupInterval: cleanupInterval,
	}

	go store.cleanup(ctx, cleanupInterval)

	return store
}

func (mc *MemoryCache) Get(ctx context.Context, key string) (any, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	val, found := mc.items[key]

	if !found || time.Now().After(val.expiresAt) {
		return nil, ErrCacheMiss
	}

	return val.value, nil
}

func (mc *MemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.items, key)

	return nil
}

func (mc *MemoryCache) cleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mc.mu.Lock()

			now := time.Now()

			for key, item := range mc.items {
				if now.After(item.expiresAt) {
					delete(mc.items, key)
				}
			}

			mc.mu.Unlock()
		}
	}
}
