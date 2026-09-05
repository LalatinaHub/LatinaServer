package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/geoip"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrialHandler struct {
	userRepo repository.UserRepository
}

func NewTrialHandler() *TrialHandler {
	return &TrialHandler{}
}

func NewTrialHandlerWithRepo(repo repository.UserRepository) *TrialHandler {
	return &TrialHandler{userRepo: repo}
}

// VMessStandardConfig defines standard V2Ray VMess JSON format.
type VMessStandardConfig struct {
	V    string `json:"v"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Port any    `json:"port"`
	ID   string `json:"id"`
	Aid  any    `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	ALPN string `json:"alpn"`
}

// TrialResponse returns all connection links and parameters for the client.
type TrialResponse struct {
	Status        string         `json:"status"`
	Message       string         `json:"message"`
	Protocol      string         `json:"protocol"`
	SpeedLimit    string         `json:"speed_limit"`
	LimitDetails  map[string]any `json:"limit_details"`
	VMessLink     string         `json:"vmess_link"`
	SingboxConfig map[string]any `json:"singbox_config"`
	ClashConfig   string         `json:"clash_config"`
	Details       map[string]any `json:"details"`
	Persisted     bool           `json:"persisted"`
}

func (h *TrialHandler) resolveHost(c *gin.Context) string {
	rawHost := c.Request.Host
	if hostPart, _, err := net.SplitHostPort(rawHost); err == nil {
		rawHost = hostPart
	}

	// If accessed through localhost / private IP, try GeoIP public IP
	if rawHost == "" || rawHost == "localhost" || rawHost == "127.0.0.1" || strings.HasPrefix(rawHost, "192.168.") || strings.HasPrefix(rawHost, "10.") {
		ipInfo := geoip.GetIpInfo()
		if ipInfo.Ip != "" {
			return ipInfo.Ip
		}
	}

	if rawHost == "" {
		return "latinaserver.local"
	}
	return rawHost
}

func (h *TrialHandler) GenerateTrial(c *gin.Context) {
	host := h.resolveHost(c)
	clientUUID := uuid.New().String()
	tokenBytes := make([]byte, 4)
	_, _ = rand.Read(tokenBytes)
	token := "trial-" + hex.EncodeToString(tokenBytes)

	now := time.Now()
	expired := now.Add(24 * time.Hour)
	expiredStr := expired.Format("2006-01-02")

	// Try persisting user to database
	persisted := false
	var repo repository.UserRepository = h.userRepo
	if repo == nil {
		if db, err := database.GetDB(); err == nil {
			repo = repository.NewUserRepository(db)
		}
	}

	if repo != nil {
		newUser := &model.User{
			Token:      token,
			Password:   clientUUID,
			Expired:    expired,
			ServerCode: host,
			Quota:      1024, // 1 GB
			Relay:      "",
			Adblock:    false,
			VPN:        "vmess",
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if _, err := repo.CreateUser(ctx, newUser); err == nil {
			persisted = true
			logger.Info().
				Str("user_id", clientUUID).
				Str("host", host).
				Msg("Trial VMess user persisted to database")

			// Trigger service reload so sing-box dynamically picks up the user
			go systemctl.Reload(config.ServiceLatinaServer)
		} else {
			logger.Warn().Err(err).Msg("Could not persist trial user to DB, continuing with generated parameters")
		}
	}

	// 1. Construct Standard V2Ray VMess URL (v2rayNG, Shadowrocket, etc.)
	v2rayObj := VMessStandardConfig{
		V:    "2",
		PS:   "LatinaServer-VMess-Trial-2Mbps",
		Add:  host,
		Port: 443,
		ID:   clientUUID,
		Aid:  0,
		Scy:  "auto",
		Net:  "ws",
		Type: "none",
		Host: host,
		Path: "/vmess",
		TLS:  "tls",
		SNI:  host,
		ALPN: "",
	}

	v2rayJSON, _ := json.Marshal(v2rayObj)
	vmessLink := "vmess://" + base64.StdEncoding.EncodeToString(v2rayJSON)

	// 2. Construct sing-box Client Outbound Configuration with TCP Brutal 2Mbps rate limit
	singboxOutbound := map[string]any{
		"type":        "vmess",
		"tag":         "LatinaServer-VMess-Trial-2Mbps",
		"server":      host,
		"server_port": 443,
		"uuid":        clientUUID,
		"security":    "auto",
		"alter_id":    0,
		"transport": map[string]any{
			"type": "ws",
			"path": "/vmess",
			"headers": map[string]string{
				"Host": host,
			},
		},
		"tls": map[string]any{
			"enabled":     true,
			"server_name": host,
			"insecure":    false,
		},
		"multiplex": map[string]any{
			"enabled":         true,
			"protocol":        "h2mux",
			"max_connections": 1,
			"min_streams":     4,
			"padding":         true,
			"brutal": map[string]any{
				"enabled":   true,
				"up_mbps":   2,
				"down_mbps": 2,
			},
		},
	}

	// 3. Construct Clash / Mihomo YAML Proxy block
	clashYaml := fmt.Sprintf(`- name: "LatinaServer-VMess-Trial-2Mbps"
  type: vmess
  server: %s
  port: 443
  uuid: %s
  alterId: 0
  cipher: auto
  udp: true
  tls: true
  servername: %s
  network: ws
  ws-opts:
    path: /vmess
    headers:
      Host: %s`, host, clientUUID, host, host)

	resp := TrialResponse{
		Status:     "success",
		Message:    "VMess trial account generated successfully with 2 Mbps speed limit",
		Protocol:   "vmess",
		SpeedLimit: "2 Mbps Up / Down",
		LimitDetails: map[string]any{
			"up_mbps":         2,
			"down_mbps":       2,
			"multiplex":       "h2mux",
			"brutal":          true,
			"congestion_ctrl": "TCP Brutal / sing-mux",
		},
		VMessLink:     vmessLink,
		SingboxConfig: singboxOutbound,
		ClashConfig:   clashYaml,
		Details: map[string]any{
			"remark":     "LatinaServer-VMess-Trial-2Mbps",
			"server":     host,
			"port":       443,
			"uuid":       clientUUID,
			"alter_id":   0,
			"security":   "auto",
			"transport":  "ws (WebSocket)",
			"path":       "/vmess",
			"tls":        "TLS enabled (port 443)",
			"sni":        host,
			"quota":      "1024 MB (1 GB)",
			"duration":   "24 Hours",
			"expired_at": expiredStr,
		},
		Persisted: persisted,
	}

	c.JSON(http.StatusOK, resp)
}
