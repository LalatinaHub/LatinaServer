package web

import (
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaSub-go/geoip"
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

	r.GET("/api/v1/:path", func(c *gin.Context) {
		switch c.Param("path") {
		case password:
			helper.ReloadService([]string{CS.ServiceLatinaServer}...)
			c.Status(http.StatusOK)
		case "info":
			c.JSON(http.StatusOK, helper.GetIpInfo())
		case "relay":
			c.JSON(http.StatusOK, relay.Relays)
		case "ping":
			c.String(http.StatusOK, "Pong")
		case "myip":
			c.JSON(http.StatusOK, geoip.MyIp{
				Ip: c.ClientIP(),
			})
		case "status":
			c.JSON(http.StatusOK, helper.GetServerStatus())
		default:
			c.String(http.StatusOK, "Welcome to Gin!")
		}
	})

	return r
}
