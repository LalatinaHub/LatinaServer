package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/cloudflare"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/LalatinaHub/LatinaServer/pkg/util"
	caddy "github.com/caddyserver/caddy/v2"

	_ "github.com/caddy-dns/cloudflare"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/mholt/caddy-l4"
)

// ReadCaddyConfig loads and unmarshals Caddy configuration from JSON file.

func ReadCaddyConfig(configLocation string) (*caddy.Config, error) {
	defer systemctl.CatchError(true)

	body, err := os.ReadFile(configLocation)
	if err != nil {
		return nil, appErrors.NewConfigError("failed to read caddy config", err)
	}

	var caddyConfig caddy.Config
	if err = json.Unmarshal(body, &caddyConfig); err != nil {
		return nil, appErrors.NewConfigError("failed to unmarshal caddy config", err)
	}

	return &caddyConfig, nil
}
// GenerateCaddyConfig creates production Caddy JSON configuration based on DB records.


func GenerateCaddyConfig() error {
	db, err := database.GetDB()
	if err != nil {
		return appErrors.NewDatabaseError("failed to get database connection", err)
	}

	kvRepo := repository.NewKVRepository(db)
	kvList, err := kvRepo.GetAll(context.Background())
	if err != nil {
		return appErrors.NewDatabaseError("failed to get KV list", err)
	}

	cfClient, err := cloudflare.NewClient()
	if err != nil {
		return appErrors.NewNetworkError("failed to create cloudflare client", err)
	}

	domains, err := cfClient.GetAssignedDomain()
	if err != nil {
		return appErrors.NewNetworkError("failed to get assigned domains", err)
	}

	var (
		domain = domains[0]
		cfKey  = kvList["cfKey"].(string)
		email  = kvList["email"].(string)

		stringCaddyConfig string
	)

	buf, err := os.ReadFile(CaddyConfigPath)
	if err != nil {
		return appErrors.NewConfigError("failed to read caddy config template", err)
	}

	// Manipulate domains
	for i := range domains {
		domains[i] = fmt.Sprintf(`"%s"`, domains[i])
	}

	// Remove duplicated domains
	domains = util.RemoveDuplicate(domains)

	stringCaddyConfig = string(buf)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, `"DOMAIN_LIST"`, strings.Join(domains, ","))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "DOMAIN", domain)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CF_KEY", cfKey)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "EMAIL", email)
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CADDY_LOGFILE", CaddyLogPath)

	// Port of
	// Trojan
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_TCP_PORT", fmt.Sprint(TrojanTCPPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_WS_PORT", fmt.Sprint(TrojanWSPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_HU_PORT", fmt.Sprint(TrojanHUPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_GRPC_PORT", fmt.Sprint(TrojanGRPCPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "TROJAN_UDP_PORT", fmt.Sprint(TrojanUDPPort))

	// VMess
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_TCP_PORT", fmt.Sprint(VMessTCPPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_WS_PORT", fmt.Sprint(VMessWSPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_HU_PORT", fmt.Sprint(VMessHUPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_GRPC_PORT", fmt.Sprint(VMessGRPCPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VMESS_UDP_PORT", fmt.Sprint(VMessUDPPort))

	// VLESS
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_TCP_PORT", fmt.Sprint(VLessTCPPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_WS_PORT", fmt.Sprint(VLessWSPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_HU_PORT", fmt.Sprint(VLessHUPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_GRPC_PORT", fmt.Sprint(VLessGRPCPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "VLESS_UDP_PORT", fmt.Sprint(VLessUDPPort))

	// Services
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "CADDY_PORT", fmt.Sprint(CaddyPort))
	stringCaddyConfig = strings.ReplaceAll(stringCaddyConfig, "WEBSERVER_PORT", fmt.Sprint(WebServerPort))

	buf = []byte(stringCaddyConfig)
	var caddyConfig caddy.Config
	if err = json.Unmarshal(buf, &caddyConfig); err != nil {
		return appErrors.NewConfigError("failed to unmarshal generated caddy config", err)
	}

	changed, err := SaveJsonToFileWithCache(CaddyActiveConfigPath, caddyConfig)
	if err != nil {
		return appErrors.NewConfigError("failed to save caddy config", err)
	}
	if changed {
		logger.Info().Msg("Caddy config updated (hash changed)")
	} else {
		logger.Debug().Msg("Caddy config unchanged (hash match), skipped write")
	}

	return nil
}

