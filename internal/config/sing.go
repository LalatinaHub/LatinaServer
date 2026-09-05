package config

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/LalatinaHub/LatinaServer/internal/config/relay"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

// ReadSingConfig loads and unmarshals sing-box JSON configuration using registries.

func ReadSingConfig(configLocation string) (option.Options, error) {
	defer systemctl.CatchError(true)

	body, err := os.ReadFile(configLocation)
	if err != nil {
		return option.Options{}, appErrors.NewConfigError("failed to read sing config", err)
	}

	var (
		ctx     context.Context = context.Background()
		options option.Options
	)

	ctx = include.Context(ctx)
	err = options.UnmarshalJSONContext(ctx, body)
	if err != nil {
		return option.Options{}, appErrors.NewConfigError("failed to unmarshal sing config", err)
	}

	return options, nil
}
// GenerateSingConfig compiles active user credentials and dynamic relay outbounds into sing-box config.


func GenerateSingConfig() error {
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

	var relayOutbounds = relay.GetRelayOutbounds()


	for i, inbound := range options.Inbounds {
		if strings.HasSuffix(inbound.Tag, "-udp") {
			continue
		}

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

	for _, list := range premiumList {
		for _, user := range list {
			options.Experimental.V2RayAPI.Stats.Users = append(options.Experimental.V2RayAPI.Stats.Users, strconv.Itoa(int(user.ID)))
		}
	}

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

	// List inbounds and outbounds tag to v2ray_api field
	// Sync shadowsocks outbound password with inbound
	ssPass := GetSSPassword()
	for i, outbound := range options.Outbounds {
		if outbound.Tag == "ss-out" && outbound.Type == C.TypeShadowsocks {
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
		logger.Info().Msg("Sing-box config updated (hash changed)")
	} else {
		logger.Debug().Msg("Sing-box config unchanged (hash match), skipped write")
	}

	return nil
}
