package constant

var (
	baseLocation = "/usr/local/etc/latinaserver/"

	CADDY_CONFIG_PATH        = baseLocation + "caddy.json"
	CADDY_ACTIVE_CONFIG_PATH = baseLocation + "caddy-active.json"
	SING_CONFIG_PATH         = baseLocation + "config.json"
	SING_ACTIVE_CONFIG_PATH  = baseLocation + "config-active.json"

	SING_LOG_PATH  = baseLocation + "singbox.log"
	CADDY_LOG_PATH = baseLocation + "caddy.log"
)
