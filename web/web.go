package web

import (
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaServer/web/reality"
	"github.com/gin-gonic/gin"
)

var (
	password = os.Getenv("PASSWORD")
)

func WebServer() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	if password == "" {
		password = "reload"
	}

	r.GET("/*path", func(c *gin.Context) {
		switch c.Param("path") {
		case "/" + password:
			config.Write()
			helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
			c.Status(http.StatusOK)
		case "/info":
			c.JSON(http.StatusOK, helper.GetIpInfo())
		case "/relay":
			c.JSON(http.StatusOK, relay.Relays)
		case "/reality":
			c.String(http.StatusOK, reality.RealityHandler())
		case "/ping":
			c.String(http.StatusOK, "Pong")
		default:
			c.Status(http.StatusSwitchingProtocols)
		}
	})

	return r
}
