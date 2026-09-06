package config

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/adblock"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/warp"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func TestReadSingConfig_ResourceTemplate(t *testing.T) {
	// Find repo root relative to internal/config
	templatePath := filepath.Join("..", "..", "resources", "sing-box", "config.json")
	opts, err := ReadSingConfig(templatePath)
	if err != nil {
		t.Fatalf("Failed to read sing-box resource config: %v", err)
	}

	if len(opts.Inbounds) == 0 {
		t.Fatalf("Expected inbounds in sing-box config, got 0")
	}

	if len(opts.Outbounds) == 0 {
		t.Fatalf("Expected outbounds in sing-box config, got 0")
	}

	if opts.Log == nil {
		t.Fatalf("Expected log options in sing-box config")
	}

	if opts.Log.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", opts.Log.Level)
	}

	if opts.Log.Output != "" {
		t.Errorf("Expected empty log output for terminal streaming, got '%s'", opts.Log.Output)
	}

	if len(opts.Endpoints) == 0 {
		t.Fatalf("Expected endpoints in sing-box config template, got 0")
	}

	hasAdguardDNS := false
	for _, s := range opts.DNS.Servers {
		if s.Tag == "adguard-dns" {
			hasAdguardDNS = true
			break
		}
	}
	if !hasAdguardDNS {
		t.Errorf("Expected 'adguard-dns' server in template DNS options")
	}

	hasWarpEndpoint := false
	for _, ep := range opts.Endpoints {
		if ep.Tag == "warp-out" && ep.Type == "wireguard" {
			hasWarpEndpoint = true
			break
		}
	}
	if !hasWarpEndpoint {
		t.Errorf("Expected 'warp-out' wireguard endpoint in template")
	}

	hasWarpRule := false
	for _, r := range opts.Route.Rules {
		if r.DefaultOptions.RuleAction.RouteOptions.Outbound == "warp-out" && len(r.DefaultOptions.RawDefaultRule.DomainSuffix) > 0 {
			hasWarpRule = true
			break
		}
	}
	if !hasWarpRule {
		t.Errorf("Expected route rule targeting 'warp-out' with unlocker domains")
	}
}

func TestSingContext_WireGuardSupport(t *testing.T) {
	ctx := SingContext(context.Background())
	if ctx == nil {
		t.Fatalf("SingContext returned nil")
	}

	privKey, pubKey, err := warp.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate keypair: %v", err)
	}

	cfg := &warp.Config{
		PrivateKey:    privKey,
		PublicKey:     pubKey,
		IPv4:          warp.DefaultIPv4,
		IPv6:          warp.DefaultIPv6,
		Endpoint:      warp.DefaultWarpEndpoint,
		PeerPublicKey: warp.DefaultPeerPublicKey,
	}

	ep, err := warp.BuildEndpoint(cfg)
	if err != nil {
		t.Fatalf("Failed to build WARP endpoint: %v", err)
	}

	opts := option.Options{
		Endpoints: []option.Endpoint{*ep},
		Outbounds: []option.Outbound{
			{
				Type: C.TypeDirect,
				Tag:  "direct",
			},
		},
		Route: &option.RouteOptions{
			Final: "direct",
		},
	}

	instance, err := box.New(box.Options{
		Context: ctx,
		Options: opts,
	})
	if err != nil {
		if strings.Contains(err.Error(), "gVisor is not included in this build") {
			t.Skip("Skipping WireGuard device instantiation test: build with '-tags with_gvisor' to test userspace gVisor device")
		}
		t.Fatalf("box.New failed with WireGuard endpoint: %v", err)
	}
	defer instance.Close()
}

func TestAdblockRules_Integration(t *testing.T) {
	adUsers := []string{"101", "102"}
	dnsRule := adblock.BuildAdblockDNSRule(adUsers)
	routeRule := adblock.BuildAdblockRouteRule(adUsers)

	// Verify DNS rule
	if dnsRule.DefaultOptions.RouteOptions.Server != adblock.AdGuardDNSTag {
		t.Errorf("Expected DNS server %s, got %s", adblock.AdGuardDNSTag, dnsRule.DefaultOptions.RouteOptions.Server)
	}
	if len(dnsRule.DefaultOptions.AuthUser) != 2 {
		t.Errorf("Expected 2 auth users in DNS rule, got %d", len(dnsRule.DefaultOptions.AuthUser))
	}

	// Verify Route rule
	if routeRule.DefaultOptions.RuleAction.Action != C.RuleActionTypeReject {
		t.Errorf("Expected reject action in route rule, got %s", routeRule.DefaultOptions.RuleAction.Action)
	}
	if len(routeRule.DefaultOptions.AuthUser) != 2 {
		t.Errorf("Expected 2 auth users in route rule, got %d", len(routeRule.DefaultOptions.AuthUser))
	}
	if len(routeRule.DefaultOptions.DomainSuffix) == 0 {
		t.Errorf("Expected non-empty DomainSuffix in adblock route rule")
	}
}


