package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	RPS   int
	Burst int
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
	return &RateLimiter{RPS: rps, Burst: burst}
}

func (rl *RateLimiter) RateLimiting(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rl.RPS), rl.Burst) // 100 запросов в секунду с burst 100
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
