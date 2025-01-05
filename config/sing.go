package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func ReadSingConfig(configLocation string) option.Options {
	defer helper.CatchError(true)

	body, err := os.ReadFile(configLocation)
	if err != nil {
		panic(err)
	}

	var options option.Options
	err = options.UnmarshalJSON(body)
	if err != nil {
		panic(err)
	}

	return options
}

func GenerateSingConfig() {
	var (
		premiumList    = db.GetPremiumList()
		relayOutbounds = relay.GetRelayOutbounds()
		options        = ReadSingConfig(CS.SING_CONFIG_PATH)
	)

	for i, inbound := range options.Inbounds {
		if strings.HasSuffix(inbound.Tag, "-udp") {
			continue
		}

		switch inbound.Type {
		case C.TypeTrojan:
			inbound.TrojanOptions.Users = []option.TrojanUser{}
			for _, user := range premiumList[C.TypeTrojan] {
				inbound.TrojanOptions.Users = append(inbound.TrojanOptions.Users, option.TrojanUser{
					Name:     strconv.Itoa(int(user.ID)),
					Password: user.Password,
				})
			}
		case C.TypeVMess:
			inbound.VMessOptions.Users = []option.VMessUser{}
			for _, user := range premiumList[C.TypeVMess] {
				inbound.VMessOptions.Users = append(inbound.VMessOptions.Users, option.VMessUser{
					Name: strconv.Itoa(int(user.ID)),
					UUID: user.Password,
				})
			}
		case C.TypeVLESS:
			inbound.VLESSOptions.Users = []option.VLESSUser{}
			for _, user := range premiumList[C.TypeVLESS] {
				inbound.VLESSOptions.Users = append(inbound.VLESSOptions.Users, option.VLESSUser{
					Name: strconv.Itoa(int(user.ID)),
					UUID: user.Password,
				})
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

	// Adblock for specific user
	adblockRules := option.Rule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			Geosite:  option.Listable[string]{"rule-ads", "oisd-full"},
			AuthUser: option.Listable[string]{},
			Outbound: C.TypeBlock,
		},
	}
	for _, premium := range premiumList {
		for _, user := range premium {
			if user.Adblock {
				adblockRules.DefaultOptions.AuthUser = append(adblockRules.DefaultOptions.AuthUser, strconv.Itoa(int(user.ID)))
			}
		}
	}
	if len(adblockRules.DefaultOptions.AuthUser) > 0 {
		options.Route.Rules = append(options.Route.Rules, adblockRules)
	}

	// Relay for specific user
	for _, outbound := range relayOutbounds {
		if len(outbound.Tag) < 5 {
			rule := option.Rule{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					AuthUser: option.Listable[string]{},
					Network:  option.Listable[string]{"tcp"},
					Outbound: outbound.Tag,
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
	for _, inbound := range options.Inbounds {
		options.Experimental.V2RayAPI.Stats.Inbounds = append(options.Experimental.V2RayAPI.Stats.Inbounds, inbound.Tag)
	}

	for _, outbound := range options.Outbounds {
		options.Experimental.V2RayAPI.Stats.Outbounds = append(options.Experimental.V2RayAPI.Stats.Outbounds, outbound.Tag)
	}

	// Write new config
	SaveJsonToFile(CS.SING_ACTIVE_CONFIG_PATH, options)
}
