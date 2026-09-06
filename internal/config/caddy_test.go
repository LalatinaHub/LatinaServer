package config

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaddyConfigTemplate_Validity(t *testing.T) {
	// Locate resources/caddy/caddy.json
	templateCandidates := []string{
		"../../resources/caddy/caddy.json",
		"resources/caddy/caddy.json",
	}

	var data []byte
	var err error
	for _, cand := range templateCandidates {
		data, err = os.ReadFile(cand)
		if err == nil {
			break
		}
	}
	require.NoError(t, err, "caddy template file must be readable")

	raw := string(data)
	assert.Contains(t, raw, "WEBSERVER_PORT", "template should reference WEBSERVER_PORT")
	assert.Contains(t, raw, "/portal", "template should route /portal to webserver")

	// Perform dummy replacement to check full JSON schema unmarshal
	dummyConfig := raw
	dummyConfig = strings.ReplaceAll(dummyConfig, `"DOMAIN_LIST"`, `"example.com"`)
	dummyConfig = strings.ReplaceAll(dummyConfig, "DOMAIN", "example.com")
	dummyConfig = strings.ReplaceAll(dummyConfig, "CF_KEY", "dummy-key")
	dummyConfig = strings.ReplaceAll(dummyConfig, "EMAIL", "admin@example.com")
	dummyConfig = strings.ReplaceAll(dummyConfig, "CADDY_LOGFILE", "/tmp/caddy.log")
	dummyConfig = strings.ReplaceAll(dummyConfig, "TROJAN_TCP_PORT", "10001")
	dummyConfig = strings.ReplaceAll(dummyConfig, "TROJAN_WS_PORT", "10002")
	dummyConfig = strings.ReplaceAll(dummyConfig, "TROJAN_HU_PORT", "10003")
	dummyConfig = strings.ReplaceAll(dummyConfig, "TROJAN_GRPC_PORT", "10004")
	dummyConfig = strings.ReplaceAll(dummyConfig, "TROJAN_UDP_PORT", "10005")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VMESS_TCP_PORT", "10006")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VMESS_WS_PORT", "10007")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VMESS_HU_PORT", "10008")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VMESS_GRPC_PORT", "10009")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VMESS_UDP_PORT", "10010")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VLESS_TCP_PORT", "10011")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VLESS_WS_PORT", "10012")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VLESS_HU_PORT", "10013")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VLESS_GRPC_PORT", "10014")
	dummyConfig = strings.ReplaceAll(dummyConfig, "VLESS_UDP_PORT", "10015")
	dummyConfig = strings.ReplaceAll(dummyConfig, "CADDY_PORT", "8080")
	dummyConfig = strings.ReplaceAll(dummyConfig, "WEBSERVER_PORT", "8081")

	var parsed map[string]any
	err = json.Unmarshal([]byte(dummyConfig), &parsed)
	require.NoError(t, err, "processed caddy template must be valid JSON")

	// Ensure fallback route points to WEBSERVER_PORT
	apps, ok := parsed["apps"].(map[string]any)
	require.True(t, ok)
	httpApp, ok := apps["http"].(map[string]any)
	require.True(t, ok)
	servers, ok := httpApp["servers"].(map[string]any)
	require.True(t, ok)
	srv0, ok := servers["srv0"].(map[string]any)
	require.True(t, ok)
	routes, ok := srv0["routes"].([]any)
	require.True(t, ok)

	// Last route is fallback
	lastRoute := routes[len(routes)-1].(map[string]any)
	handles := lastRoute["handle"].([]any)
	firstHandle := handles[0].(map[string]any)
	assert.Equal(t, "reverse_proxy", firstHandle["handler"])
	upstreams := firstHandle["upstreams"].([]any)
	firstUpstream := upstreams[0].(map[string]any)
	assert.Equal(t, "localhost:8081", firstUpstream["dial"])
}
