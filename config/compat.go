package config

import (
	intConfig "github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	caddy "github.com/caddyserver/caddy/v2"
	"github.com/sagernet/sing-box/option"
)

// ReadCaddyConfig - backward compatibility wrapper
func ReadCaddyConfig(configLocation string) *caddy.Config {
	cfg, err := intConfig.ReadCaddyConfig(configLocation)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read caddy config")
	}
	return cfg
}

// GenerateCaddyConfig - backward compatibility wrapper
func GenerateCaddyConfig() {
	if err := intConfig.GenerateCaddyConfig(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate caddy config")
	}
}

// ReadSingConfig - backward compatibility wrapper
func ReadSingConfig(configLocation string) option.Options {
	opts, err := intConfig.ReadSingConfig(configLocation)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read sing config")
	}
	return opts
}

// GenerateSingConfig - backward compatibility wrapper
func GenerateSingConfig() {
	if err := intConfig.GenerateSingConfig(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate sing config")
	}
}

// SaveJsonToFile - backward compatibility wrapper
func SaveJsonToFile(filename string, content any) {
	if err := intConfig.SaveJsonToFile(filename, content); err != nil {
		logger.Fatal().Err(err).Msg("Failed to save JSON to file")
	}
}

