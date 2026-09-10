package constant

// File paths
var (
	baseLocation = "/usr/local/etc/latinaserver/"

	CaddyConfigPath       = baseLocation + "caddy.json"
	CaddyActiveConfigPath = baseLocation + "caddy-active.json"
	SingConfigPath        = baseLocation + "config.json"
	SingActiveConfigPath  = baseLocation + "config-active.json"

	SingLogPath  = baseLocation + "singbox.log"
	CaddyLogPath = baseLocation + "caddy.log"
)

// Network addresses & ports
const (
	V2RayAPIAddress = "0.0.0.0:5555"
	ClashAPIAddress = "0.0.0.0:9090"

	CaddyPort     = 50000
	MixedPort     = 7878
	WebServerPort = 5000

	TrojanTCPPort  = 52001
	TrojanWSPort   = 52002
	TrojanHUPort   = 52003
	TrojanGRPCPort = 52004
	VMessTCPPort   = 52005
	VMessWSPort    = 52006
	VMessHUPort    = 52007
	VMessGRPCPort  = 52008
	VLessTCPPort   = 52009
	VLessWSPort    = 52010
	VLessHUPort    = 52011
	VLessGRPCPort  = 52012

	SSTCPPort     = 53000
	TrojanUDPPort = 53001
	VLessUDPPort  = 53002
	VMessUDPPort  = 53003
	MixedInPort   = 53004
	SSTCPBrutalPort = 53005
)

// External services
const (
	IPResolverDomain = "myip.ipeek.workers.dev"
	IPResolverPath   = "/"
)

// System services
const (
	ServiceLatinaServer = "latinaserver"
)

// Backward compatibility aliases
const (
	VLESS_TCP_PORT_COMPAT  = VLessTCPPort
	VLESS_WS_PORT_COMPAT   = VLessWSPort
	VLESS_HU_PORT_COMPAT   = VLessHUPort
	VLESS_GRPC_PORT_COMPAT = VLessGRPCPort
)
