package cache

import (
	"context"
	"sync"
	"time"
)

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]item
}

type item struct {
	value     any
	expiresAt time.Time
}

func NewMemoryCache() *MemoryCache {
	store := &MemoryCache{
		items: make(map[string]item),
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
	const cleanupInterval = 5 * time.Minute

	ticker := time.NewTicker(cleanupInterval)

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
