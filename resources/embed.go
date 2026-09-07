package resources

import (
	_ "embed"
)

// DefaultCaddyTemplate is the embedded base Caddy JSON configuration template.
//
//go:embed caddy/caddy.json
var DefaultCaddyTemplate []byte

// DefaultSingBoxTemplate is the embedded base Sing-box JSON configuration template.
//
//go:embed sing-box/config.json
var DefaultSingBoxTemplate []byte
