package web

import (
	"net/http"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/geoip"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	cache *MemoryCache
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		cache: NewMemoryCache(),
	}
}

func (h *HealthHandler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func (h *HealthHandler) Info(c *gin.Context) {
	const cacheKey = "ip_info"
	if cached, ok := h.cache.Get(cacheKey); ok {
		c.JSON(http.StatusOK, cached)
		return
	}

	info := geoip.GetIpInfo()
	h.cache.Set(cacheKey, info, 5*time.Minute)
	c.JSON(http.StatusOK, info)
}

