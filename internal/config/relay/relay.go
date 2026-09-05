package relay

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/geoip"
	"github.com/LalatinaHub/LatinaServer/internal/proxy"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	relays            []model.ProxyNode
	relaysMu          sync.RWMutex
	excludedRelayCode = []string{geoip.GetIpInfo().CountryCode}
	lastFetchTime     time.Time
	fetcherRunning    bool
	fetcherMu         sync.Mutex
)

const (
	relayTTL           = 15 * time.Minute
	refreshInterval    = 15 * time.Minute
	maxBackoffDuration = 1 * time.Minute
)

// GetRelays returns the current cached relay list (thread-safe).
func GetRelays() []model.ProxyNode {
	relaysMu.RLock()
	defer relaysMu.RUnlock()
	return relays
}

// setRelays updates the relay cache (thread-safe).
func setRelays(newRelays []model.ProxyNode) {
	relaysMu.Lock()
	defer relaysMu.Unlock()
	relays = newRelays
	lastFetchTime = time.Now()
}

// GatherRelays fetches relays once synchronously (for initial startup or manual trigger).
func GatherRelays() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.GetDB()
	if err != nil {
		logger.Error().Err(err).Msg("GatherRelays: error getting DB")
		return
	}
	proxyRepo := repository.NewProxyRepository(db)

	relaysData, err := proxyRepo.GetRelays(ctx, excludedRelayCode, 10)
	if err != nil {
		logger.Warn().Err(err).Msg("GatherRelays: error gathering relays (using stale data)")
		// Don't clear existing relays on failure, use stale data
		return
	}

	setRelays(relaysData)
	logger.Info().Int("count", len(relaysData)).Msg("Relay cache updated")
}

// StartBackgroundRelayFetcher starts a background goroutine that periodically refreshes relay data.
// It uses exponential backoff on failures and does not block the caller.
func StartBackgroundRelayFetcher(ctx context.Context) {
	fetcherMu.Lock()
	if fetcherRunning {
		fetcherMu.Unlock()
		return
	}
	fetcherRunning = true
	fetcherMu.Unlock()

	go func() {
		defer func() {
			fetcherMu.Lock()
			fetcherRunning = false
			fetcherMu.Unlock()
		}()

		timer := time.NewTimer(refreshInterval)
		defer timer.Stop()

		backoff := 1 * time.Second

		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("Relay fetcher stopped")
				return
			case <-timer.C:
				fetchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

				db, err := database.GetDB()
				var relaysData []model.ProxyNode
				if err == nil {
					proxyRepo := repository.NewProxyRepository(db)
					relaysData, err = proxyRepo.GetRelays(fetchCtx, excludedRelayCode, 10)
				}
				cancel()

				if err != nil {
					logger.Warn().
						Err(err).
						Dur("retry_in", backoff).
						Msg("Relay fetcher: error fetching relays, will retry soon")
					timer.Reset(backoff)
					backoff = min(backoff*2, maxBackoffDuration)
				} else {
					setRelays(relaysData)
					logger.Info().Int("count", len(relaysData)).Msg("Relay cache refreshed")
					backoff = 1 * time.Second // Reset backoff on success
					timer.Reset(refreshInterval)
				}
			}
		}
	}()
}

// GetRelayOutbounds converts cached relay nodes to sing-box outbound configurations.
func GetRelayOutbounds() []option.Outbound {
	var (
		proxies      = GetRelays()
		outbounds    = []option.Outbound{}
		outboundsMap = map[string][]option.Outbound{}
		converter    = proxy.NewConverterSingbox()
	)

	if len(proxies) == 0 {
		return outbounds
	}

	for _, proxyNode := range proxies {
		if len(outboundsMap[proxyNode.CountryCode]) < 5 {
			// Generate unique tag from remark or fallback to country code + index
			tag := proxyNode.Remark
			if tag == "" {
				tag = fmt.Sprintf("%s-relay-%d", proxyNode.CountryCode, len(outboundsMap[proxyNode.CountryCode]))
			}

			outbound, err := converter.ConvertToOutbound(proxyNode, tag)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("proxy", proxyNode.Remark).
					Msg("Error converting proxy to outbound")
				continue
			}

			outboundsMap[proxyNode.CountryCode] = append(outboundsMap[proxyNode.CountryCode], outbound)
		}
	}

	for cc, out := range outboundsMap {
		var outboundTags = []string{}
		for _, outbound := range out {
			if !slices.Contains(outboundTags, outbound.Tag) {
				outboundTags = append(outboundTags, outbound.Tag)
			}
		}

		urltest := option.Outbound{
			Tag:  cc,
			Type: C.TypeURLTest,
			Options: option.URLTestOutboundOptions{
				Outbounds: outboundTags,
			},
		}

		// Check tag
		if urltest.Tag != "" {
			outbounds = append(outbounds, urltest)
			outbounds = append(outbounds, out...)
		}
	}

	return outbounds
}


