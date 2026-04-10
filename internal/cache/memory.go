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

func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	if cleanupInterval <= 0 {
		cleanupInterval = DefaultCleanupInterval
	}

	store := &MemoryCache{
		items:           map[string]item{},
		cleanupInterval: cleanupInterval,
	}

	go store.cleanup()

	return store
}

func (m *MemoryCache) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, found := m.items[key]

	if !found || time.Now().After(val.expiresAt) {
		return nil, ErrCacheMiss
	}

	return val.value, nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)

	return nil
}

func (m *MemoryCache) cleanup() {
	ticker := time.NewTicker(m.cleanupInterval)

	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		now := time.Now()

		for key, item := range m.items {
			if now.After(item.expiresAt) {
				delete(m.items, key)
			}
		}

		m.mu.Unlock()
	}
}
