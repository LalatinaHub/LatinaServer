package constant

var (
	baseLocation = "/usr/local/etc/latinaserver/"

	CaddyConfigPath       = baseLocation + "caddy.json"
	CaddyActiveConfigPath = baseLocation + "caddy-active.json"
	SingConfigPath        = baseLocation + "config.json"
	SingActiveConfigPath  = baseLocation + "config-active.json"

	SingLogPath  = baseLocation + "singbox.log"
	CaddyLogPath = baseLocation + "caddy.log"
)
