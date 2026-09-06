package adblock

import (
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

const (
	AdGuardDNSTag = "adguard-dns"
	AdGuardDNSIP  = "94.140.14.14"
)

// BuildAdGuardDNSServer creates an upstream AdGuard DNS server option.
func BuildAdGuardDNSServer() option.DNSServerOptions {
	return option.DNSServerOptions{
		Type: C.DNSTypeUDP,
		Tag:  AdGuardDNSTag,
		Options: &option.RemoteDNSServerOptions{
			DNSServerAddressOptions: option.DNSServerAddressOptions{
				Server:     AdGuardDNSIP,
				ServerPort: 53,
			},
		},
	}
}

// BuildAdblockDNSRule creates a DNS rule routing queries from adblock users to AdGuard DNS.
func BuildAdblockDNSRule(authUsers []string) option.DNSRule {
	return option.DNSRule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultDNSRule{
			RawDefaultDNSRule: option.RawDefaultDNSRule{
				AuthUser: badoption.Listable[string](authUsers),
			},
			DNSRuleAction: option.DNSRuleAction{
				Action: C.RuleActionTypeRoute,
				RouteOptions: option.DNSRouteActionOptions{
					Server: AdGuardDNSTag,
				},
			},
		},
	}
}

// BuildAdblockRouteRule creates a route rule rejecting connections from adblock users to ad/malware domains.
func BuildAdblockRouteRule(authUsers []string) option.Rule {
	domains := GetAllAdblockDomains()

	return option.Rule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			RawDefaultRule: option.RawDefaultRule{
				AuthUser:     badoption.Listable[string](authUsers),
				DomainSuffix: badoption.Listable[string](domains),
			},
			RuleAction: option.RuleAction{
				Action: C.RuleActionTypeReject,
			},
		},
	}
}
