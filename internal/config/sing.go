package config

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LalatinaHub/LatinaServer/internal/config/relay"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/adblock"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/warp"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/LalatinaHub/LatinaServer/resources"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/protocol/wireguard"
	"github.com/sagernet/sing/common/json/badoption"
)

// SingContext initializes sing-box registries with WireGuard endpoint support.
func SingContext(ctx context.Context) context.Context {
	epReg := include.EndpointRegistry()
	wireguard.RegisterEndpoint(epReg)
	return box.Context(
		ctx,
		include.InboundRegistry(),
		include.OutboundRegistry(),
		epReg,
		include.DNSTransportRegistry(),
		include.ServiceRegistry(),
		include.CertificateProviderRegistry(),
	)
}

// ReadSingConfig loads and unmarshals sing-box JSON configuration using registries.
func ReadSingConfig(configLocation string) (option.Options, error) {
	defer systemctl.CatchError(true)

	body, err := os.ReadFile(configLocation)
	if err != nil {
		if os.IsNotExist(err) && configLocation == SingConfigPath && len(resources.DefaultSingBoxTemplate) > 0 {
			logger.Info().Str("path", SingConfigPath).Msg("Sing-box config template not found on disk; seeding from embedded default")
			body = resources.DefaultSingBoxTemplate
			_ = os.MkdirAll(filepath.Dir(SingConfigPath), 0755)
			_ = os.WriteFile(SingConfigPath, body, 0644)
		} else {
			return option.Options{}, appErrors.NewConfigError("failed to read sing config", err)
		}
	}

	var (
		ctx     context.Context = context.Background()
		options option.Options
	)

	ctx = SingContext(ctx)
	err = options.UnmarshalJSONContext(ctx, body)
	if err != nil {
		return option.Options{}, appErrors.NewConfigError("failed to unmarshal sing config", err)
	}

	return options, nil
}
// GenerateSingConfig compiles active user credentials, dynamic relays, and WARP into sing-box config.
func GenerateSingConfig() error {
	return GenerateSingConfigWithAllOptions(true, true)
}

// GenerateSingConfigWithOptions compiles active user credentials and optionally dynamic relay outbounds into sing-box config.
func GenerateSingConfigWithOptions(includeRelays bool) error {
	return GenerateSingConfigWithAllOptions(includeRelays, true)
}

// GenerateSingConfigWithAllOptions compiles active user credentials, optional dynamic relays, and optional WARP unlocker into sing-box config.
func GenerateSingConfigWithAllOptions(includeRelays bool, includeWarp bool) error {
	db, err := database.GetDB()
	if err != nil {
		return appErrors.NewDatabaseError("failed to get database connection", err)
	}

	userRepo := repository.NewUserRepository(db)
	premiumList, err := userRepo.GetActiveUsersGroupedByVPN(context.Background())
	if err != nil {
		return appErrors.NewDatabaseError("failed to get active users grouped by VPN", err)
	}

	options, err := ReadSingConfig(SingConfigPath)
	if err != nil {
		return err
	}

	// Configure sing-box logging to display complete connection logs in terminal
	if options.Log == nil {
		options.Log = &option.LogOptions{}
	}
	if logOutput := os.Getenv("SINGBOX_LOG_OUTPUT"); logOutput != "" {
		options.Log.Output = logOutput
	} else if os.Getenv("LOG_TERMINAL") != "false" {
		// Empty string outputs directly to console / terminal (stderr)
		options.Log.Output = ""
	}
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		options.Log.Level = strings.ToLower(envLevel)
	} else if options.Log.Level == "" {
		options.Log.Level = "info"
	}
	options.Log.Timestamp = true

	var relayOutbounds []option.Outbound
	if includeRelays {
		relayOutbounds = relay.GetRelayOutbounds()
	}

	for i, inbound := range options.Inbounds {
		switch inbound.Type {
		case C.TypeTrojan:
			var trojanOptions = inbound.Options.(*option.TrojanInboundOptions)
			trojanOptions.Users = []option.TrojanUser{}
			for _, user := range premiumList[C.TypeTrojan] {
				trojanOptions.Users = append(trojanOptions.Users, option.TrojanUser{
					Name:     strconv.Itoa(int(user.ID)),
					Password: user.Password,
				})
			}
			inbound.Options = trojanOptions
		case C.TypeVMess:
			var vmessOptions = inbound.Options.(*option.VMessInboundOptions)
			vmessOptions.Users = []option.VMessUser{}
			for _, user := range premiumList[C.TypeVMess] {
				vmessOptions.Users = append(vmessOptions.Users, option.VMessUser{
					Name: strconv.Itoa(int(user.ID)),
					UUID: user.Password,
				})
			}
			inbound.Options = vmessOptions
		case C.TypeVLESS:
			var vlessOptions = inbound.Options.(*option.VLESSInboundOptions)
			vlessOptions.Users = []option.VLESSUser{}
			for _, user := range premiumList[C.TypeVLESS] {
				vlessOptions.Users = append(vlessOptions.Users, option.VLESSUser{
					Name: strconv.Itoa(int(user.ID)),
					UUID: user.Password,
				})
			}
			inbound.Options = vlessOptions
		case C.TypeShadowsocks:
			if ssOptions, ok := inbound.Options.(*option.ShadowsocksInboundOptions); ok {
				ssOptions.Password = GetSSPassword()
				inbound.Options = ssOptions
			}
		}

		options.Inbounds[i] = inbound
	}

	var adblockUsers []string
	for _, list := range premiumList {
		for _, user := range list {
			options.Experimental.V2RayAPI.Stats.Users = append(options.Experimental.V2RayAPI.Stats.Users, strconv.Itoa(int(user.ID)))
			if user.Adblock {
				adblockUsers = append(adblockUsers, strconv.Itoa(int(user.ID)))
			}
		}
	}

	// AdBlock & Malware Shield DNS integration
	if options.DNS != nil {
		// Ensure AdGuard DNS server is registered in DNS servers
		hasAdGuardDNS := false
		for _, server := range options.DNS.Servers {
			if server.Tag == adblock.AdGuardDNSTag {
				hasAdGuardDNS = true
				break
			}
		}
		if !hasAdGuardDNS {
			options.DNS.Servers = append(options.DNS.Servers, adblock.BuildAdGuardDNSServer())
		}

		// Inject or update AdGuard DNS rule for users with Adblock enabled
		if len(adblockUsers) > 0 {
			dnsRule := adblock.BuildAdblockDNSRule(adblockUsers)
			ruleIndex := -1
			for idx, r := range options.DNS.Rules {
				if r.DefaultOptions.RouteOptions.Server == adblock.AdGuardDNSTag {
					ruleIndex = idx
					break
				}
			}
			if ruleIndex >= 0 {
				options.DNS.Rules[ruleIndex] = dnsRule
			} else {
				options.DNS.Rules = append([]option.DNSRule{dnsRule}, options.DNS.Rules...)
			}
		} else {
			var filteredDNSRules []option.DNSRule
			for _, r := range options.DNS.Rules {
				if r.DefaultOptions.RouteOptions.Server != adblock.AdGuardDNSTag {
					filteredDNSRules = append(filteredDNSRules, r)
				}
			}
			options.DNS.Rules = filteredDNSRules
		}
	}

	if options.Route == nil {
		options.Route = &option.RouteOptions{}
	}

	// AdBlock & Malware Shield Routing integration: reject ads and malware for adblock users
	if len(adblockUsers) > 0 {
		adblockRouteRule := adblock.BuildAdblockRouteRule(adblockUsers)
		ruleIndex := -1
		for idx, r := range options.Route.Rules {
			if r.DefaultOptions.RuleAction.Action == C.RuleActionTypeReject && len(r.DefaultOptions.RawDefaultRule.AuthUser) > 0 {
				ruleIndex = idx
				break
			}
		}
		if ruleIndex >= 0 {
			options.Route.Rules[ruleIndex] = adblockRouteRule
		} else {
			insertIdx := 0
			for rIdx, r := range options.Route.Rules {
				if r.DefaultOptions.RuleAction.Action == "sniff" ||
					r.DefaultOptions.RuleAction.Action == "hijack-dns" ||
					(r.DefaultOptions.RuleAction.Action == "reject" && len(r.DefaultOptions.RawDefaultRule.AuthUser) == 0) {
					insertIdx = rIdx + 1
				}
			}
			if insertIdx >= len(options.Route.Rules) {
				options.Route.Rules = append(options.Route.Rules, adblockRouteRule)
			} else {
				options.Route.Rules = append(options.Route.Rules[:insertIdx], append([]option.Rule{adblockRouteRule}, options.Route.Rules[insertIdx:]...)...)
			}
		}
	} else {
		var filteredRules []option.Rule
		for _, r := range options.Route.Rules {
			if !(r.DefaultOptions.RuleAction.Action == C.RuleActionTypeReject && len(r.DefaultOptions.RawDefaultRule.AuthUser) > 0) {
				filteredRules = append(filteredRules, r)
			}
		}
		options.Route.Rules = filteredRules
	}

	// Cloudflare WARP smart routing for AI & streaming unlocker
	if includeWarp && os.Getenv("WARP_ENABLED") != "false" {
		warpCfg, err := warp.GetOrInitWarpConfig(context.Background())
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to initialize WARP configuration; proceeding without WARP")
		} else {
			ep, epErr := warp.BuildEndpoint(warpCfg)
			if epErr != nil {
				logger.Warn().Err(epErr).Msg("Failed to build WireGuard endpoint for WARP")
			} else {
				// Ensure warp-out endpoint is in options.Endpoints
				warpIndex := -1
				for idx, existingEp := range options.Endpoints {
					if existingEp.Tag == warp.EndpointTag {
						warpIndex = idx
						break
					}
				}
				if warpIndex >= 0 {
					options.Endpoints[warpIndex] = *ep
				} else {
					options.Endpoints = append(options.Endpoints, *ep)
				}

				// Determine target outbound based on WARP health status
				targetOutbound := warp.EndpointTag
				if !warp.IsHealthy() {
					targetOutbound = "direct"
					logger.Warn().Msg("WARP is currently unhealthy; routing unlocker domains to direct")
				}

				// Update existing unlocker rule or insert new one
				unlockerDomains := warp.GetAllUnlockerDomains()
				ruleUpdated := false
				for rIdx, r := range options.Route.Rules {
					if r.DefaultOptions.RuleAction.RouteOptions.Outbound == warp.EndpointTag ||
						(len(r.DefaultOptions.RawDefaultRule.DomainSuffix) > 0 && len(unlockerDomains) > 0 && r.DefaultOptions.RawDefaultRule.DomainSuffix[0] == unlockerDomains[0]) {
						options.Route.Rules[rIdx].DefaultOptions.RuleAction.RouteOptions.Outbound = targetOutbound
						ruleUpdated = true
						break
					}
				}
				if !ruleUpdated {
					smartRule := warp.BuildSmartRule(targetOutbound)
					options.Route.Rules = append(options.Route.Rules, smartRule)
				}
			}
		}
	} else {
		// WARP explicitly disabled: remove warp endpoint and fallback any warp rules to direct
		var filteredEndpoints []option.Endpoint
		for _, ep := range options.Endpoints {
			if ep.Tag != warp.EndpointTag {
				filteredEndpoints = append(filteredEndpoints, ep)
			}
		}
		options.Endpoints = filteredEndpoints

		for rIdx, r := range options.Route.Rules {
			if r.DefaultOptions.RuleAction.RouteOptions.Outbound == warp.EndpointTag {
				options.Route.Rules[rIdx].DefaultOptions.RuleAction.RouteOptions.Outbound = "direct"
			}
		}
	}

	if includeRelays && len(relayOutbounds) > 0 {
		options.Outbounds = append(options.Outbounds, relayOutbounds...)

		// Relay for specific user
		for _, outbound := range relayOutbounds {
			if len(outbound.Tag) < 5 {
				rule := option.Rule{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						RawDefaultRule: option.RawDefaultRule{
							AuthUser: badoption.Listable[string]{},
							Network:  badoption.Listable[string]{"tcp"},
						},
						RuleAction: option.RuleAction{
							Action: "route",
							RouteOptions: option.RouteActionOptions{
								Outbound: outbound.Tag,
							},
						},
					},
				}

				for _, premium := range premiumList {
					for _, user := range premium {
						if user.Relay == outbound.Tag {
							rule.DefaultOptions.AuthUser = append(rule.DefaultOptions.AuthUser, strconv.Itoa(int(user.ID)))
						}
					}
				}

				if len(rule.DefaultOptions.AuthUser) > 0 {
					options.Route.Rules = append(options.Route.Rules, rule)
				}
			}
		}
	}

	// List inbounds and outbounds tag to v2ray_api field
	// Sync shadowsocks outbound password with inbound
	ssPass := GetSSPassword()
	for i, outbound := range options.Outbounds {
		if (outbound.Tag == "ss-out" || outbound.Tag == "ss-out-brutal") && outbound.Type == C.TypeShadowsocks {
			if ssOutOptions, ok := outbound.Options.(*option.ShadowsocksOutboundOptions); ok {
				ssOutOptions.Password = ssPass
				options.Outbounds[i] = outbound
			}
		}
	}

	// Update Clash API secret if configured
	if clashSecret := GetClashSecret(); clashSecret != "" && options.Experimental != nil && options.Experimental.ClashAPI != nil {
		options.Experimental.ClashAPI.Secret = clashSecret
	}

	for _, inbound := range options.Inbounds {
		options.Experimental.V2RayAPI.Stats.Inbounds = append(options.Experimental.V2RayAPI.Stats.Inbounds, inbound.Tag)
	}

	for _, outbound := range options.Outbounds {
		options.Experimental.V2RayAPI.Stats.Outbounds = append(options.Experimental.V2RayAPI.Stats.Outbounds, outbound.Tag)
	}

	// Write new config
	changed, err := SaveJsonToFileWithCache(SingActiveConfigPath, options)
	if err != nil {
		return appErrors.NewConfigError("failed to save sing config", err)
	}
	if changed {
		logger.Info().Bool("include_relays", includeRelays).Bool("include_warp", includeWarp).Msg("Sing-box config updated (hash changed)")
	} else {
		logger.Debug().Bool("include_relays", includeRelays).Bool("include_warp", includeWarp).Msg("Sing-box config unchanged (hash match), skipped write")
	}

	return nil
}
