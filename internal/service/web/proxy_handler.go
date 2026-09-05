package web

import (
	"net/http"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	"github.com/LalatinaHub/LatinaServer/pkg/util"
	"github.com/gin-gonic/gin"
)

type ProxyHandler struct{}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{}
}

func (h *ProxyHandler) Check(c *gin.Context) {
	proxy := c.Query("ip")
	if proxy == "" {
		c.String(http.StatusBadRequest, "No proxy provided!")
		return
	}

	proxyIP, err := util.CheckProxyIP(proxy)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, proxyIP)
}

func (h *ProxyHandler) Relays(c *gin.Context) {
	c.JSON(http.StatusOK, relay.GetRelays())
}
