package config

import (
	"os"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

var LogOptions = &option.LogOptions{
	Level:  "info",
	Output: "/usr/local/etc/latinaserver/singbox.log",
}
var DNSOptions = &option.DNSOptions{
	Servers: []option.DNSServerOptions{
		{
			Address: "tls://1.1.1.1",
		},
	},
}
var NTPOptions = &option.NTPOptions{
	Enabled: true,
	ServerOptions: option.ServerOptions{
		Server:     "time.apple.com",
		ServerPort: 123,
	},
}
var ListenOptions = option.ListenOptions{
	Listen:             &badoption.Addr{},
	ListenPort:         52001,
	TCPFastOpen:        true,
	UDPFragmentDefault: true,
}
var InboundsOptions = []option.Inbound{
	{
		Type: C.TypeMixed,
		Tag:  C.TypeMixed,
		Options: option.HTTPMixedInboundOptions{
			ListenOptions: ListenOptions,
		},
	},
	{
		Type: C.TypeTrojan,
		Tag:  C.TypeTrojan,
		Options: option.TrojanInboundOptions{
			ListenOptions: ListenOptions,
		},
	},
	{
		Type: C.TypeTrojan,
		Tag:  C.TypeTrojan + "-ws",
		Options: option.TrojanInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:             C.V2RayTransportTypeWebsocket,
				WebsocketOptions: WSOptions,
			},
		},
	},
	{
		Type: C.TypeTrojan,
		Tag:  C.TypeTrojan + "-hu",
		Options: option.TrojanInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:               C.V2RayTransportTypeHTTPUpgrade,
				HTTPUpgradeOptions: HUOptions,
			},
		},
	},
	{
		Type: C.TypeTrojan,
		Tag:  C.TypeTrojan + "-grpc",
		Options: option.TrojanInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:        C.V2RayTransportTypeGRPC,
				GRPCOptions: GRPCOptions,
			},
		},
	},
	{
		Type: C.TypeVMess,
		Tag:  C.TypeVMess,
		Options: option.VMessInboundOptions{
			ListenOptions: ListenOptions,
		},
	},
	{
		Type: C.TypeVMess,
		Tag:  C.TypeVMess + "-ws",
		Options: option.VMessInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:             C.V2RayTransportTypeWebsocket,
				WebsocketOptions: WSOptions,
			},
		},
	},
	{
		Type: C.TypeVMess,
		Tag:  C.TypeVMess + "-hu",
		Options: option.VMessInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:               C.V2RayTransportTypeHTTPUpgrade,
				HTTPUpgradeOptions: HUOptions,
			},
		},
	},
	{
		Type: C.TypeVMess,
		Tag:  C.TypeVMess + "-grpc",
		Options: option.VMessInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:        C.V2RayTransportTypeGRPC,
				GRPCOptions: GRPCOptions,
			},
		},
	},
	{
		Type: C.TypeVLESS,
		Tag:  C.TypeVLESS,
		Options: option.VLESSInboundOptions{
			ListenOptions: ListenOptions,
		},
	},
	{
		Type: C.TypeVLESS,
		Tag:  C.TypeVLESS + "-ws",
		Options: option.VLESSInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:             C.V2RayTransportTypeWebsocket,
				WebsocketOptions: WSOptions,
			},
		},
	},
	{
		Type: C.TypeVLESS,
		Tag:  C.TypeVLESS + "-hu",
		Options: option.VLESSInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:               C.V2RayTransportTypeHTTPUpgrade,
				HTTPUpgradeOptions: HUOptions,
			},
		},
	},
	{
		Type: C.TypeVLESS,
		Tag:  C.TypeVLESS + "-grpc",
		Options: option.VLESSInboundOptions{
			ListenOptions: ListenOptions,
			Transport: &option.V2RayTransportOptions{
				Type:        C.V2RayTransportTypeGRPC,
				GRPCOptions: GRPCOptions,
			},
		},
	},
}
var OutboundsOptions = []option.Outbound{
	{
		Type: C.TypeDirect,
		Tag:  C.TypeDirect,
	},
	{
		Type: C.TypeBlock,
		Tag:  C.TypeBlock,
	},
	{
		Type: C.TypeDNS,
		Tag:  C.TypeDNS,
	},
}
var MultiplexOptions = &option.InboundMultiplexOptions{
	Enabled: true,
	Padding: false,
	Brutal: &option.BrutalOptions{
		Enabled:  false,
		UpMbps:   100,
		DownMbps: 100,
	},
}
var RouteOptions = &option.RouteOptions{
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
				RawDefaultRule: option.RawDefaultRule{
					Protocol: badoption.Listable[string]{"dns"},
				},
				RuleAction: option.RuleAction{
					RouteOptions: option.RouteActionOptions{
						Outbound: C.TypeDNS,
					},
				},
			},
		},
		{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					Port: badoption.Listable[uint16]{53},
				},
				RuleAction: option.RuleAction{
					RouteOptions: option.RouteActionOptions{
						Outbound: C.TypeDirect,
					},
				},
			},
		},
		{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					DomainSuffix: badoption.Listable[string]{"googlesyndication.com"},
				},
				RuleAction: option.RuleAction{
					RouteOptions: option.RouteActionOptions{
						Outbound: C.TypeDirect,
					},
				},
			},
		},
	},
	Final: C.TypeDirect,
}
var ExperimentalOptions = &option.ExperimentalOptions{
	ClashAPI: &option.ClashAPIOptions{
		ExternalController: CS.ClashAPIAddress,
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

var WSOptions = option.V2RayWebsocketOptions{}
var HUOptions = option.V2RayHTTPUpgradeOptions{}
var GRPCOptions = option.V2RayGRPCOptions{}

var SingOptions = option.Options{
	Log:          LogOptions,
	DNS:          DNSOptions,
	NTP:          NTPOptions,
	Inbounds:     InboundsOptions,
	Outbounds:    OutboundsOptions,
	Route:        RouteOptions,
	Experimental: ExperimentalOptions,
}
