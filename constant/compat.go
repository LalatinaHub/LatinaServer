package constant

import "github.com/LalatinaHub/LatinaServer/internal/config"

var (
	CADDY_CONFIG_PATH        = config.CaddyConfigPath
	CADDY_ACTIVE_CONFIG_PATH = config.CaddyActiveConfigPath
	SING_CONFIG_PATH         = config.SingConfigPath
	SING_ACTIVE_CONFIG_PATH  = config.SingActiveConfigPath

	SING_LOG_PATH  = config.SingLogPath
	CADDY_LOG_PATH = config.CaddyLogPath

	IP_RESOLVER_DOMAIN = config.IPResolverDomain
	IP_RESOLVER_PATH   = config.IPResolverPath
)

const (
	V2RAY_API_ADDRESS = config.V2RayAPIAddress
	CLASH_API_ADDRESS = config.ClashAPIAddress

	CADDY_PORT     = config.CaddyPort
	MIXED_PORT     = config.MixedPort
	WEBSERVER_PORT = config.WebServerPort

	TROJAN_TCP_PORT  = config.TrojanTCPPort
	TROJAN_WS_PORT   = config.TrojanWSPort
	TROJAN_HU_PORT   = config.TrojanHUPort
	TROJAN_GRPC_PORT = config.TrojanGRPCPort
	VMESS_TCP_PORT   = config.VMessTCPPort
	VMESS_WS_PORT    = config.VMessWSPort
	VMESS_HU_PORT    = config.VMessHUPort
	VMESS_GRPC_PORT  = config.VMessGRPCPort
	VLESS_TCP_PORT   = config.VLESS_TCP_PORT_COMPAT
	VLESS_WS_PORT    = config.VLESS_WS_PORT_COMPAT
	VLESS_HU_PORT    = config.VLESS_HU_PORT_COMPAT
	VLESS_GRPC_PORT  = config.VLESS_GRPC_PORT_COMPAT

	SS_TCP_PORT     = config.SSTCPPort
	TROJAN_UDP_PORT = config.TrojanUDPPort
	VLESS_UDP_PORT  = config.VLessUDPPort
	VMess_UDP_PORT  = config.VMessUDPPort
	MIXED_IN_PORT   = config.MixedInPort

	SERVICE_LATINASERVER = config.ServiceLatinaServer
)
