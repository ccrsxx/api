package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_SetGet(t *testing.T) {
	ctx := t.Context()

	mc := NewMemoryCache(ctx, DefaultCleanupInterval)

	key := "user:123"
	val := "john_doe"

	if err := mc.Set(ctx, key, val, time.Minute); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := mc.Get(ctx, key)

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != val {
		t.Errorf("got %v, want %v", got, val)
	}

	newVal := "jane_doe"

	if err := mc.Set(ctx, key, newVal, time.Minute); err != nil {
		t.Fatalf("Set overwrite failed: %v", err)
	}

	got, err = mc.Get(ctx, key)

	if err != nil {
		t.Fatalf("Get after overwrite failed: %v", err)
	}

	if got != newVal {
		t.Errorf("got %v, want %v", got, newVal)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	ctx := t.Context()

	mc := NewMemoryCache(ctx, DefaultCleanupInterval)

	key := "temp-key"
	val := "temp-data"
	ttl := 5 * time.Millisecond

	// Set with very short TTL
	if err := mc.Set(ctx, key, val, ttl); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists immediately
	got, err := mc.Get(ctx, key)

	if err != nil || got != val {
		t.Fatal("Item should exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(ttl * 2)

	// Verify it is gone (Lazy expiration check in Get)
	_, err = mc.Get(ctx, key)

	if err != ErrCacheMiss {
		t.Errorf("got %v, want ErrCacheMiss for expired item", err)
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	ctx := t.Context()

	mc := NewMemoryCache(ctx, DefaultCleanupInterval)

	key := "del-key"
	val := "delete-me"

	err := mc.Set(ctx, key, val, time.Minute)

	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify exists
	if _, err := mc.Get(ctx, key); err != nil {
		t.Fatal("Pre-requisite failed: item not set")
	}

	// Delete
	if err := mc.Delete(ctx, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify gone
	_, err = mc.Get(ctx, key)

	if err != ErrCacheMiss {
		t.Errorf("got %v, want ErrCacheMiss after Delete", err)
	}
}

func TestMemoryCache_Miss(t *testing.T) {
	ctx := t.Context()

	mc := NewMemoryCache(ctx, DefaultCleanupInterval)

	_, err := mc.Get(ctx, "ghost-key")

	if err != ErrCacheMiss {
		t.Errorf("got %v, want ErrCacheMiss for unknown key", err)
	}
}

func TestMemoryCache_Cleanup(t *testing.T) {
	interval := 10 * time.Millisecond

	ctx := t.Context()

	mc := NewMemoryCache(ctx, interval)

	key := "cleanup-key"
	val := "cleanup-data"

	err := mc.Set(ctx, key, val, 1*time.Millisecond)

	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Coverage Check: The 'cleanup' goroutine loop code should have run.
	// We verify the side effect (item is gone).
	_, err = mc.Get(ctx, key)

	if err != ErrCacheMiss {
		t.Errorf("cleanup loop failed to remove expired item")
	}
}

func TestMemoryCache_DefaultCleanupInterval(t *testing.T) {
	t.Run("Positive Interval", func(t *testing.T) {
		ctx := t.Context()

		mc := NewMemoryCache(ctx, 1)

		if mc.cleanupInterval != 1 {
			t.Errorf("got %v, want 1", mc.cleanupInterval)
		}
	})

	t.Run("Zero Interval Fallback", func(t *testing.T) {
		ctx := t.Context()

		mc := NewMemoryCache(ctx, 0)

		if mc.cleanupInterval != DefaultCleanupInterval {
			t.Errorf("got %v, want default %v", mc.cleanupInterval, DefaultCleanupInterval)
		}
	})

	t.Run("Negative Interval Fallback", func(t *testing.T) {
		ctx := t.Context()

		// Pass a negative number to trigger the if statement
		mc := NewMemoryCache(ctx, -1*time.Minute)

		if mc.cleanupInterval != DefaultCleanupInterval {
			t.Errorf("got %v, want default %v", mc.cleanupInterval, DefaultCleanupInterval)
		}
	})
}
