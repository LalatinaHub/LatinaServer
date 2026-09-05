package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func RunWithContext(ctx context.Context) error {
	srv := NewServerWithContext(ctx)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.WebServerPort),
		Handler:      h2c.NewHandler(srv.SetupRouter(), &http2.Server{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info().Int("port", config.WebServerPort).Msg("Starting gin web server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("Gin web server error")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Warn().Err(err).Msg("Gin server shutdown forced")
		_ = server.Close()
	}

	logger.Info().Str("reason", ctx.Err().Error()).Msg("Gin stopped")
	return ctx.Err()
}

