package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
)

var (
	cachedSSPassword string
	cachedMu         sync.Mutex
)

// GetSSPassword returns shadowsocks password from environment variable or generates a secure random one.
func GetSSPassword() string {
	cachedMu.Lock()
	defer cachedMu.Unlock()

	if val := os.Getenv("SS_PASSWORD"); val != "" {
		return val
	}

	if cachedSSPassword != "" {
		return cachedSSPassword
	}

	// Generate 16 secure random bytes (32 hex characters)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback deterministic if rand fails (extremely rare)
		return "latina-sec-ss-default-key"
	}

	cachedSSPassword = hex.EncodeToString(b)
	return cachedSSPassword
}

// GetClashSecret returns Clash API secret from environment.
func GetClashSecret() string {
	if val := os.Getenv("CLASH_SECRET"); val != "" {
		return val
	}
	if val := os.Getenv("YACD_PASSWORD"); val != "" {
		return val
	}
	return ""
}
