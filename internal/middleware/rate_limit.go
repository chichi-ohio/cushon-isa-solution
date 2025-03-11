package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	sync.RWMutex
	requests map[string][]time.Time
	window   time.Duration
	limit    int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		limit:    limit,
	}
}

func (rl *RateLimiter) cleanup(now time.Time) {
	rl.Lock()
	defer rl.Unlock()

	for ip, times := range rl.requests {
		var valid []time.Time
		for _, t := range times {
			if now.Sub(t) <= rl.window {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = valid
		}
	}
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	now := time.Now()
	rl.cleanup(now)

	rl.Lock()
	defer rl.Unlock()

	times := rl.requests[ip]
	if len(times) < rl.limit {
		rl.requests[ip] = append(times, now)
		return true
	}

	return false
}

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(time.Minute, requestsPerMinute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.isAllowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
