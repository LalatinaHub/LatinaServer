package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type portalMockRepo struct {
	mu        sync.RWMutex
	users     map[string]*model.User
	servers   []model.Server
	wildcards []string
}

func (m *portalMockRepo) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if u, ok := m.users[token]; ok {
		return u, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *portalMockRepo) UpdateUserAdblock(ctx context.Context, token string, adblock bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[token]; ok {
		u.Adblock = adblock
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *portalMockRepo) UpdateUserUUID(ctx context.Context, token string, newUUID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[token]; ok {
		u.Password = newUUID
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *portalMockRepo) UpdateUserToken(ctx context.Context, oldToken string, newToken string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[oldToken]; ok {
		delete(m.users, oldToken)
		u.Token = newToken
		m.users[newToken] = u
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *portalMockRepo) UpdateUserConfig(ctx context.Context, token string, vpn string, serverCode string, relay string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[token]; ok {
		u.VPN = vpn
		u.ServerCode = serverCode
		u.Relay = relay
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *portalMockRepo) GetServers(ctx context.Context) ([]model.Server, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.servers, nil
}

func (m *portalMockRepo) GetWildcards(ctx context.Context) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.wildcards, nil
}

func (m *portalMockRepo) getUser(token string) *model.User {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.users[token]
}

func setupPortalRouter(repo repository.UserSettingsRepository, customPath string) (*gin.Engine, *web.PortalHandler, *atomic.Bool, *atomic.Bool) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	var genCalled atomic.Bool
	var reloadCalled atomic.Bool

	singGen := func() error {
		genCalled.Store(true)
		return nil
	}
	reload := func() error {
		reloadCalled.Store(true)
		return nil
	}

	handler := web.NewPortalHandlerWithRepo(repo, singGen, reload, customPath)

	r.GET("/portal", handler.ServePortal)

	api := r.Group("/api/v1/portal")
	{
		api.GET("/profile", handler.GetProfile)
		api.POST("/reset-uuid", handler.ResetUUID)
		api.POST("/change-token", handler.ChangeToken)
		api.POST("/update-config", handler.UpdateConfig)
		api.POST("/toggle-adblock", handler.ToggleAdblock)
		api.GET("/servers", handler.GetServers)
		api.GET("/relays", handler.GetRelays)
		api.GET("/status", handler.GetStatus)
		api.GET("/wildcards", handler.GetWildcards)
		api.GET("/info", handler.GetInfo)
	}

	return r, handler, &genCalled, &reloadCalled
}

func TestPortalHandler_GetProfile(t *testing.T) {
	futureDate := time.Now().AddDate(0, 1, 0)
	repo := &portalMockRepo{
		users: map[string]*model.User{
			"valid-user": {
				ID:         100,
				Token:      "valid-user",
				Password:   "uuid-test-1234",
				Expired:    futureDate,
				ServerCode: "SG1",
				Quota:      5368709120, // 5GB
				Relay:      "SG",
				Adblock:    true,
				VPN:        "vless",
			},
		},
		servers: []model.Server{
			{ID: 1, Code: "SG1", Domain: "sg1.example.com", Country: "SG", UsersCount: 10, UsersMax: 100},
		},
		wildcards: []string{"bug.vidio.com"},
	}

	router, _, _, _ := setupPortalRouter(repo, "")

	t.Run("missing token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/profile", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unknown user returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/profile?token=ghost-token", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("success via query param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/profile?token=valid-user", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp web.ProfileResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, int64(100), resp.ID)
		assert.Equal(t, "valid-user", resp.Token)
		assert.Equal(t, "uuid-test-1234", resp.Password)
		assert.Equal(t, "vless", resp.VPN)
		assert.True(t, resp.Adblock)
		assert.True(t, resp.IsDonator)
		assert.False(t, resp.IsExpired)
		assert.Contains(t, resp.QuotaHuman, "5.00 GB")
		assert.Contains(t, resp.SubURL, "/sub?token=valid-user")
		assert.Contains(t, resp.ClashURL, "clash://install-config")
		assert.Contains(t, resp.SingboxURL, "sing-box://import-remote-profile")
		assert.Contains(t, resp.ShadowrocketURL, "shadowrocket://add/sub://")
		assert.Contains(t, resp.RawURI, "vless://uuid-test-1234@sg1.example.com:443")
		assert.Len(t, resp.Servers, 1)
		assert.Len(t, resp.Wildcards, 1)
	})

	t.Run("success via Authorization Bearer header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/profile", nil)
		req.Header.Set("Authorization", "Bearer valid-user")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestPortalHandler_ResetUUID(t *testing.T) {
	repo := &portalMockRepo{
		users: map[string]*model.User{
			"user-token": {
				ID:       101,
				Token:    "user-token",
				Password: "old-uuid",
			},
		},
	}

	router, handler, genCalled, reloadCalled := setupPortalRouter(repo, "")

	t.Run("missing token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/reset-uuid", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success reset uuid via query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/reset-uuid?token=user-token", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res map[string]any
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)

		newPassword := res["password"].(string)
		assert.NotEmpty(t, newPassword)
		assert.NotEqual(t, "old-uuid", newPassword)

		// Wait deterministically for async reload
		handler.WaitForAsyncReload()
		assert.True(t, genCalled.Load())
		assert.True(t, reloadCalled.Load())
	})
}

func TestPortalHandler_ChangeToken(t *testing.T) {
	repo := &portalMockRepo{
		users: map[string]*model.User{
			"token-old": {
				ID:    102,
				Token: "token-old",
			},
		},
	}

	router, _, _, _ := setupPortalRouter(repo, "")

	t.Run("missing token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/change-token", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success change token", func(t *testing.T) {
		body := []byte(`{"token": "token-old"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/change-token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res map[string]any
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)

		newToken := res["token"].(string)
		assert.Len(t, newToken, 8)
		assert.NotEqual(t, "token-old", newToken)
	})
}

func TestPortalHandler_UpdateConfig(t *testing.T) {
	repo := &portalMockRepo{
		users: map[string]*model.User{
			"user-config-tok": {
				ID:         103,
				Token:      "user-config-tok",
				VPN:        "vmess",
				ServerCode: "ID1",
				Relay:      "",
			},
		},
	}

	router, handler, _, _ := setupPortalRouter(repo, "")
	defer handler.WaitForAsyncReload()

	t.Run("invalid protocol returns 400", func(t *testing.T) {
		body := []byte(`{"token": "user-config-tok", "vpn": "wireguard", "server_code": "SG1"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/update-config", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("success update config", func(t *testing.T) {
		body := []byte(`{"token": "user-config-tok", "vpn": "trojan", "server_code": "SG2", "relay": "Tanpa Relay"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/update-config", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		u := repo.getUser("user-config-tok")
		assert.Equal(t, "trojan", u.VPN)
		assert.Equal(t, "SG2", u.ServerCode)
		assert.Equal(t, "", u.Relay) // "Tanpa Relay" normalized to ""
	})
}

func TestPortalHandler_ToggleAdblock(t *testing.T) {
	repo := &portalMockRepo{
		users: map[string]*model.User{
			"user-adblock-tok": {
				ID:      104,
				Token:   "user-adblock-tok",
				Adblock: false,
			},
		},
	}

	router, handler, _, _ := setupPortalRouter(repo, "")
	defer handler.WaitForAsyncReload()

	t.Run("toggle to true", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/toggle-adblock?token=user-adblock-tok", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, repo.getUser("user-adblock-tok").Adblock)
	})

	t.Run("toggle to false explicitly", func(t *testing.T) {
		body := []byte(`{"token": "user-adblock-tok", "adblock": false}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/portal/toggle-adblock", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.False(t, repo.getUser("user-adblock-tok").Adblock)
	})
}

func TestPortalHandler_MetadataEndpoints(t *testing.T) {
	repo := &portalMockRepo{
		servers: []model.Server{
			{ID: 1, Code: "SG1", Domain: "sg1.test", Country: "SG", UsersCount: 1, UsersMax: 50},
		},
		wildcards: []string{"bug.domain.com"},
	}

	router, _, _, _ := setupPortalRouter(repo, "")

	t.Run("get servers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/servers", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "SG1")
	})

	t.Run("get relays", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/relays", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Tanpa Relay")
	})

	t.Run("get status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/status", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("get wildcards", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/wildcards", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "bug.domain.com")
	})

	t.Run("get info", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/info", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "t.me/foolvpn")
		assert.Contains(t, rec.Body.String(), "trakteer.id")
	})
}

func TestPortalHandler_ServePortal(t *testing.T) {
	tmpDir := t.TempDir()
	portalHTML := filepath.Join(tmpDir, "portal.html")
	err := os.WriteFile(portalHTML, []byte("<!DOCTYPE html><html><body>Portal Test</body></html>"), 0644)
	require.NoError(t, err)

	repo := &portalMockRepo{}
	router, _, _, _ := setupPortalRouter(repo, portalHTML)

	req := httptest.NewRequest(http.MethodGet, "/portal", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, rec.Body.String(), "Portal Test")
}
