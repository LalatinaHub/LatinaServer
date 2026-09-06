package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserSettingsRepo struct {
	users map[string]*model.User
}

func (m *mockUserSettingsRepo) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	if u, ok := m.users[token]; ok {
		return u, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserSettingsRepo) UpdateUserAdblock(ctx context.Context, token string, adblock bool) error {
	if u, ok := m.users[token]; ok {
		u.Adblock = adblock
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *mockUserSettingsRepo) UpdateUserUUID(ctx context.Context, token string, newUUID string) error {
	if u, ok := m.users[token]; ok {
		u.Password = newUUID
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *mockUserSettingsRepo) UpdateUserToken(ctx context.Context, oldToken string, newToken string) error {
	if u, ok := m.users[oldToken]; ok {
		delete(m.users, oldToken)
		u.Token = newToken
		m.users[newToken] = u
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *mockUserSettingsRepo) UpdateUserConfig(ctx context.Context, token string, vpn string, serverCode string, relay string) error {
	if u, ok := m.users[token]; ok {
		u.VPN = vpn
		u.ServerCode = serverCode
		u.Relay = relay
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *mockUserSettingsRepo) GetServers(ctx context.Context) ([]model.Server, error) {
	return []model.Server{
		{ID: 1, Code: "SG1", Domain: "sg1.test", Country: "SG", UsersCount: 5, UsersMax: 100},
	}, nil
}

func (m *mockUserSettingsRepo) GetWildcards(ctx context.Context) ([]string, error) {
	return []string{"quiz.vidio.com"}, nil
}

func setupUserSettingsRouter(repo repository.UserSettingsRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := web.NewUserSettingsHandlerWithRepo(repo)
	api := r.Group("/api/v1")
	{
		api.GET("/user/settings", handler.GetSettings)
		api.POST("/user/settings", handler.UpdateSettings)
	}
	return r
}

func TestUserSettingsHandler_GetSettings(t *testing.T) {
	repo := &mockUserSettingsRepo{
		users: map[string]*model.User{
			"valid-token": {
				ID:         42,
				Token:      "valid-token",
				ServerCode: "SG",
				VPN:        "vmess",
				Quota:      5000000,
				Expired:    time.Now().Add(24 * time.Hour),
				Adblock:    true,
			},
		},
	}
	r := setupUserSettingsRouter(repo)

	t.Run("missing token returns 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/settings", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Missing token")
	})

	t.Run("user not found returns 404", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/settings?token=non-existent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})

	t.Run("success fetching user via query token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/settings?token=valid-token", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp web.UserSettingsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, int64(42), resp.Data.ID)
		assert.True(t, resp.Data.Adblock)
		assert.Equal(t, "SG", resp.Data.ServerCode)
	})

	t.Run("success fetching user via Bearer header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/settings", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp web.UserSettingsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.True(t, resp.Data.Adblock)
	})
}

func TestUserSettingsHandler_UpdateSettings(t *testing.T) {
	repo := &mockUserSettingsRepo{
		users: map[string]*model.User{
			"toggle-token": {
				ID:         100,
				Token:      "toggle-token",
				ServerCode: "ID",
				VPN:        "trojan",
				Quota:      10000000,
				Expired:    time.Now().Add(48 * time.Hour),
				Adblock:    false,
			},
		},
	}
	r := setupUserSettingsRouter(repo)

	t.Run("missing adblock field returns 400", func(t *testing.T) {
		payload := []byte(`{"token": "toggle-token"}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/settings", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing token returns 400", func(t *testing.T) {
		payload := []byte(`{"adblock": true}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/settings", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Missing token")
	})

	t.Run("user not found returns 404", func(t *testing.T) {
		payload := []byte(`{"token": "unknown", "adblock": true}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/settings", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success updating adblock to true", func(t *testing.T) {
		payload := []byte(`{"token": "toggle-token", "adblock": true}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/settings", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp web.UserSettingsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.True(t, resp.Data.Adblock)
		assert.True(t, repo.users["toggle-token"].Adblock)
	})

	t.Run("success updating adblock via query token and body", func(t *testing.T) {
		payload := []byte(`{"adblock": false}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/settings?token=toggle-token", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp web.UserSettingsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.False(t, resp.Data.Adblock)
		assert.False(t, repo.users["toggle-token"].Adblock)
	})
}
