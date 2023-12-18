package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaSub-go/provider"
	"github.com/dickymuliafiqri/BenchBox/modules/benchmark"
	singboxBench "github.com/dickymuliafiqri/BenchBox/modules/sing-box"
	"github.com/dickymuliafiqri/BenchBox/server/api/bench"
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
			helper.ReloadService([]string{CS.ServiceLatinaServer}...)
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

	r.POST("/*path", func(c *gin.Context) {
		switch c.Param("path") {
		case "/bench":
			node := c.PostForm("url")

			result, err := benchAccount(node)
			if err != nil || result == nil {
				c.String(http.StatusInternalServerError, err.Error())
			}

			c.JSON(http.StatusOK, result)
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

func benchAccount(node string) ([]bench.ResultType, error) {
	node = strings.ReplaceAll(node, ",", "\n")
	result := []bench.ResultType{}
	outbounds, err := provider.Parse(node)
	if err != nil {
		return nil, err
	}

	for _, outbound := range outbounds {
		opt, listenPort := singboxBench.GenerateConfig(&outbound)
		box, err := singboxBench.Create(opt)
		if err != nil {
			return nil, err
		}
		defer box.Close()

		time.Sleep(1 * time.Second)

		result = append(result, bench.ResultType{
			Node:   outbound.Tag,
			Result: benchmark.StartBenchmark(listenPort),
		})
	}

	return result, err
}
