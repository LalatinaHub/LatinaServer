package config

import (
	"encoding/json"
	"os"
	"strings"

	_ "github.com/caddy-dns/cloudflare"
	caddy "github.com/caddyserver/caddy/v2"
	_ "github.com/mholt/caddy-l4"
)

func LoadCaddyConfig() *caddy.Config {
	var (
		domain = os.Getenv("DOMAIN")
		cfKey  = os.Getenv("CF_KEY")
		email  = os.Getenv("EMAIL")

		stringCaddyConfig string
	)

	buf, err := os.ReadFile("/usr/local/etc/latinaserver/caddy.json")
	if err != nil {
		panic(err)
	}

	stringCaddyConfig = string(buf)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "DOMAIN", domain)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CF_KEY", cfKey)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "EMAIL", email)

	buf = []byte(stringCaddyConfig)
	var caddyConfig caddy.Config
	if err = json.Unmarshal(buf, &caddyConfig); err != nil {
		panic(err)
	}

	return &caddyConfig
}
