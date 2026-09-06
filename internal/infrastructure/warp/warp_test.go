package warp

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair(t *testing.T) {
	priv, pub, err := GenerateKeyPair()
	require.NoError(t, err)
	assert.NotEmpty(t, priv)
	assert.NotEmpty(t, pub)

	privBytes, err := base64.StdEncoding.DecodeString(priv)
	require.NoError(t, err)
	assert.Len(t, privBytes, 32)

	pubBytes, err := base64.StdEncoding.DecodeString(pub)
	require.NoError(t, err)
	assert.Len(t, pubBytes, 32)
}

func TestDomains(t *testing.T) {
	domains := GetAllUnlockerDomains()
	assert.NotEmpty(t, domains)
	assert.Contains(t, domains, "openai.com")
	assert.Contains(t, domains, "chatgpt.com")
	assert.Contains(t, domains, "netflix.com")
	assert.Contains(t, domains, "disneyplus.com")
	assert.Contains(t, domains, "reddit.com")
}

func TestBuildEndpointAndRule(t *testing.T) {
	cfg := &Config{
		PrivateKey:    "cGFzc3dvcmQxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ=",
		PublicKey:     "cHVibGljMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY=",
		IPv4:          "172.16.0.2/32",
		IPv6:          "2606:4700:110:8f81:85e8:bb2d:5c55:6834/128",
		Endpoint:      "162.159.192.1:2408",
		PeerPublicKey: "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
		Reserved:      [3]uint8{0, 0, 0},
	}

	ep, err := BuildEndpoint(cfg)
	require.NoError(t, err)
	assert.Equal(t, "wireguard", ep.Type)
	assert.Equal(t, EndpointTag, ep.Tag)
	assert.NotNil(t, ep.Options)

	rule := BuildSmartRule(EndpointTag)
	assert.Equal(t, "default", rule.Type)
	assert.Equal(t, "route", rule.DefaultOptions.RuleAction.Action)
	assert.Equal(t, EndpointTag, rule.DefaultOptions.RuleAction.RouteOptions.Outbound)
	assert.Contains(t, []string(rule.DefaultOptions.RawDefaultRule.DomainSuffix), "openai.com")
}

func TestSaveAndLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test_warp.json")

	cfg := &Config{
		PrivateKey:    "priv-key",
		PublicKey:     "pub-key",
		IPv4:          "172.16.0.2/32",
		Endpoint:      DefaultWarpEndpoint,
		PeerPublicKey: DefaultPeerPublicKey,
	}

	err := SaveConfig(path, cfg)
	require.NoError(t, err)

	loaded, err := LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, cfg.PrivateKey, loaded.PrivateKey)
	assert.Equal(t, cfg.Endpoint, loaded.Endpoint)
}

func TestHealthStatus(t *testing.T) {
	SetHealthy(true)
	assert.True(t, IsHealthy())

	SetHealthy(false)
	assert.False(t, IsHealthy())

	SetHealthy(true)
	assert.True(t, IsHealthy())
}
