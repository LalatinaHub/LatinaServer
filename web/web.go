package web

import (
	"net/http"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	web_helper "github.com/LalatinaHub/LatinaServer/web/helper"
	"github.com/gin-gonic/gin"
)

type UDPMessage struct {
	Target string `json:"target"`
	Data   string `json:"data"`
}

func WebServer() http.Handler {
	var (
		kvList = db.GetKVList()
		r      = gin.Default()
	)

	r.Use(gin.Recovery())

	r.POST("/api/v1/:path", func(c *gin.Context) {
		switch c.Param("path") {
		case "udp-proxy":
			proxyRequest := web_helper.ProxyRequest{}
			c.BindJSON(&proxyRequest)
			resBuffer := web_helper.HandleUDPForwarding(proxyRequest)
			c.Data(http.StatusOK, "application/x-binary", resBuffer)
		}
	})

	r.GET("/api/v1/:path", func(c *gin.Context) {
		switch c.Param("path") {
		case kvList["apiToken"]:
			helper.ReloadService([]string{CS.SERVICE_LATINASERVER}...)
			c.Status(http.StatusOK)
		case "check":
			proxy := c.Query("ip")
			if proxy == "" {
				c.String(http.StatusBadRequest, "No proxy provided!")
				return
			}

			proxyIP, err := helper.CheckProxyIP(proxy)
			if err != nil {
				c.String(500, err.Error())
				return
			}
			c.JSON(http.StatusOK, proxyIP)
		case "info":
			c.JSON(http.StatusOK, helper.GetIpInfo())
		case "relay":
			c.JSON(http.StatusOK, relay.Relays)
		case "ping":
			c.String(http.StatusOK, "Pong")
		case "status":
			c.JSON(http.StatusOK, helper.GetServerStatus())
		default:
			c.String(http.StatusOK, "Welcome to Gin!")
		}
	})

	return r
}
