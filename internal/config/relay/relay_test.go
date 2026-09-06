package relay

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestGetRelayOutbounds(t *testing.T) {
	t.Run("empty relays returns empty slice", func(t *testing.T) {
		setRelays(nil)
		outbounds := GetRelayOutbounds()
		assert.Empty(t, outbounds)
	})

	t.Run("skips invalid relay nodes and avoids duplicate tags", func(t *testing.T) {
		validUUID := "11111111-2222-3333-4444-555555555555"

		mockNodes := []model.ProxyNode{
			// Valid vmess node
			{
				VPN:         "vmess",
				CountryCode: "SG",
				Server:      "sg1.example.com",
				ServerPort:  443,
				UUID:        validUUID,
				Remark:      "SG-Relay",
			},
			// Duplicate remark with node 1 (should be deduplicated to SG-Relay-2)
			{
				VPN:         "vmess",
				CountryCode: "SG",
				Server:      "sg2.example.com",
				ServerPort:  443,
				UUID:        validUUID,
				Remark:      "SG-Relay",
			},
			// Node with reserved tag "direct" (should be sanitized to direct-2)
			{
				VPN:         "trojan",
				CountryCode: "US",
				Server:      "us1.example.com",
				ServerPort:  443,
				Password:    "password123",
				Remark:      "direct",
			},
			// Broken node: invalid UUID (must be skipped)
			{
				VPN:         "vmess",
				CountryCode: "JP",
				Server:      "jp1.example.com",
				ServerPort:  443,
				UUID:        "invalid-uuid-format",
				Remark:      "JP-Broken",
			},
			// Broken node: missing server (must be skipped)
			{
				VPN:         "trojan",
				CountryCode: "JP",
				Server:      "",
				ServerPort:  443,
				Password:    "pass",
				Remark:      "JP-No-Server",
			},
			// Broken node: invalid shadowsocks cipher (must be skipped)
			{
				VPN:         "shadowsocks",
				CountryCode: "JP",
				Server:      "jp2.example.com",
				ServerPort:  8388,
				Method:      "rc4-md5",
				Password:    "pass",
				Remark:      "JP-Bad-Cipher",
			},
		}

		setRelays(mockNodes)
		defer setRelays(nil)

		outbounds := GetRelayOutbounds()
		assert.NotEmpty(t, outbounds)

		// Verify that none of the outbounds have duplicate tags
		seenTags := make(map[string]bool)
		for _, ob := range outbounds {
			assert.False(t, seenTags[ob.Tag], "duplicate tag detected: %s", ob.Tag)
			seenTags[ob.Tag] = true
			// Assert that reserved tags were not used directly as outbounds
			assert.NotEqual(t, "direct", ob.Tag)
			assert.NotEqual(t, "ss-out", ob.Tag)
		}

		// Broken nodes should not exist in the outbounds
		assert.False(t, seenTags["JP-Broken"])
		assert.False(t, seenTags["JP-No-Server"])
		assert.False(t, seenTags["JP-Bad-Cipher"])

		// Valid nodes should be present with deduplicated tags
		assert.True(t, seenTags["SG-Relay"])
		assert.True(t, seenTags["SG-Relay-2"])
		assert.True(t, seenTags["direct-2"])

		// URLTest tags should exist for SG and US
		assert.True(t, seenTags["SG"])
		assert.True(t, seenTags["US"])
	})
}
