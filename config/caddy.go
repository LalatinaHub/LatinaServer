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
		caddyConfigPath   string = "/usr/local/etc/latinaserver/caddy.json"
	)

	buf, err := os.ReadFile(caddyConfigPath)
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

	// Write edited config
	f, err := os.Create(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if b, err := json.MarshalIndent(caddyConfig, "", "\t"); err == nil {
		f.WriteString(string(b))
	}

	return &caddyConfig
}
