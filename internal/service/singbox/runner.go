package singbox

import (
	"context"
	"os"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/experimental"
	"github.com/sagernet/sing-box/experimental/clashapi"
	"github.com/sagernet/sing-box/experimental/v2rayapi"
	"github.com/sagernet/sing-box/include"
)

func init() {
	experimental.RegisterClashServerConstructor(clashapi.NewServer)
	experimental.RegisterV2RayServerConstructor(v2rayapi.NewServer)
}

// RunWithContext generates configuration, bootstraps sing-box core, and gracefully stops on ctx cancel.
func RunWithContext(ctx context.Context) error {
	if _, err := os.Stat(config.SingLogPath); err == nil {
		os.Remove(config.SingLogPath)
	}

	if err := config.GenerateSingConfig(); err != nil {
		logger.Error().Err(err).Msg("Failed to generate sing-box config")
		return err
	}

	singOptions, err := config.ReadSingConfig(config.SingActiveConfigPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read sing-box config")
		return err
	}

	singCtx := include.Context(context.Background())
	instance, err := box.New(box.Options{
		Context: singCtx,
		Options: singOptions,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create sing-box instance")
		return err
	}

	defer instance.Close()

	go func() {
		logger.Info().Msg("Starting sing-box...")

		if err = instance.Start(); err != nil {
			logger.Error().Err(err).Msg("Sing-box runtime error")
		} else {
			logger.Info().Msg("Sing-box started!")
		}
	}()

	<-ctx.Done()

	logger.Info().Str("reason", ctx.Err().Error()).Msg("Sing-box stopped")
	return ctx.Err()
}

