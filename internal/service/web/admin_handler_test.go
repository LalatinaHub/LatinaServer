package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_AuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		configured   string
		authHeader   string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "no token configured",
			configured:   "",
			authHeader:   "Bearer some-token",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Admin API token not configured",
		},
		{
			name:         "missing authorization header",
			configured:   "secret-token-123",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Missing Authorization header",
		},
		{
			name:         "malformed header - not bearer",
			configured:   "secret-token-123",
			authHeader:   "Basic secret-token-123",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid Authorization header format",
		},
		{
			name:         "wrong token",
			configured:   "secret-token-123",
			authHeader:   "Bearer wrong-token",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid or unauthorized API token",
		},
		{
			name:         "valid token",
			configured:   "secret-token-123",
			authHeader:   "Bearer secret-token-123",
			expectedCode: http.StatusOK,
			expectedBody: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			handler := web.NewAdminHandlerWithToken(tt.configured)

			r.POST("/admin/reload", handler.AuthMiddleware(), func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			})

			req, _ := http.NewRequest(http.MethodPost, "/admin/reload", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestAdminHandler_Reload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewAdminHandlerWithToken("test-token")
	r.POST("/reload", handler.Reload)

	req, _ := http.NewRequest(http.MethodPost, "/reload", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp["status"])
}
