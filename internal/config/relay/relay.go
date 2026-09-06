package relay

import (
	"context"
	"fmt"
	"slices"
	"strings"
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

var reservedOutboundTags = map[string]bool{
	"direct":    true,
	"ss-out":    true,
	"final-dns": true,
	"block":     true,
	"dns-out":   true,
}

// GetRelayOutbounds converts cached relay nodes to sing-box outbound configurations.
func GetRelayOutbounds() []option.Outbound {
	var (
		proxies      = GetRelays()
		outbounds    = []option.Outbound{}
		outboundsMap = map[string][]option.Outbound{}
		converter    = proxy.NewConverterSingbox()
		usedTags     = make(map[string]bool)
	)

	if len(proxies) == 0 {
		return outbounds
	}

	for tag := range reservedOutboundTags {
		usedTags[tag] = true
	}

	for _, proxyNode := range proxies {
		cc := strings.ToUpper(strings.TrimSpace(proxyNode.CountryCode))
		if cc == "" {
			cc = "OTHER"
		}

		if len(outboundsMap[cc]) < 5 {
			baseTag := strings.TrimSpace(proxyNode.Remark)
			if baseTag == "" {
				baseTag = fmt.Sprintf("%s-relay-%d", cc, len(outboundsMap[cc])+1)
			}

			tag := baseTag
			suffix := 1
			for usedTags[tag] {
				suffix++
				tag = fmt.Sprintf("%s-%d", baseTag, suffix)
			}

			outbound, err := converter.ConvertToOutbound(proxyNode, tag)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("proxy", proxyNode.Remark).
					Msg("Error converting proxy to outbound (skipping)")
				continue
			}

			usedTags[tag] = true
			outboundsMap[cc] = append(outboundsMap[cc], outbound)
		}
	}

	for cc, out := range outboundsMap {
		var outboundTags = []string{}
		for _, outbound := range out {
			if !slices.Contains(outboundTags, outbound.Tag) {
				outboundTags = append(outboundTags, outbound.Tag)
			}
		}

		if len(outboundTags) == 0 {
			continue
		}

		urlTestTag := cc
		suffix := 1
		for usedTags[urlTestTag] {
			suffix++
			urlTestTag = fmt.Sprintf("%s-group-%d", cc, suffix)
		}
		usedTags[urlTestTag] = true

		urltest := option.Outbound{
			Tag:  urlTestTag,
			Type: C.TypeURLTest,
			Options: option.URLTestOutboundOptions{
				Outbounds: outboundTags,
			},
		}

		outbounds = append(outbounds, urltest)
		outbounds = append(outbounds, out...)
	}

	return outbounds
}


