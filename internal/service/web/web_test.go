package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache()

	// Cache miss
	if _, ok := cache.Get("non_existent"); ok {
		t.Fatal("expected cache miss for non_existent key")
	}

	// Cache set and hit
	cache.Set("key1", "val1", 50*time.Millisecond)
	val, ok := cache.Get("key1")
	if !ok || val != "val1" {
		t.Fatalf("expected val1, got %v", val)
	}

	// Cache expiration
	time.Sleep(60 * time.Millisecond)
	if _, ok := cache.Get("key1"); ok {
		t.Fatal("expected key1 to be expired")
	}
}

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewRateLimiter(3, 100*time.Millisecond)
	defer limiter.Stop()

	ip := "192.168.1.1"

	// First 3 requests should be allowed
	for i := 1; i <= 3; i++ {
		if !limiter.Allow(ip) {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	// 4th request should be blocked
	if limiter.Allow(ip) {
		t.Fatal("4th request should be rate-limited")
	}

	// After window expiration, should be allowed again
	time.Sleep(120 * time.Millisecond)
	if !limiter.Allow(ip) {
		t.Fatal("request after window expiry should be allowed")
	}
}

func TestGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GzipMiddleware())

	r.GET("/test-gzip", func(c *gin.Context) {
		c.String(http.StatusOK, "hello compressed world!")
	})

	// Without Accept-Encoding
	reqNoGzip := httptest.NewRequest("GET", "/test-gzip", nil)
	wNoGzip := httptest.NewRecorder()
	r.ServeHTTP(wNoGzip, reqNoGzip)

	if wNoGzip.Header().Get("Content-Encoding") == "gzip" {
		t.Fatal("expected uncompressed response when Accept-Encoding is omitted")
	}

	// With Accept-Encoding: gzip
	reqGzip := httptest.NewRequest("GET", "/test-gzip", nil)
	reqGzip.Header.Set("Accept-Encoding", "gzip")
	wGzip := httptest.NewRecorder()
	r.ServeHTTP(wGzip, reqGzip)

	if wGzip.Header().Get("Content-Encoding") != "gzip" {
		t.Fatal("expected gzip Content-Encoding header")
	}
}
