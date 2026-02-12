package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGetCachedData(t *testing.T) {
	realCache := cacheManager.memory

	defer func() { cacheManager.memory = realCache }()

	ctx := context.Background()

	t.Run("Cache Miss -> Fetch -> Set", func(t *testing.T) {
		cacheManager.memory = newMemoryCache()

		key := "miss-key"
		want := "fetched-want"

		fetcher := func() (string, error) {
			return want, nil
		}

		got, err := GetCachedData(ctx, key, ProviderMemory, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		time.Sleep(10 * time.Millisecond)

		if cached, _ := cacheManager.memory.Get(ctx, key); cached != want {
			t.Error("value was not saved to cache")
		}
	})

	t.Run("Cache Hit -> Return Immediate", func(t *testing.T) {
		cacheManager.memory = newMemoryCache()

		key := "hit-key"
		want := "cached-want"

		_ = cacheManager.memory.Set(ctx, key, want, time.Minute)

		fetcher := func() (string, error) {
			t.Fatal("fetcher should not be called")
			return "wrong", nil
		}

		got, err := GetCachedData(ctx, key, ProviderMemory, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Fetcher Error -> Return Error", func(t *testing.T) {
		cacheManager.memory = newMemoryCache()
		expectedErr := errors.New("db dead")

		fetcher := func() (string, error) {
			return "", expectedErr
		}

		_, err := GetCachedData(ctx, "err-key", ProviderMemory, fetcher, StaticTTL[string](time.Minute))

		if !errors.Is(err, expectedErr) {
			t.Errorf("got error %v, want %v", err, expectedErr)
		}
	})

	t.Run("Provider Fallback (Cloudflare -> Memory)", func(t *testing.T) {
		cacheManager.memory = newMemoryCache()

		key := "fallback-key"
		want := "data"

		fetcher := func() (string, error) {
			return want, nil
		}

		got, err := GetCachedData(ctx, key, ProviderCloudflare, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestStaticTTL(t *testing.T) {
	ttl := 10 * time.Minute

	if got := StaticTTL[string](ttl)("any"); got != ttl {
		t.Errorf("got %v, want %v", got, ttl)
	}
}

type faultyCache struct{}

func (f *faultyCache) Get(ctx context.Context, key string) (any, error) {
	return nil, ErrCacheMiss
}

func (f *faultyCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return errors.New("forced storage error")
}

func (f *faultyCache) Delete(ctx context.Context, key string) error {
	return nil
}

func TestGetCachedDataCoverage(t *testing.T) {
	realCache := cacheManager.memory

	defer func() { cacheManager.memory = realCache }()

	ctx := context.Background()

	t.Run("Async Set Error (Log Warning)", func(t *testing.T) {
		cacheManager.memory = &faultyCache{}

		key := "log-key"
		want := "log-want"

		fetcher := func() (string, error) {
			return want, nil
		}

		got, err := GetCachedData(ctx, key, ProviderMemory, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		// Wait for goroutine to hit the error block
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Type Assertion Failure (Wrong Type in Cache)", func(t *testing.T) {
		cacheManager.memory = newMemoryCache()

		key := "wrong-type-key"

		_ = cacheManager.memory.Set(ctx, key, 999, time.Minute)

		want := "correct-string"

		fetcher := func() (string, error) {
			return want, nil
		}

		got, err := GetCachedData(ctx, key, ProviderMemory, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
