package middleware

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/utils"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	limit    rate.Limit
	burst    int
	policy   string
	window   time.Duration
	visitors map[string]*visitor
}

func newLimiterFromConfig(requests int, window time.Duration) *RateLimiter {
	limit := rate.Limit(float64(requests) / window.Seconds())
	burst := requests
	policy := fmt.Sprintf("%d; w=%d", requests, int(window.Seconds()))

	rl := &RateLimiter{
		limit:    limit,
		burst:    burst,
		policy:   policy,
		window:   window,
		visitors: map[string]*visitor{},
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
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

func (rl *RateLimiter) handleRateLimit(w http.ResponseWriter, r *http.Request) error {
	ip := utils.GetIpAddressFromRequest(r)

	limiter := rl.getVisitor(ip)

	allowed := limiter.Allow()

	currentTokens := limiter.Tokens()

	remaining := int(math.Max(0, currentTokens))

	resetTime := 0.0

	if rl.limit > 0 {
		resetTime = (float64(rl.burst) - currentTokens) / float64(rl.limit)
	}

	w.Header().Set("RateLimit-Limit", strconv.Itoa(rl.burst))
	w.Header().Set("RateLimit-Reset", strconv.Itoa(int(math.Ceil(resetTime))))
	w.Header().Set("RateLimit-Policy", rl.policy)
	w.Header().Set("RateLimit-Remaining", strconv.Itoa(remaining))

	if !allowed {
		retryAfter := 1.0 / float64(rl.limit)
		waitSecs := int(math.Max(1, math.Ceil(retryAfter)))

		w.Header().Set("Retry-After", strconv.Itoa(waitSecs))

		return api.NewHttpError(http.StatusTooManyRequests, "Too many requests, please try again later.", nil)
	}

	return nil
}

func (rl *RateLimiter) cleanup() {
	interval := max(rl.window*4, time.Minute)

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()

		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > interval {
				delete(rl.visitors, ip)
			}
		}

		rl.mu.Unlock()
	}
}

func GlobalRateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	rl := newLimiterFromConfig(requests, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := rl.handleRateLimit(w, r); err != nil {
				api.HandleHttpError(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func HandlerRateLimit(requests int, window time.Duration) func(api.HTTPHandlerWithErr) api.HTTPHandlerWithErr {
	rl := newLimiterFromConfig(requests, window)

	return func(next api.HTTPHandlerWithErr) api.HTTPHandlerWithErr {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := rl.handleRateLimit(w, r); err != nil {
				return err
			}

			return next(w, r)
		}
	}
}
