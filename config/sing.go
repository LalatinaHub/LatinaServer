package config

import (
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
)

var (
	RealityPublicKey  = "dSurRwxcBfR-kZGO6UEb8EeweJjE4HyVKpUJOGZSXQs"
	RealityPrivateKey = "GHTprpUhfzbhJrtcAPrDKFJt6URah5VJN-39jFOOmVI"
	RealityShortID    = option.Listable[string]{"193ad0acc0a872d8"}
)

func ReadSingConfig() option.Options {
	body, err := os.ReadFile("/usr/local/etc/latinaserver/config.json")
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

func WriteSingConfig() option.Options {
	premiumList := db.GetPremiumList()
	sniList := db.GetSniList()
	relayOutbounds := relay.GetRelayOutbounds()
	options := ReadSingConfig()
	options.Experimental = &option.ExperimentalOptions{
		ClashAPI: &option.ClashAPIOptions{
			ExternalController: "0.0.0.0:9090",
			ExternalUI:         "/usr/local/latinaserver/dashboard/",
			Secret:             os.Getenv("PASSWORD"),
		},
		V2RayAPI: &option.V2RayAPIOptions{
			Listen: CS.V2rayAPIAddress,
			Stats: &option.V2RayStatsServiceOptions{
				Enabled:   true,
				Inbounds:  []string{},
				Outbounds: []string{},
				Users:     []string{},
			},
		},
	}

	var inbounds []option.Inbound
	for _, inbound := range options.Inbounds {
		var (
			port = 52000 + len(inbounds)
		)

		if strings.Contains(inbound.Tag, "reality") {
			continue
		}

		switch inbound.Type {
		case C.TypeTrojan:
			inbound.TrojanOptions.ListenPort = uint16(port)
			inbound.TrojanOptions.Users = []option.TrojanUser{}

			for _, user := range premiumList[C.TypeTrojan] {
				inbound.TrojanOptions.Users = append(inbound.TrojanOptions.Users, option.TrojanUser{
					Name:     strconv.Itoa(int(user.Id)),
					Password: user.Password,
				})
			}

			if inbound.TrojanOptions.Transport != nil {
				switch inbound.TrojanOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					inbound.TrojanOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeHTTPUpgrade:
					inbound.TrojanOptions.Transport.HTTPUpgradeOptions.Path = "/" + inbound.Type
				}
			}
		case C.TypeVMess:
			inbound.VMessOptions.ListenPort = uint16(port)
			inbound.VMessOptions.Users = []option.VMessUser{}

			for _, user := range premiumList[C.TypeVMess] {
				inbound.VMessOptions.Users = append(inbound.VMessOptions.Users, option.VMessUser{
					Name: strconv.Itoa(int(user.Id)),
					UUID: user.Password,
				})
			}

			if inbound.VMessOptions.Transport != nil {
				switch inbound.VMessOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					inbound.VMessOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeHTTPUpgrade:
					inbound.VMessOptions.Transport.HTTPUpgradeOptions.Path = "/" + inbound.Type
				}
			}
		case C.TypeVLESS:
			inbound.VLESSOptions.ListenPort = uint16(port)
			inbound.VLESSOptions.Users = []option.VLESSUser{}

			for _, user := range premiumList[C.TypeVLESS] {
				inbound.VLESSOptions.Users = append(inbound.VLESSOptions.Users, option.VLESSUser{
					Name: strconv.Itoa(int(user.Id)),
					UUID: user.Password,
				})
			}

			if inbound.VLESSOptions.Transport != nil {
				switch inbound.VLESSOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					inbound.VLESSOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeHTTPUpgrade:
					inbound.VLESSOptions.Transport.HTTPUpgradeOptions.Path = "/" + inbound.Type
				}
			}
		case C.TypeHysteria2:
			inbound.Hysteria2Options.ListenPort = uint16(port)
			inbound.Hysteria2Options.Users = []option.Hysteria2User{}
			inbound.Hysteria2Options.TLS.ServerName = os.Getenv("DOMAIN")

			for _, user := range premiumList[C.TypeVLESS] {
				inbound.Hysteria2Options.Users = append(inbound.Hysteria2Options.Users, option.Hysteria2User{
					Name:     strconv.Itoa(int(user.Id)),
					Password: user.Password,
				})
			}
		}

		inbounds = append(inbounds, inbound)
	}

	// Generate reality inbounds
	for i, inbound := range inbounds {
		if !strings.Contains(inbound.Tag, "-") {
			for x, sni := range sniList {
				port := 53000 + (i * 1000) + x
				tlsOptions := &option.InboundTLSOptions{
					Enabled:    true,
					ServerName: sni,
					Reality: &option.InboundRealityOptions{
						Enabled: true,
						Handshake: option.InboundRealityHandshakeOptions{
							ServerOptions: option.ServerOptions{
								Server:     sni,
								ServerPort: 443,
							},
						},
						PrivateKey: RealityPrivateKey,
						ShortID:    RealityShortID,
					},
				}

				generatedInbound := inbound
				generatedInbound.Tag = generatedInbound.Type + "-reality-" + sni + " : " + strconv.Itoa(port)
				switch inbound.Type {
				case C.TypeTrojan:
					generatedInbound.TrojanOptions.ListenPort = uint16(port)
					generatedInbound.TrojanOptions.TLS = tlsOptions
				case C.TypeVMess:
					generatedInbound.VMessOptions.ListenPort = uint16(port)
					generatedInbound.VMessOptions.TLS = tlsOptions
				case C.TypeVLESS:
					generatedInbound.VLESSOptions.ListenPort = uint16(port)
					generatedInbound.VLESSOptions.TLS = tlsOptions
				}

				// Ignore some protocol
				switch generatedInbound.Type {
				case C.TypeHysteria2, C.TypeVMess:
				default:
					inbounds = append(inbounds, generatedInbound)
				}

			}
		}
	}

	for _, list := range premiumList {
		for _, user := range list {
			options.Experimental.V2RayAPI.Stats.Users = append(options.Experimental.V2RayAPI.Stats.Users, strconv.Itoa(int(user.Id)))
		}
	}

	options.Inbounds = inbounds
	options.Outbounds = []option.Outbound{
		{
			Type: C.TypeDirect,
			Tag:  "direct",
		},
		{
			Type: C.TypeBlock,
			Tag:  "block",
		},
		{
			Type: C.TypeDNS,
			Tag:  "dns-out",
		},
	}
	options.Outbounds = append(options.Outbounds, relayOutbounds...)

	options.Route = &option.RouteOptions{
		GeoIP: &option.GeoIPOptions{
			Path:           "/usr/local/etc/sing-box/geoip.db",
			DownloadURL:    "https://github.com/malikshi/sing-box-geo/releases/latest/download/geoip.db",
			DownloadDetour: "direct",
		},
		Geosite: &option.GeositeOptions{
			Path:           "/usr/local/etc/sing-box/geosite.db",
			DownloadURL:    "https://github.com/malikshi/sing-box-geo/releases/latest/download/geosite.db",
			DownloadDetour: "direct",
		},
		Rules: []option.Rule{
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Protocol: option.Listable[string]{"dns"},
					Outbound: "dns-out",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					IPCIDR:   option.Listable[string]{"1.1.1.1", "8.8.8.8"},
					Outbound: "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Port:     option.Listable[uint16]{53},
					Outbound: "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					DomainSuffix: option.Listable[string]{"googlesyndication.com"},
					Outbound:     "direct",
				},
			},
		},
		Final: "direct",
	}

	// Spesific route each server
	serverInfo := helper.GetIpInfo()
	serverCode := strings.Split(serverInfo.Org, " ")[0]
	switch serverCode {
	case "AS133800":
		for _, outbound := range options.Outbounds {
			if outbound.Tag == "SG" {
				options.Route.Rules = append(options.Route.Rules, []option.Rule{
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							Geosite:  option.Listable[string]{"google", "rule-playstore", "rule-streaming"},
							Outbound: "SG",
						},
					},
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							GeoIP:    option.Listable[string]{"google"},
							Outbound: "SG",
						},
					}}...)
				break
			}
		}
	}

	// Adblock for specific user
	adblockRules := option.Rule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			Geosite:  option.Listable[string]{"rule-ads", "oisd-full"},
			AuthUser: option.Listable[string]{},
			Outbound: "block",
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
					AuthUser: option.Listable[string]{},
					Network:  option.Listable[string]{"tcp"},
					Outbound: outbound.Tag,
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
	f, err := os.Create("/usr/local/etc/latinaserver/config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(options, "", "\t")
	if err != nil {
		panic(err)
	}
	f.WriteString(string(b))

	return options
}
