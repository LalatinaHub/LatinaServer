package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/gin-gonic/gin"
)

// UserSettingsHandler handles client self-service user configuration.
type UserSettingsHandler struct {
	repo repository.UserSettingsRepository
}

// NewUserSettingsHandler creates a new UserSettingsHandler.
func NewUserSettingsHandler() *UserSettingsHandler {
	return &UserSettingsHandler{}
}

// NewUserSettingsHandlerWithRepo creates a UserSettingsHandler with a custom repository (for testing).
func NewUserSettingsHandlerWithRepo(repo repository.UserSettingsRepository) *UserSettingsHandler {
	return &UserSettingsHandler{repo: repo}
}

func (h *UserSettingsHandler) getRepo() (repository.UserSettingsRepository, error) {
	if h.repo != nil {
		return h.repo, nil
	}
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}
	return repository.NewUserSettingsRepository(db), nil
}

// UpdateUserSettingsRequest represents the payload to update user settings.
type UpdateUserSettingsRequest struct {
	Token   string `json:"token" form:"token"`
	Adblock *bool  `json:"adblock" binding:"required"`
}

// UserSettingsResponse wraps API response for settings.
type UserSettingsResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *UserSettingsData `json:"data,omitempty"`
}

// UserSettingsData contains user preference and status info.
type UserSettingsData struct {
	ID         int64  `json:"id"`
	Token      string `json:"token"`
	ServerCode string `json:"server_code"`
	VPN        string `json:"vpn"`
	Quota      int64  `json:"quota"`
	Expired    string `json:"expired"`
	Adblock    bool   `json:"adblock"`
}

func extractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}

	return ""
}

// GetSettings retrieves user settings and adblock status.
func (h *UserSettingsHandler) GetSettings(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing token: provide token query parameter or Authorization Bearer header",
		})
		return
	}

	repo, err := h.getRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database connection unavailable",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := repo.GetUserByToken(ctx, token)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "User not found or invalid token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch user settings",
		})
		return
	}

	c.JSON(http.StatusOK, UserSettingsResponse{
		Success: true,
		Message: "User settings fetched successfully",
		Data: &UserSettingsData{
			ID:         user.ID,
			Token:      user.Token,
			ServerCode: user.ServerCode,
			VPN:        user.VPN,
			Quota:      user.Quota,
			Expired:    user.Expired.Format(time.RFC3339),
			Adblock:    user.Adblock,
		},
	})
}

// UpdateSettings updates user settings such as Adblock toggle.
func (h *UserSettingsHandler) UpdateSettings(c *gin.Context) {
	var req UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload: 'adblock' boolean is required",
		})
		return
	}

	token := req.Token
	if token == "" {
		token = extractToken(c)
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing token: provide token in JSON body, query parameter, or Bearer header",
		})
		return
	}

	repo, err := h.getRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database connection unavailable",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := repo.GetUserByToken(ctx, token)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "User not found or invalid token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to find user",
		})
		return
	}

	if err := repo.UpdateUserAdblock(ctx, token, *req.Adblock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update adblock setting",
		})
		return
	}

	// Re-generate sing-box configuration with updated AdBlock rules
	if genErr := config.GenerateSingConfig(); genErr != nil {
		logger.Warn().Err(genErr).Msg("Failed to re-generate sing-box config after user settings change")
	} else {
		logger.Info().Str("token", token).Bool("adblock", *req.Adblock).Msg("User adblock setting updated, sing-box config regenerated")
		systemctl.Reload(config.ServiceLatinaServer)
	}

	c.JSON(http.StatusOK, UserSettingsResponse{
		Success: true,
		Message: "User settings updated successfully",
		Data: &UserSettingsData{
			ID:         user.ID,
			Token:      user.Token,
			ServerCode: user.ServerCode,
			VPN:        user.VPN,
			Quota:      user.Quota,
			Expired:    user.Expired.Format(time.RFC3339),
			Adblock:    *req.Adblock,
		},
	})
}
