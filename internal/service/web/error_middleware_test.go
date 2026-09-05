package web_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/gin-gonic/gin"
)

func TestErrorMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(web.RequestLoggerMiddleware())
	r.Use(web.ErrorMiddleware())

	r.GET("/validation-error", func(c *gin.Context) {
		c.Error(appErrors.NewValidationError("invalid param format", nil))
	})

	r.GET("/database-error", func(c *gin.Context) {
		c.Error(appErrors.NewDatabaseError("database connection timeout", errors.New("timeout")))
	})

	r.GET("/network-error", func(c *gin.Context) {
		c.Error(appErrors.NewNetworkError("cloudflare unreachable", errors.New("timeout")))
	})

	r.GET("/generic-error", func(c *gin.Context) {
		c.Error(errors.New("something crashed"))
	})

	// Test validation error
	req, _ := http.NewRequest(http.MethodGet, "/validation-error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp web.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Success || resp.Type != string(appErrors.ErrorTypeValidation) {
		t.Fatalf("unexpected response: %+v", resp)
	}

	// Test request ID header propagation
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}

	// Test database error
	reqDB, _ := http.NewRequest(http.MethodGet, "/database-error", nil)
	wDB := httptest.NewRecorder()
	r.ServeHTTP(wDB, reqDB)

	if wDB.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", wDB.Code)
	}

	// Test network error
	reqNet, _ := http.NewRequest(http.MethodGet, "/network-error", nil)
	wNet := httptest.NewRecorder()
	r.ServeHTTP(wNet, reqNet)

	if wNet.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", wNet.Code)
	}
}
