package config

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

var (
	localCtx   context.Context
	configPath = "/usr/local/etc/latinaserver/config.json"
)

func ReadSingConfig() option.Options {
	defer helper.CatchError(true)

	body, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var options option.Options
	err = options.UnmarshalJSONContext(localCtx, body)
	if err != nil {
		panic(err)
	}

	return options
}

func GenerateSingConfig() option.Options {
	premiumList := db.GetPremiumList()
	relayOutbounds := relay.GetRelayOutbounds()
	options := SingOptions

	for i, inbound := range options.Inbounds {
		var (
			port              = 52000 + i
			anyInboundOptions = inbound.Options
		)

		switch inbound.Type {
		case C.TypeMixed:
			anyInboundOptions = option.HTTPMixedInboundOptions{
				ListenOptions: option.ListenOptions{
					ListenPort: CS.MixedPort,
				},
			}
		case C.TypeTrojan:
			anyInboundOptions := option.TrojanInboundOptions{
				ListenOptions: option.ListenOptions{
					ListenPort: uint16(port),
				},
				Users:     []option.TrojanUser{},
				Multiplex: MultiplexOptions,
			}

			for _, user := range premiumList[C.TypeTrojan] {
				anyInboundOptions.Users = append(anyInboundOptions.Users, option.TrojanUser{
					Name:     strconv.Itoa(int(user.Id)),
					Password: user.Password,
				})
			}

			if anyInboundOptions.Transport != nil {
				switch anyInboundOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					anyInboundOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeHTTPUpgrade:
					anyInboundOptions.Transport.HTTPUpgradeOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeGRPC:
					anyInboundOptions.Transport.GRPCOptions.ServiceName = inbound.Type
				}
			}
		case C.TypeVMess:
			anyInboundOptions := option.VMessInboundOptions{
				ListenOptions: option.ListenOptions{
					ListenPort: uint16(port),
				},
				Users:     []option.VMessUser{},
				Multiplex: MultiplexOptions,
			}

			for _, user := range premiumList[C.TypeVMess] {
				anyInboundOptions.Users = append(anyInboundOptions.Users, option.VMessUser{
					Name: strconv.Itoa(int(user.Id)),
					UUID: user.Password,
				})
			}

			if anyInboundOptions.Transport != nil {
				switch anyInboundOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					anyInboundOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeHTTPUpgrade:
					anyInboundOptions.Transport.HTTPUpgradeOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeGRPC:
					anyInboundOptions.Transport.GRPCOptions.ServiceName = inbound.Type
				}
			}
		case C.TypeVLESS:
			anyInboundOptions := option.VLESSInboundOptions{
				ListenOptions: option.ListenOptions{
					ListenPort: uint16(port),
				},
				Users:     []option.VLESSUser{},
				Multiplex: MultiplexOptions,
			}

			for _, user := range premiumList[C.TypeVLESS] {
				anyInboundOptions.Users = append(anyInboundOptions.Users, option.VLESSUser{
					Name: strconv.Itoa(int(user.Id)),
					UUID: user.Password,
				})
			}

			if anyInboundOptions.Transport != nil {
				switch anyInboundOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					anyInboundOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeHTTPUpgrade:
					anyInboundOptions.Transport.HTTPUpgradeOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeGRPC:
					anyInboundOptions.Transport.GRPCOptions.ServiceName = inbound.Type
				}
			}

		}

		inbound.Options = anyInboundOptions
		options.Inbounds[i] = inbound
	}

	for _, list := range premiumList {
		for _, user := range list {
			options.Experimental.V2RayAPI.Stats.Users = append(options.Experimental.V2RayAPI.Stats.Users, strconv.Itoa(int(user.Id)))
		}
	}

	options.Outbounds = append(options.Outbounds, relayOutbounds...)

	// Eliminate existing rules if exists
	tempRules := []option.Rule{}
	for _, rule := range options.Route.Rules {
		switch rule.DefaultOptions.RouteOptions.Outbound {
		case C.TypeDNS, C.TypeDirect, C.TypeBlock:
			if rule.DefaultOptions.AuthUser != nil {
				continue
			}
			tempRules = append(tempRules, rule)
		}
	}
	options.Route.Rules = tempRules

	// Spesific route each server
	var (
		domains    = db.GetDomainList()
		serverInfo = helper.GetIpInfo()
		serverCode = strings.Split(serverInfo.Org, " ")[0]
	)
	switch serverCode {
	case "AS133800":
		for _, outbound := range options.Outbounds {
			if outbound.Tag == "SG" {
				options.Route.Rules = append(options.Route.Rules, []option.Rule{
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							RawDefaultRule: option.RawDefaultRule{
								Geosite: badoption.Listable[string]{"google", "rule-playstore", "rule-streaming"},
							},
							RuleAction: option.RuleAction{
								RouteOptions: option.RouteActionOptions{
									Outbound: "SG",
								},
							},
						},
					},
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							RawDefaultRule: option.RawDefaultRule{
								GeoIP: badoption.Listable[string]{"google"},
							},
							RuleAction: option.RuleAction{
								RouteOptions: option.RouteActionOptions{
									Outbound: "SG",
								},
							},
						},
					}}...)
				break
			}
		}
	}
	switch serverInfo.CountryCode {
	case "SG":
		var (
			assignedDomainsTag = []string{}
			proxyTag           = "ID-Socks-Proxy"
		)
		for _, domain := range domains {
			isExists := false
			for _, tag := range assignedDomainsTag {
				if tag == domain.Code {
					isExists = true
					break
				}
			}

			if isExists {
				continue
			}

			if domain.Location == "ID" {
				assignedDomainsTag = append(assignedDomainsTag, domain.Code)
				options.Outbounds = append(options.Outbounds, option.Outbound{
					Type: C.TypeSOCKS,
					Tag:  domain.Code,
					Options: option.SOCKSOutboundOptions{
						ServerOptions: option.ServerOptions{
							Server:     domain.Domain,
							ServerPort: CS.MixedPort,
						},
					},
				})
			}
		}

		options.Outbounds = append(options.Outbounds, option.Outbound{
			Type: C.TypeURLTest,
			Tag:  proxyTag,
			Options: option.URLTestOutboundOptions{
				Outbounds: assignedDomainsTag,
			},
		})

		options.Route.Rules = append(options.Route.Rules, []option.Rule{
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					RawDefaultRule: option.RawDefaultRule{
						Geosite: badoption.Listable[string]{"youtube"},
					},
					RuleAction: option.RuleAction{
						RouteOptions: option.RouteActionOptions{
							Outbound: proxyTag,
						},
					},
				},
			}}...)
	}

	// Adblock for specific user
	adblockRules := option.Rule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			RawDefaultRule: option.RawDefaultRule{
				Geosite:  badoption.Listable[string]{"rule-ads", "oisd-full"},
				AuthUser: badoption.Listable[string]{},
			},
			RuleAction: option.RuleAction{
				RouteOptions: option.RouteActionOptions{
					Outbound: C.TypeBlock,
				},
			},
		},
	}
	for _, premium := range premiumList {
		for _, user := range premium {
			if user.Adblock {
				adblockRules.DefaultOptions.AuthUser = append(adblockRules.DefaultOptions.AuthUser, strconv.Itoa(int(user.Id)))
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
					RawDefaultRule: option.RawDefaultRule{
						AuthUser: badoption.Listable[string]{},
						Network:  badoption.Listable[string]{"tcp"},
					},
					RuleAction: option.RuleAction{
						RouteOptions: option.RouteActionOptions{
							Outbound: outbound.Tag,
						},
					},
				},
			}

			for _, premium := range premiumList {
				for _, user := range premium {
					if user.CC == outbound.Tag {
						rule.DefaultOptions.AuthUser = append(rule.DefaultOptions.AuthUser, strconv.Itoa(int(user.Id)))
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
	f, err := os.Create(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(options, "", "\t")
	f.WriteString(string(b))

	if err != nil {
		panic(err)
	}

	return options
}
