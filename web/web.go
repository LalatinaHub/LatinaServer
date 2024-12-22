package web

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaSub-go/geoip"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	password = os.Getenv("PASSWORD")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type UDPMessage struct {
	Target string `json:"target"`
	Data   string `json:"data"`
}

func WebServer() http.Handler {
	r := gin.Default()
	r.Use(gin.Recovery())

	if password == "" {
		password = "reload"
	}

	r.GET("/api/v1/:path", func(c *gin.Context) {
		switch c.Param("path") {
		case password:
			helper.ReloadService([]string{CS.SERVICE_LATINASERVER}...)
			c.Status(http.StatusOK)
		case "udp-proxy":
			udpProxyHandler(c)
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

func udpProxyHandler(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer wsConn.Close()

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from WebSocket: %v", err)
			break
		}

		var msg UDPMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		udpAddr, err := net.ResolveUDPAddr("udp", msg.Target)
		if err != nil {
			log.Printf("Failed to resolve UDP address: %v", err)
			continue
		}

		udpConn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			log.Printf("Failed to connect to UDP server: %v", err)
			continue
		}
		defer udpConn.Close()

		_, err = udpConn.Write([]byte(msg.Data))
		if err != nil {
			log.Printf("Error sending to UDP server: %v", err)
		}
	}
}
