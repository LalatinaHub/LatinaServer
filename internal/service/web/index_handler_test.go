package web_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIndexHandler_ServeIndex(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewIndexHandler()
	r.GET("/", handler.ServeIndex)

	// Case 2: GET / with token query param redirects to /portal?token=xxx
	reqToken, _ := http.NewRequest(http.MethodGet, "/?token=SAMPLE_TOKEN", nil)
	wToken := httptest.NewRecorder()
	r.ServeHTTP(wToken, reqToken)

	assert.Equal(t, http.StatusTemporaryRedirect, wToken.Code)
	assert.Equal(t, "/portal?token=SAMPLE_TOKEN", wToken.Header().Get("Location"))

	// Case 1: Plain GET / without token returns 200 OK and serves camouflage HTML (when static assets exist)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Skip("Skipping camouflage index assertion: web/dist static assets not generated (run hugo or make build-web)")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, w.Body.String(), "Lalatina Systems")
}

func TestIndexHandler_CustomPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewIndexHandlerWithPath("non_existent_file.html")
	r.GET("/", handler.ServeIndex)

	// Case 1: Non-existent custom path returns 404
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Case 2: Existing custom path returns 200 OK
	tmpFile := t.TempDir() + "/index.html"
	_ = os.WriteFile(tmpFile, []byte("<html>Custom Index</html>"), 0644)
	handler2 := web.NewIndexHandlerWithPath(tmpFile)
	r2 := gin.New()
	r2.GET("/", handler2.ServeIndex)

	req2, _ := http.NewRequest(http.MethodGet, "/", nil)
	w2 := httptest.NewRecorder()
	r2.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Custom Index")
}
