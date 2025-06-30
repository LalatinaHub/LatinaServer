package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LalatinaHub/LatinaServer/cf"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	caddy "github.com/caddyserver/caddy/v2"

	_ "github.com/caddy-dns/cloudflare"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
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
		kvList  = db.GetKVList()
		domains = cf.MakeCloudflareAPIClient().GetAssignedDomain()
		domain  = domains[0]
		cfKey   = kvList["cfKey"].(string)
		email   = kvList["email"].(string)

		stringCaddyConfig string
	)

	buf, err := os.ReadFile(CS.CADDY_CONFIG_PATH)
	if err != nil {
		panic(err)
	}

	// Manipulate domains
	for i := range domains {
		domains[i] = fmt.Sprintf(`"%s"`, domains[i])
	}

	// Remove duplicated domains
	domains = helper.RemoveDuplicate(domains)

	stringCaddyConfig = string(buf)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, `"DOMAIN_LIST"`, strings.Join(domains, ","))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "DOMAIN", domain)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CF_KEY", cfKey)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "EMAIL", email)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CADDY_LOGFILE", CS.CADDY_LOG_PATH)

	// Port of
	// Trojan
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_TCP_PORT", fmt.Sprint(CS.TROJAN_TCP_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_WS_PORT", fmt.Sprint(CS.TROJAN_WS_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_HU_PORT", fmt.Sprint(CS.TROJAN_HU_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_GRPC_PORT", fmt.Sprint(CS.TROJAN_GRPC_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_UDP_PORT", fmt.Sprint(CS.TROJAN_UDP_PORT))

	// VMess
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_TCP_PORT", fmt.Sprint(CS.VMESS_TCP_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_WS_PORT", fmt.Sprint(CS.VMESS_WS_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_HU_PORT", fmt.Sprint(CS.VMESS_HU_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_GRPC_PORT", fmt.Sprint(CS.VMESS_GRPC_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_UDP_PORT", fmt.Sprint(CS.VMESS_UDP_PORT))

	// VLESS
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_TCP_PORT", fmt.Sprint(CS.VLESS_TCP_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_WS_PORT", fmt.Sprint(CS.VLESS_WS_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_HU_PORT", fmt.Sprint(CS.VLESS_HU_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_GRPC_PORT", fmt.Sprint(CS.VLESS_GRPC_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_UDP_PORT", fmt.Sprint(CS.VLESS_UDP_PORT))

	// Services
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CADDY_PORT", fmt.Sprint(CS.CADDY_PORT))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "WEBSERVER_PORT", fmt.Sprint(CS.WEBSERVER_PORT))

	buf = []byte(stringCaddyConfig)
	var caddyConfig caddy.Config
	if err = json.Unmarshal(buf, &caddyConfig); err != nil {
		panic(err)
	}

	SaveJsonToFile(CS.CADDY_ACTIVE_CONFIG_PATH, caddyConfig)
}
