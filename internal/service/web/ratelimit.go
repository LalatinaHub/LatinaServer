package web

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements token bucket rate limiting per IP.
type RateLimiter struct {
	clients  map[string]*clientBucket
	mu       sync.RWMutex
	rate     int           // requests per minute
	window   time.Duration // time window
	stopOnce sync.Once
	stopChan chan struct{}
}

type clientBucket struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter with default background context.
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return NewRateLimiterWithContext(context.Background(), rate, window)
}

// NewRateLimiterWithContext creates a new rate limiter and gracefully stops cleanup on ctx.Done() or Stop().
func NewRateLimiterWithContext(ctx context.Context, rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*clientBucket),
		rate:     rate,
		window:   window,
		stopChan: make(chan struct{}),
	}

	// Background cleanup for expired clients every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-rl.stopChan:
				return
			case <-ticker.C:
				rl.cleanup()
			}
		}
	}()

	return rl
}

// Stop terminates the background cleanup goroutine.
func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopChan)
	})
}

// Allow checks if request from IP is allowed under rate limit.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[ip]

	if !exists {
		rl.clients[ip] = &clientBucket{
			tokens:    rl.rate - 1,
			lastReset: now,
		}
		return true
	}

	// Reset bucket if window expired
	if now.Sub(bucket.lastReset) > rl.window {
		bucket.tokens = rl.rate - 1
		bucket.lastReset = now
		return true
	}

	// Check if tokens available
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup removes old client buckets to prevent memory leak.
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, bucket := range rl.clients {
		if now.Sub(bucket.lastReset) > rl.window*2 {
			delete(rl.clients, ip)
		}
	}
}

// Middleware returns Gin middleware for rate limiting.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
