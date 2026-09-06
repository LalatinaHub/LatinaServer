package warp

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

const (
	EndpointTag = "warp-out"
)

// BuildEndpoint converts WarpConfig into sing-box WireGuard Endpoint configuration.
func BuildEndpoint(cfg *Config) (*option.Endpoint, error) {
	if cfg == nil {
		return nil, fmt.Errorf("warp config cannot be nil")
	}

	// Parse Local IPv4 / IPv6 addresses
	var prefixes []netip.Prefix
	if cfg.IPv4 != "" {
		p4, err := netip.ParsePrefix(cfg.IPv4)
		if err == nil {
			prefixes = append(prefixes, p4)
		}
	}
	if cfg.IPv6 != "" {
		p6, err := netip.ParsePrefix(cfg.IPv6)
		if err == nil {
			prefixes = append(prefixes, p6)
		}
	}

	if len(prefixes) == 0 {
		// Fallback to default Cloudflare WARP local address
		prefixes = []netip.Prefix{
			netip.MustParsePrefix(DefaultIPv4),
		}
	}

	// Parse peer endpoint host and port
	host, portStr, err := net.SplitHostPort(cfg.Endpoint)
	var port uint16 = 2408
	if err == nil {
		if p, parseErr := strconv.Atoi(portStr); parseErr == nil && p > 0 && p <= 65535 {
			port = uint16(p)
		}
	} else {
		host = cfg.Endpoint
	}

	peer := option.WireGuardPeer{
		Address:   host,
		Port:      port,
		PublicKey: cfg.PeerPublicKey,
		AllowedIPs: badoption.Listable[netip.Prefix]{
			netip.MustParsePrefix("0.0.0.0/0"),
			netip.MustParsePrefix("::/0"),
		},
		Reserved: cfg.Reserved[:],
	}

	endpointOpts := option.WireGuardEndpointOptions{
		Address:    badoption.Listable[netip.Prefix](prefixes),
		PrivateKey: cfg.PrivateKey,
		Peers:      []option.WireGuardPeer{peer},
		MTU:        1280,
	}

	return &option.Endpoint{
		Type:    C.TypeWireGuard,
		Tag:     EndpointTag,
		Options: &endpointOpts,
	}, nil
}

// BuildSmartRule constructs the routing rule for AI and streaming services.
func BuildSmartRule(targetOutbound string) option.Rule {
	if targetOutbound == "" {
		targetOutbound = EndpointTag
	}

	domains := GetAllUnlockerDomains()

	return option.Rule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			RawDefaultRule: option.RawDefaultRule{
				DomainSuffix: badoption.Listable[string](domains),
			},
			RuleAction: option.RuleAction{
				Action: "route",
				RouteOptions: option.RouteActionOptions{
					Outbound: targetOutbound,
				},
			},
		},
	}
}
