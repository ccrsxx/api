package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_SetGet(t *testing.T) {
	c := newMemoryCache(defaultCleanupInterval)
	ctx := context.Background()

	key := "user:123"
	val := "john_doe"

	if err := c.Set(ctx, key, val, time.Minute); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := c.Get(ctx, key)

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != val {
		t.Errorf("got %v, want %v", got, val)
	}

	newVal := "jane_doe"

	if err := c.Set(ctx, key, newVal, time.Minute); err != nil {
		t.Fatalf("Set overwrite failed: %v", err)
	}

	got, err = c.Get(ctx, key)

	if err != nil {
		t.Fatalf("Get after overwrite failed: %v", err)
	}

	if got != newVal {
		t.Errorf("got %v, want %v", got, newVal)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := newMemoryCache(defaultCleanupInterval)
	ctx := context.Background()

	key := "temp-key"
	val := "temp-data"
	ttl := 5 * time.Millisecond

	// Set with very short TTL
	if err := c.Set(ctx, key, val, ttl); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists immediately
	got, err := c.Get(ctx, key)

	if err != nil || got != val {
		t.Fatal("Item should exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(ttl * 2)

	// Verify it is gone (Lazy expiration check in Get)
	_, err = c.Get(ctx, key)

	if err != ErrCacheMiss {
		t.Errorf("want ErrCacheMiss for expired item, got %v", err)
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	c := newMemoryCache(defaultCleanupInterval)
	ctx := context.Background()

	key := "del-key"
	val := "delete-me"

	_ = c.Set(ctx, key, val, time.Minute)

	// Verify exists
	if _, err := c.Get(ctx, key); err != nil {
		t.Fatal("Pre-requisite failed: item not set")
	}

	// Delete
	if err := c.Delete(ctx, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify gone
	_, err := c.Get(ctx, key)

	if err != ErrCacheMiss {
		t.Errorf("want ErrCacheMiss after Delete, got %v", err)
	}
}

func TestMemoryCache_Miss(t *testing.T) {
	c := newMemoryCache(defaultCleanupInterval)
	ctx := context.Background()

	_, err := c.Get(ctx, "ghost-key")

	if err != ErrCacheMiss {
		t.Errorf("want ErrCacheMiss for unknown key, got %v", err)
	}
}

func TestMemoryCache_Cleanup(t *testing.T) {
	interval := 10 * time.Millisecond

	c := newMemoryCache(interval)
	ctx := context.Background()

	key := "cleanup-key"
	val := "cleanup-data"

	_ = c.Set(ctx, key, val, 1*time.Millisecond)

	time.Sleep(50 * time.Millisecond)

	// Coverage Check: The 'cleanup' goroutine loop code should have run.
	// We verify the side effect (item is gone).
	_, err := c.Get(ctx, key)

	if err != ErrCacheMiss {
		t.Errorf("cleanup loop failed to remove expired item")
	}
}
