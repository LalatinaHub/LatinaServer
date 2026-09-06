package adblock

import (
	"testing"

	C "github.com/sagernet/sing-box/constant"
)

func TestDomains(t *testing.T) {
	if len(AdAndTrackingDomains) == 0 {
		t.Fatal("AdAndTrackingDomains should not be empty")
	}
	if len(MalwareAndPhishingDomains) == 0 {
		t.Fatal("MalwareAndPhishingDomains should not be empty")
	}
	if len(GamblingAndScamDomains) == 0 {
		t.Fatal("GamblingAndScamDomains should not be empty")
	}

	all := GetAllAdblockDomains()
	if len(all) == 0 {
		t.Fatal("GetAllAdblockDomains should not be empty")
	}

	seen := make(map[string]bool)
	for _, d := range all {
		if seen[d] {
			t.Errorf("Duplicate domain in GetAllAdblockDomains: %s", d)
		}
		seen[d] = true
	}
}

func TestBuildAdGuardDNSServer(t *testing.T) {
	server := BuildAdGuardDNSServer()
	if server.Tag != AdGuardDNSTag {
		t.Errorf("Expected tag %s, got %s", AdGuardDNSTag, server.Tag)
	}
	if server.Type != C.DNSTypeUDP {
		t.Errorf("Expected type %s, got %s", C.DNSTypeUDP, server.Type)
	}
}

func TestBuildAdblockDNSRule(t *testing.T) {
	users := []string{"1", "2", "3"}
	rule := BuildAdblockDNSRule(users)

	if rule.Type != C.RuleTypeDefault {
		t.Errorf("Expected rule type %s, got %s", C.RuleTypeDefault, rule.Type)
	}
	if len(rule.DefaultOptions.AuthUser) != 3 {
		t.Errorf("Expected 3 auth users, got %d", len(rule.DefaultOptions.AuthUser))
	}
	if rule.DefaultOptions.RouteOptions.Server != AdGuardDNSTag {
		t.Errorf("Expected DNS server %s, got %s", AdGuardDNSTag, rule.DefaultOptions.RouteOptions.Server)
	}
}

func TestBuildAdblockRouteRule(t *testing.T) {
	users := []string{"10", "20"}
	rule := BuildAdblockRouteRule(users)

	if rule.Type != C.RuleTypeDefault {
		t.Errorf("Expected rule type %s, got %s", C.RuleTypeDefault, rule.Type)
	}
	if rule.DefaultOptions.RuleAction.Action != C.RuleActionTypeReject {
		t.Errorf("Expected action %s, got %s", C.RuleActionTypeReject, rule.DefaultOptions.RuleAction.Action)
	}
	if len(rule.DefaultOptions.AuthUser) != 2 {
		t.Errorf("Expected 2 auth users, got %d", len(rule.DefaultOptions.AuthUser))
	}
	if len(rule.DefaultOptions.DomainSuffix) == 0 {
		t.Errorf("Expected non-empty DomainSuffix in adblock rule")
	}
}
