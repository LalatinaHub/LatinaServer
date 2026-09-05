package config

import (
	"path/filepath"
	"testing"
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

	if opts.Route == nil || len(opts.Route.Rules) == 0 {
		t.Fatalf("Expected route rules in sing-box config")
	}
}
