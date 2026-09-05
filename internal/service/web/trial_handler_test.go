package web_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	repoMocks "github.com/LalatinaHub/LatinaServer/internal/repository/mocks"
	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTrialHandler_GenerateTrial_WithoutRepo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := web.NewTrialHandler()
	r.POST("/api/v1/trial", handler.GenerateTrial)
	r.GET("/api/v1/trial", handler.GenerateTrial)

	// Test POST
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trial", nil)
	req.Host = "example.com:443"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp web.TrialResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "vmess", resp.Protocol)
	assert.Equal(t, "2 Mbps Up / Down", resp.SpeedLimit)

	// Verify limit details has brutal 2 mbps
	assert.Equal(t, float64(2), resp.LimitDetails["up_mbps"])
	assert.Equal(t, float64(2), resp.LimitDetails["down_mbps"])

	// Verify VMess Link format
	assert.True(t, strings.HasPrefix(resp.VMessLink, "vmess://"))
	b64Data := strings.TrimPrefix(resp.VMessLink, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	assert.NoError(t, err)

	var vmessObj web.VMessStandardConfig
	err = json.Unmarshal(decoded, &vmessObj)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", vmessObj.Add)
	assert.Equal(t, "ws", vmessObj.Net)
	assert.Equal(t, "/vmess", vmessObj.Path)
	assert.Equal(t, "tls", vmessObj.TLS)

	// Verify sing-box config contains multiplex with brutal rate limiting
	singboxCfg := resp.SingboxConfig
	assert.Equal(t, "vmess", singboxCfg["type"])
	multiplex, ok := singboxCfg["multiplex"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, true, multiplex["enabled"])
	brutal, ok := multiplex["brutal"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, true, brutal["enabled"])
	assert.Equal(t, float64(2), brutal["up_mbps"])
	assert.Equal(t, float64(2), brutal["down_mbps"])

	// Verify Clash config
	assert.Contains(t, resp.ClashConfig, "type: vmess")
	assert.Contains(t, resp.ClashConfig, "network: ws")

	// Test GET endpoint as well
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/trial", nil)
	reqGet.Host = "test.domain"
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)
}

func TestTrialHandler_GenerateTrial_WithRepo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockRepo := new(repoMocks.MockUserRepository)
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.VPN == "vmess" && u.Quota == 1024
	})).Return(int64(1), nil)

	handler := web.NewTrialHandlerWithRepo(mockRepo)
	r.POST("/api/v1/trial", handler.GenerateTrial)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trial", nil)
	req.Host = "server.example.org"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp web.TrialResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Persisted)
	mockRepo.AssertExpectations(t)
}
