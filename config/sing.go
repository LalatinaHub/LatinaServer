package config

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

func ReadSingConfig(configLocation string) option.Options {
	defer helper.CatchError(true)

	body, err := os.ReadFile(configLocation)
	if err != nil {
		panic(err)
	}

	var (
		ctx     context.Context = context.Background()
		options option.Options
	)

	ctx = box.Context(ctx, include.InboundRegistry(), include.OutboundRegistry(), include.EndpointRegistry(), include.DNSTransportRegistry())
	err = options.UnmarshalJSONContext(ctx, body)
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
	for _, inbound := range options.Inbounds {
		options.Experimental.V2RayAPI.Stats.Inbounds = append(options.Experimental.V2RayAPI.Stats.Inbounds, inbound.Tag)
	}

	for _, outbound := range options.Outbounds {
		options.Experimental.V2RayAPI.Stats.Outbounds = append(options.Experimental.V2RayAPI.Stats.Outbounds, outbound.Tag)
	}

	// Write new config
	SaveJsonToFile(CS.SING_ACTIVE_CONFIG_PATH, options)
}
