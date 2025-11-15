package middleware

import (
	"net/http"

	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/config"
	"golang.org/x/time/rate"
)

// RateLimit limits the number of requests per second
func RateLimit(cfg config.RateLimitConfig) Middleware {
	limiter := rate.NewLimiter(
		rate.Limit(cfg.RequestsPerSecond),
		cfg.Burst,
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
