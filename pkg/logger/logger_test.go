package logger_test

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/pkg/logger"
)

func TestLogger(t *testing.T) {
	// Setup development logger
	logger.SetupLogger("debug", false)
	if logger.Logger() == nil {
		t.Fatal("expected logger to not be nil")
	}

	logger.Debug().Msg("debug test")
	logger.Info().Msg("info test")
	logger.Warn().Msg("warn test")
	logger.Error().Msg("error test")

	// Context logger
	ctxLogger := logger.WithContext(map[string]interface{}{
		"request_id": "test-req-123",
		"trace_id":   "test-trace-456",
	})
	ctxLogger.Info().Msg("context test")

	// Setup production logger
	logger.SetupLogger("info", true)
	logger.Info().Msg("prod logger test")
}
