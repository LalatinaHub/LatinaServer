package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
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
		case "/ping":
			c.String(http.StatusOK, "Pong")
		default:
			if proxy, err := reverse(c, "https://fool.azurewebsites.net/get"); err == nil {
				proxy.ServeHTTP(c.Writer, c.Request)
			}
		}
	})

	return r
}

func reverse(c *gin.Context, target string) (*httputil.ReverseProxy, error) {
	remote, err := url.Parse(target)
	if err != nil {
		fmt.Println(err)
		return &httputil.ReverseProxy{}, err
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
	}

	return proxy, err
}
