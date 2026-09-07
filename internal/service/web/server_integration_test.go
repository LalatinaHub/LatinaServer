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

	t.Run("root with token query redirects to /portal?token=xxx", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/?token=MEMBER88", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/portal?token=MEMBER88", w.Header().Get("Location"))
	})

	t.Run("static camouflage publication site", func(t *testing.T) {
		reqCheck, _ := http.NewRequest(http.MethodGet, "/", nil)
		wCheck := httptest.NewRecorder()
		r.ServeHTTP(wCheck, reqCheck)
		if wCheck.Code == http.StatusNotFound {
			t.Skip("Skipping static file integration checks: web/dist not found (run hugo or make build-web)")
		}

		assert.Equal(t, http.StatusOK, wCheck.Code)
		assert.Contains(t, wCheck.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, wCheck.Body.String(), "Lalatina Systems")

		reqPortal, _ := http.NewRequest(http.MethodGet, "/portal", nil)
		wPortal := httptest.NewRecorder()
		r.ServeHTTP(wPortal, reqPortal)
		assert.Equal(t, http.StatusOK, wPortal.Code)
		assert.Contains(t, wPortal.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, wPortal.Body.String(), "Portal Layanan Mandiri")

		reqSub, _ := http.NewRequest(http.MethodGet, "/portal/session", nil)
		wSub := httptest.NewRecorder()
		r.ServeHTTP(wSub, reqSub)
		assert.Equal(t, http.StatusOK, wSub.Code)
		assert.Contains(t, wSub.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, wSub.Body.String(), "Portal Layanan Mandiri")

		reqPosts, _ := http.NewRequest(http.MethodGet, "/posts/", nil)
		wPosts := httptest.NewRecorder()
		r.ServeHTTP(wPosts, reqPosts)
		assert.Equal(t, http.StatusOK, wPosts.Code)
		assert.Contains(t, wPosts.Body.String(), "Katalog publikasi ilmiah")

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
