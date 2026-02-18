package middleware

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/utils"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	limit    rate.Limit
	burst    int
	policy   string
	window   time.Duration
	visitors map[string]*visitor
}

func newLimiterFromConfig(requests int, window time.Duration) *rateLimiter {
	// Convert requests/window to requests/second (Token Bucket Rate)
	// Example: 100 reqs / 60s = 1.666 reqs/sec
	limit := rate.Limit(float64(requests) / window.Seconds())
	burst := requests

	// RFC RateLimit-Policy format: "limit; w=window_in_seconds"
	// Example: "100; w=60"
	policy := fmt.Sprintf("%d; w=%d", requests, int(window.Seconds()))

	rl := &rateLimiter{
		limit:    limit,
		burst:    burst,
		policy:   policy,
		window:   window,
		visitors: map[string]*visitor{},
	}

	// Calculate cleanup interval based on window (min 1 minute)
	interval := max(window*4, time.Minute)

	// Start the background cleanup routine.
	// We pass 'nil' for the stop channel so it runs forever in production.
	go rl.cleanup(interval, nil)

	return rl
}

func (rl *rateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]

	if !exists {
		v = &visitor{
			limiter:  rate.NewLimiter(rl.limit, rl.burst),
			lastSeen: time.Now(),
		}

		rl.visitors[ip] = v

		return v.limiter

	}

	v.lastSeen = time.Now()

	return v.limiter
}

func (rl *rateLimiter) handleRateLimit(w http.ResponseWriter, r *http.Request) error {
	ip := utils.GetIpAddressFromRequest(r)

	limiter := rl.getVisitor(ip)

	allowed := limiter.Allow()

	currentTokens := limiter.Tokens()

	// Ensure we don't display negative tokens due to float precision
	remaining := int(math.Max(0, currentTokens))

	resetTime := 0.0

	if rl.limit > 0 {
		// Calculate time until the bucket is completely full again
		// Formula: (Capacity - CurrentTokens) / RefillRate
		resetTime = (float64(rl.burst) - currentTokens) / float64(rl.limit)
	}

	w.Header().Set("RateLimit-Limit", strconv.Itoa(rl.burst))
	w.Header().Set("RateLimit-Reset", strconv.Itoa(int(math.Ceil(resetTime))))
	w.Header().Set("RateLimit-Policy", rl.policy)
	w.Header().Set("RateLimit-Remaining", strconv.Itoa(remaining))

	if !allowed {
		// Calculate exact wait time to regenerate ONE new token
		// Formula: 1.0 / RefillRate (tokens per second)
		retryAfter := 1.0 / float64(rl.limit)

		// Round up to ensure the server has definitely refilled by the time they retry
		waitSecs := int(math.Max(1, math.Ceil(retryAfter)))

		w.Header().Set("Retry-After", strconv.Itoa(waitSecs))

		return &api.HttpError{
			StatusCode: http.StatusTooManyRequests,
			Message:    "Too many requests, please try again later.",
		}
	}

	return nil
}

func (rl *rateLimiter) pruneVisitors(interval time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > interval {
			delete(rl.visitors, ip)
		}
	}
}

func (rl *rateLimiter) cleanup(interval time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			rl.pruneVisitors(interval)
		}
	}
}

func RateLimit(requests int, window time.Duration) func(next http.Handler) http.Handler {
	rl := newLimiterFromConfig(requests, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := rl.handleRateLimit(w, r); err != nil {
				// No need to wrap error, as it's already returning http error from handleRateLimit
				api.HandleHttpError(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
