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

	singCtx := config.SingContext(context.Background())
	instance, err := box.New(box.Options{
		Context: singCtx,
		Options: singOptions,
	})
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Failed to create sing-box instance; attempting Tier-1 fallback without relays...")

		fallbackErr := config.GenerateSingConfigWithOptions(false)
		if fallbackErr == nil {
			if fallbackOptions, readErr := config.ReadSingConfig(config.SingActiveConfigPath); readErr == nil {
				instance, err = box.New(box.Options{
					Context: singCtx,
					Options: fallbackOptions,
				})
			}
		}

		// Tier-2 fallback: If still failing, fallback to pure direct mode (no relays, no WARP)
		if err != nil {
			logger.Warn().
				Err(err).
				Msg("Tier-1 fallback failed; activating Tier-2 safe mode (direct routing only)...")

			if safeErr := config.GenerateSingConfigWithAllOptions(false, false); safeErr != nil {
				logger.Error().Err(safeErr).Msg("Failed to generate safe mode sing-box config")
				return err
			}

			safeOptions, readErr := config.ReadSingConfig(config.SingActiveConfigPath)
			if readErr != nil {
				logger.Error().Err(readErr).Msg("Failed to read safe mode sing-box config")
				return err
			}

			instance, err = box.New(box.Options{
				Context: singCtx,
				Options: safeOptions,
			})
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create safe mode sing-box instance")
				return err
			}
		}

		logger.Warn().Msg("Sing-box self-healing recovery succeeded (active in safe fallback mode)")
	}

	defer instance.Close()

	startErrChan := make(chan error, 1)
	go func() {
		logger.Info().Msg("Starting sing-box...")

		if err = instance.Start(); err != nil {
			logger.Error().Err(err).Msg("Sing-box runtime error")
			startErrChan <- err
		} else {
			logger.Info().Msg("Sing-box started!")
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info().Str("reason", ctx.Err().Error()).Msg("Sing-box stopped")
		return ctx.Err()
	case err := <-startErrChan:
		return err
	}
}

