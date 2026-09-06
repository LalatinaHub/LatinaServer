package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServer_FullIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := web.NewServer()
	r := server.SetupRouter()

	t.Run("root serves camouflage publication with 200 OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "Lalatina Systems")
	})

	t.Run("root with token query redirects to /portal?token=xxx", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/?token=MEMBER88", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/portal?token=MEMBER88", w.Header().Get("Location"))
	})

	t.Run("portal route serves portal SPA with 200 OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/portal", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "Portal Layanan Mandiri")
	})

	t.Run("portal subpath serves portal SPA with 200 OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/portal/session", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "Portal Layanan Mandiri")
	})

	t.Run("static camouflage research posts served", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/posts/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Katalog publikasi ilmiah")

		// Also check specific post
		reqPost, _ := http.NewRequest(http.MethodGet, "/posts/resilient-packet-routing-under-adversarial-conditions/", nil)
		wPost := httptest.NewRecorder()
		r.ServeHTTP(wPost, reqPost)

		assert.Equal(t, http.StatusOK, wPost.Code)
		assert.Contains(t, wPost.Body.String(), "Perutean Paket")
	})

	t.Run("api v1 ping returns pong", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Pong", w.Body.String())
	})

	t.Run("api v1 trial returns 200 OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/trial", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "vmess://")
	})

	t.Run("api v1 portal metadata endpoints", func(t *testing.T) {
		endpoints := []string{
			"/api/v1/portal/servers",
			"/api/v1/portal/relays",
			"/api/v1/portal/status",
			"/api/v1/portal/wildcards",
			"/api/v1/portal/info",
		}

		for _, ep := range endpoints {
			req, _ := http.NewRequest(http.MethodGet, ep, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Endpoint %s failed", ep)
		}
	})

	t.Run("api v1 portal protected endpoints require token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/portal/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
