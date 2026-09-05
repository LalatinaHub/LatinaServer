package caddy

import (
	"context"
	"os"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	caddy "github.com/caddyserver/caddy/v2"
)

// RunWithContext generates configuration, starts Caddy service, and shuts down gracefully on ctx cancel.

func RunWithContext(ctx context.Context) error {
	if _, err := os.Stat(config.CaddyLogPath); err == nil {
		os.Remove(config.CaddyLogPath)
	}
	defer caddy.Stop()

	if err := config.GenerateCaddyConfig(); err != nil {
		logger.Error().Err(err).Msg("Failed to generate caddy config")
		return err
	}

	cfg, err := config.ReadCaddyConfig(config.CaddyActiveConfigPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read caddy config")
		return err
	}

	go func() {
		logger.Info().Msg("Starting caddy...")
		if err := caddy.Run(cfg); err != nil {
			logger.Error().Err(err).Msg("Caddy runtime error")
		} else {
			logger.Info().Msg("Caddy started!")
		}
	}()

	<-ctx.Done()

	logger.Info().Str("reason", ctx.Err().Error()).Msg("Caddy stopped")
	return ctx.Err()
}



