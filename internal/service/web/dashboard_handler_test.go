package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDashboardHandler_ServeDashboard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewDashboardHandler()
	r.GET("/dashboard", handler.ServeDashboard)

	req, _ := http.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	
	body := w.Body.String()
	assert.Contains(t, body, "LatinaServer")
	assert.Contains(t, body, "CPU Usage")
	assert.Contains(t, body, "RAM Memory")
	assert.Contains(t, body, "Disk Storage")
	assert.Contains(t, body, "Network I/O")
	assert.Contains(t, body, "2 Mbps")
	assert.Contains(t, body, "/api/v1/status")
	assert.Contains(t, body, "/api/v1/trial")
}

func TestServer_DashboardAndRootRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := web.NewServer()
	r := server.SetupRouter()

	// Test /dashboard directly
	reqDash, _ := http.NewRequest(http.MethodGet, "/dashboard", nil)
	wDash := httptest.NewRecorder()
	r.ServeHTTP(wDash, reqDash)

	assert.Equal(t, http.StatusOK, wDash.Code)
	assert.Contains(t, wDash.Header().Get("Content-Type"), "text/html")

	// Test / root with token redirects to /portal?token=xxx
	reqToken, _ := http.NewRequest(http.MethodGet, "/?token=LATINA01", nil)
	wToken := httptest.NewRecorder()
	r.ServeHTTP(wToken, reqToken)

	assert.Equal(t, http.StatusTemporaryRedirect, wToken.Code)
	assert.Equal(t, "/portal?token=LATINA01", wToken.Header().Get("Location"))

	// Test / root serves camouflage publication index with 200 OK (when static assets exist)
	reqRoot, _ := http.NewRequest(http.MethodGet, "/", nil)
	wRoot := httptest.NewRecorder()
	r.ServeHTTP(wRoot, reqRoot)

	if wRoot.Code == http.StatusNotFound {
		t.Skip("Skipping camouflage index assertion: web/dist not generated (run hugo or make build-web)")
	}

	assert.Equal(t, http.StatusOK, wRoot.Code)
	assert.Contains(t, wRoot.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, wRoot.Body.String(), "Lalatina Systems")
}
