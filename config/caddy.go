package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	caddy "github.com/caddyserver/caddy/v2"

	_ "github.com/caddy-dns/cloudflare"
	_ "github.com/mholt/caddy-l4"
)

func ReadCaddyConfig(configLocation string) *caddy.Config {
	defer helper.CatchError(true)

	body, err := os.ReadFile(configLocation)
	if err != nil {
		panic(err)
	}

	var caddyConfig caddy.Config
	if err = json.Unmarshal(body, &caddyConfig); err != nil {
		panic(err)
	}

	return &caddyConfig
}

func GenerateCaddyConfig() {
	var (
		domain = os.Getenv("DOMAIN")
		cfKey  = os.Getenv("CF_KEY")
		email  = os.Getenv("EMAIL")

		stringCaddyConfig string
	)

	buf, err := os.ReadFile(CS.CaddyConfigPath)
	if err != nil {
		panic(err)
	}

	stringCaddyConfig = string(buf)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "DOMAIN", domain)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CF_KEY", cfKey)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "EMAIL", email)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CADDY_LOGFILE", CS.CaddyLogPath)

	// Port of
	// Trojan
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_TCP_PORT", fmt.Sprint(CS.TrojanTCPPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_WS_PORT", fmt.Sprint(CS.TrojanWSPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_HU_PORT", fmt.Sprint(CS.TrojanHUPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_GRPC_PORT", fmt.Sprint(CS.TrojanGRPCPort))

	// VMess
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_TCP_PORT", fmt.Sprint(CS.VMessTCPPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_WS_PORT", fmt.Sprint(CS.VMessWSPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_HU_PORT", fmt.Sprint(CS.VMessHUPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_GRPC_PORT", fmt.Sprint(CS.VMessGRPCPort))

	// VLESS
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_TCP_PORT", fmt.Sprint(CS.VLESSTCPPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_WS_PORT", fmt.Sprint(CS.VLESSWSPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_HU_PORT", fmt.Sprint(CS.VLESSHUPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_GRPC_PORT", fmt.Sprint(CS.VLESSGRPCPort))

	// Services
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CADDY_PORT", fmt.Sprint(CS.CaddyPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "WEBSERVER_PORT", fmt.Sprint(CS.WebServerPort))

	buf = []byte(stringCaddyConfig)
	var caddyConfig caddy.Config
	if err = json.Unmarshal(buf, &caddyConfig); err != nil {
		panic(err)
	}

	SaveJsonToFile(CS.CaddyActiveConfigPath, caddyConfig)
}
