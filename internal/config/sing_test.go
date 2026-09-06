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

	if opts.Log == nil {
		t.Fatalf("Expected log options in sing-box config")
	}

	if opts.Log.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", opts.Log.Level)
	}

	if opts.Log.Output != "" {
		t.Errorf("Expected empty log output for terminal streaming, got '%s'", opts.Log.Output)
	}
}
