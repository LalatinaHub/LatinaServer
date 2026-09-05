package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Ping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewHealthHandler()
	r.GET("/ping", handler.Ping)

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Pong", w.Body.String())
}

func TestProxyHandler_Check_NoIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewProxyHandler()
	r.GET("/check", handler.Check)

	req, _ := http.NewRequest(http.MethodGet, "/check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "No proxy provided")
}

func TestProxyHandler_Relays(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewProxyHandler()
	r.GET("/relays", handler.Relays)

	req, _ := http.NewRequest(http.MethodGet, "/relays", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServer_AdminReloadRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := web.NewServer()
	r := server.SetupRouter()

	// Request without authorization header should be 401
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/reload", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

