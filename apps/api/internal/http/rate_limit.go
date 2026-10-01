package httpserver

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	entries     map[string]rateLimitEntry
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		entries:     make(map[string]rateLimitEntry),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			next.ServeHTTP(w, r)
			return
		}

		key := clientAddress(r)
		now := time.Now()

		rl.mu.Lock()

		rl.cleanup(now)

		entry, exists := rl.entries[key]

		if !exists || now.Sub(entry.windowStart) >= rl.window {
			entry = rateLimitEntry{
				windowStart: now,
			}
		}

		entry.count++
		rl.entries[key] = entry

		remaining := rl.maxRequests - entry.count
		if remaining < 0 {
			remaining = 0
		}

		allowed := entry.count <= rl.maxRequests

		rl.mu.Unlock()

		w.Header().Set("X-RateLimit-Limit", formatInt(rl.maxRequests))
		w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))

		if !allowed {
			retryAfter := int(rl.window - now.Sub(entry.windowStart))
			if retryAfter < 1 {
				retryAfter = 1
			}

			w.Header().Set("Retry-After", formatInt(retryAfter))
			writeError(
				w,
				http.StatusTooManyRequests,
				"rate limit exceeded",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup(now time.Time) {
	for key, entry := range rl.entries {
		if now.Sub(entry.windowStart) >= rl.window {
			delete(rl.entries, key)
		}
	}
}

func clientAddress(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

func formatInt(value int) string {
	return strconv.Itoa(value)
}
