package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGetOrFetch(t *testing.T) {
	ctx := context.Background()

	t.Run("Cache Miss -> Fetch -> Set", func(t *testing.T) {
		ctx := t.Context()

		c := NewMemoryCache(ctx, DefaultCleanupInterval)

		key := "miss-key"
		want := "fetched-want"

		fetcher := func() (string, error) {
			return want, nil
		}

		got, err := GetOrFetch(ctx, c, key, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		time.Sleep(10 * time.Millisecond) // Wait for async set

		cached, err := c.Get(ctx, key)
		if err != nil {
			t.Fatalf("error getting from cache: %v", err)
		}

		if cached != want {
			t.Error("value was not saved to cache")
		}
	})

	t.Run("Cache Hit -> Return Immediate", func(t *testing.T) {
		ctx := t.Context()

		c := NewMemoryCache(ctx, DefaultCleanupInterval)

		key := "hit-key"
		want := "cached-want"

		err := c.Set(ctx, key, want, time.Minute)

		if err != nil {
			t.Fatalf("error setting up cache: %v", err)
		}

		fetcher := func() (string, error) {
			t.Fatal("fetcher should not be called")
			return "", nil
		}

		got, err := GetOrFetch(ctx, c, key, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Fetcher Error -> Return Error", func(t *testing.T) {
		ctx := t.Context()

		c := NewMemoryCache(ctx, DefaultCleanupInterval)

		wantErr := errors.New("db dead")

		fetcher := func() (string, error) {
			return "", wantErr
		}

		_, err := GetOrFetch(ctx, c, "err-key", fetcher, StaticTTL[string](time.Minute))

		if !errors.Is(err, wantErr) {
			t.Errorf("got error %v, want %v", err, wantErr)
		}
	})

	t.Run("Nil Cache Fallback", func(t *testing.T) {
		want := "data"
		fetcher := func() (string, error) { return want, nil }

		got, err := GetOrFetch(ctx, nil, "fallback", fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
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

func (f *faultyCache) Get(ctx context.Context, key string) (any, error) { return nil, ErrCacheMiss }
func (f *faultyCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return errors.New("forced storage error")
}
func (f *faultyCache) Delete(ctx context.Context, key string) error { return nil }

func TestGetOrFetchCoverage(t *testing.T) {
	ctx := context.Background()

	t.Run("Async Set Error (Log Warning)", func(t *testing.T) {
		c := &faultyCache{}
		want := "log-want"
		fetcher := func() (string, error) { return want, nil }

		got, err := GetOrFetch(ctx, c, "log-key", fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		// Wait for goroutine to hit the error block
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Type Assertion Failure (Wrong Type in Cache)", func(t *testing.T) {
		ctx := t.Context()

		c := NewMemoryCache(ctx, DefaultCleanupInterval)

		key := "wrong-type-key"
		err := c.Set(ctx, key, 999, time.Minute)

		if err != nil {
			t.Fatalf("error setting up cache: %v", err)
		}

		want := "correct-string"
		fetcher := func() (string, error) { return want, nil }

		got, err := GetOrFetch(ctx, c, key, fetcher, StaticTTL[string](time.Minute))

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
