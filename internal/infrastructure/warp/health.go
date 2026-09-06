package warp

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/logger"
)

var (
	isHealthy      atomic.Bool
	running        atomic.Bool
	statusCallback func(healthy bool)
	callbackMu     sync.RWMutex
)

func init() {
	isHealthy.Store(true) // Start optimistic
}

// SetStatusChangeCallback registers a callback invoked whenever WARP health status transitions.
func SetStatusChangeCallback(fn func(healthy bool)) {
	callbackMu.Lock()
	defer callbackMu.Unlock()
	statusCallback = fn
}

func notifyStatusChange(healthy bool) {
	callbackMu.RLock()
	fn := statusCallback
	callbackMu.RUnlock()
	if fn != nil {
		go fn(healthy)
	}
}

// IsHealthy returns the current health status of Cloudflare WARP.
func IsHealthy() bool {
	return isHealthy.Load()
}

// SetHealthy manually sets the health status (useful for tests or overrides).
func SetHealthy(healthy bool) {
	old := isHealthy.Swap(healthy)
	if old != healthy {
		notifyStatusChange(healthy)
	}
}

// CheckConnectivity verifies network connectivity to the WARP endpoint.
func CheckConnectivity(ctx context.Context, endpoint string) error {
	if endpoint == "" {
		endpoint = DefaultWarpEndpoint
	}

	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.DialContext(ctx, "udp", endpoint)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// StartHealthChecker starts a background goroutine that periodically monitors WARP connectivity.
func StartHealthChecker(ctx context.Context, interval time.Duration) {
	if !running.CompareAndSwap(false, true) {
		return
	}

	if interval <= 0 {
		interval = 30 * time.Second
	}

	go func() {
		defer running.Store(false)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var consecutiveFailures int

		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("WARP health checker stopped")
				return
			case <-ticker.C:
				probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				cfg, err := GetOrInitWarpConfig(probeCtx)
				endpoint := DefaultWarpEndpoint
				if err == nil && cfg != nil && cfg.Endpoint != "" {
					endpoint = cfg.Endpoint
				}

				checkErr := CheckConnectivity(probeCtx, endpoint)
				cancel()

				if checkErr != nil {
					consecutiveFailures++
					if consecutiveFailures >= 3 && isHealthy.Load() {
						SetHealthy(false)
						logger.Warn().
							Err(checkErr).
							Int("failures", consecutiveFailures).
							Msg("WARP endpoint unreachable; routing fallback to direct activated")
					}
				} else {
					if !isHealthy.Load() {
						SetHealthy(true)
						logger.Info().
							Int("cleared_failures", consecutiveFailures).
							Msg("WARP connectivity restored; smart routing reactivated")
					}
					consecutiveFailures = 0
				}
			}
		}
	}()
}
