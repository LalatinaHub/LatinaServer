package systemctl

import (
	"os/exec"

	"github.com/LalatinaHub/LatinaServer/pkg/logger"
)

func Reload(names ...string) {
	for _, name := range names {
		logger.Info().Str("service", name).Msg("Reloading system service...")
		_, err := exec.Command("systemctl", "reload", name).Output()
		if err != nil {
			if err.Error() == "exit status 1" {
				logger.Warn().Str("service", name).Msg("Reload failed, attempting restart...")
				if err := exec.Command("systemctl", "restart", name).Run(); err != nil {
					logger.Error().Err(err).Str("service", name).Msg("Failed to restart system service")
				}
			} else {
				logger.Error().Err(err).Str("service", name).Msg("Failed to reload system service")
			}
			continue
		}

		logger.Info().Str("service", name).Msg("System service successfully reloaded")
	}
}

func CatchError(print bool) any {
	message := recover()

	if message != nil && print {
		logger.Error().Interface("recovered", message).Msg("Recovered from panic")
	}
	return message
}

