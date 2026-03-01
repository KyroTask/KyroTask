package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiters per IP address.
type RateLimiter struct {
	limiters sync.Map
	rps      rate.Limit // Requests per second
	burst    int        // Max burst size
}

// NewRateLimiter creates a new thread-safe rate limiter.
func NewRateLimiter(rps rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		rps:   rps,
		burst: burst,
	}
}

// getLimiter retrieves or creates a limiter for the given IP address.
func (r *RateLimiter) getLimiter(ip string) *rate.Limiter {
	limiter, exists := r.limiters.Load(ip)
	if !exists {
		newLimiter := rate.NewLimiter(r.rps, r.burst)
		r.limiters.Store(ip, newLimiter)
		return newLimiter
	}
	return limiter.(*rate.Limiter)
}

// Limit is the Gin middleware function that enforces rate limiting.
func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := r.getLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
