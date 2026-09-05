package main

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	runtimeDebug "runtime/debug"
	"syscall"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/internal/config/relay"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/v2ray"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/internal/service/caddy"
	"github.com/LalatinaHub/LatinaServer/internal/service/singbox"
	"github.com/LalatinaHub/LatinaServer/internal/service/web"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/go-co-op/gocron"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

func updateUsersQuota() {
	defer systemctl.CatchError(true)

	db, err := database.GetDB()
	if err != nil {
		logger.Error().Err(err).Msg("Error getting DB connection for quota update")
		return
	}

	userRepo := repository.NewUserRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	singConfig, err := config.ReadSingConfig(config.SingActiveConfigPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read sing-box config for quota update")
		return
	}

	if singConfig.Experimental == nil || singConfig.Experimental.V2RayAPI == nil {
		return
	}

	usages := v2ray.GetUsersStatsBatch(singConfig.Experimental.V2RayAPI.Stats.Users)
	if len(usages) == 0 {
		return
	}

	depletedUsers, err := userRepo.DeductQuotaBatch(ctx, usages)
	if err != nil {
		logger.Error().Err(err).Msg("Error batch updating user quota")
		return
	}

	if len(depletedUsers) > 0 {
		logger.Warn().
			Int("depleted_count", len(depletedUsers)).
			Msg("Depleted quota detected, reloading service...")
		systemctl.Reload(config.ServiceLatinaServer)
	}
}

func main() {
	// Initialize structured logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	isProd := os.Getenv("ENV") == "production"
	logger.SetupLogger(logLevel, isProd)

	logger.Info().
		Str("log_level", logLevel).
		Bool("production", isProd).
		Msg("Starting LatinaServer...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database indexes for performance
	if err := database.InitIndexes(ctx); err != nil {
		logger.Warn().Err(err).Msg("Failed to initialize database indexes")
	} else {
		logger.Info().Msg("Database indexes initialized successfully")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	s := gocron.NewScheduler(loc)
	s.Every(5).Minutes().Tag("update-quota").Do(updateUsersQuota)
	s.Every(30).Minutes().Tag("memory-cleanup").Do(func() {
		runtimeDebug.FreeOSMemory()
		logger.Debug().Msg("Scheduled memory cleanup performed")
	})
	s.Every(1).Day().At("03:00").Tag("reboot").Do(func() {
		if err := exec.Command("reboot").Run(); err != nil {
			logger.Error().Err(err).Msg("Failed to reboot the server")
		}
	})

	// Start async functions
	s.StartAsync()

	// Initial relay fetch (non-blocking)
	go relay.GatherRelays()
	
	// Start background relay fetcher with 15-minute refresh interval
	relay.StartBackgroundRelayFetcher(ctx)

	for {
		cancelCtx, cancelServices := context.WithCancel(ctx)

		// Pre-flight generate Caddy & Sing-box configs concurrently
		if err := config.GenerateConfigsParallel(cancelCtx); err != nil {
			logger.Error().Err(err).Msg("Failed to generate configurations concurrently")
		}

		go caddy.RunWithContext(cancelCtx)
		go singbox.RunWithContext(cancelCtx)
		go web.RunWithContext(cancelCtx)

		for {
			osSignal := <-c
			logger.Info().Str("signal", osSignal.String()).Msg("Stopping services...")
			time.Sleep(2 * time.Second)

			if osSignal == syscall.SIGHUP {
				cancelServices()
				break
			}

			// Graceful shutdown
			cancelServices()
			cancel() // Stop relay fetcher and other background tasks
			_ = v2ray.Close()
			_ = database.Close()
			logger.Info().Msg("LatinaServer shut down successfully")
			return
		}
	}
}


