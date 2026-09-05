package web

import (
	"net/http"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/util"
	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	cache *MemoryCache
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{
		cache: NewMemoryCache(),
	}
}

func (h *StatusHandler) Status(c *gin.Context) {
	const cacheKey = "server_status"
	if cached, ok := h.cache.Get(cacheKey); ok {
		c.JSON(http.StatusOK, cached)
		return
	}

	status := util.GetServerStatus()
	h.cache.Set(cacheKey, status, 30*time.Second)
	c.JSON(http.StatusOK, status)
}

