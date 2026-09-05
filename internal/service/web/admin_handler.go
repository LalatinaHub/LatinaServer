package web

import (
	"context"
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	apiToken string
}

func NewAdminHandler() *AdminHandler {
	// 1. Check environment variable first
	if envToken := os.Getenv("ADMIN_API_TOKEN"); envToken != "" {
		return &AdminHandler{apiToken: envToken}
	}

	// 2. Fallback to database KV store
	db, err := database.GetDB()
	if err != nil {
		return &AdminHandler{apiToken: ""}
	}

	kvRepo := repository.NewKVRepository(db)
	kvList, err := kvRepo.GetAll(context.Background())
	if err != nil {
		return &AdminHandler{apiToken: ""}
	}

	token, _ := kvList["apiToken"].(string)
	return &AdminHandler{apiToken: token}
}

func NewAdminHandlerWithToken(token string) *AdminHandler {
	return &AdminHandler{apiToken: token}
}

// AuthMiddleware validates Authorization: Bearer <token> using constant time comparison.
func (h *AdminHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.apiToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Admin API token not configured",
			})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format, expected 'Bearer <token>'",
			})
			return
		}

		token := strings.TrimSpace(parts[1])
		if subtle.ConstantTimeCompare([]byte(token), []byte(h.apiToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or unauthorized API token",
			})
			return
		}

		c.Next()
	}
}

func (h *AdminHandler) Reload(c *gin.Context) {
	systemctl.Reload(config.ServiceLatinaServer)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service reload triggered",
	})
}

func (h *AdminHandler) GetAPIToken() string {
	return h.apiToken
}

