package config

import "github.com/LalatinaHub/LatinaServer/pkg/constant"

// File paths
var (
	CaddyConfigPath       = constant.CaddyConfigPath
	CaddyActiveConfigPath = constant.CaddyActiveConfigPath
	SingConfigPath        = constant.SingConfigPath
	SingActiveConfigPath  = constant.SingActiveConfigPath

	SingLogPath  = constant.SingLogPath
	CaddyLogPath = constant.CaddyLogPath
)

// Network addresses
const (
	V2RayAPIAddress = constant.V2RayAPIAddress
	ClashAPIAddress = constant.ClashAPIAddress

	CaddyPort     = constant.CaddyPort
	MixedPort     = constant.MixedPort
	WebServerPort = constant.WebServerPort

	TrojanTCPPort  = constant.TrojanTCPPort
	TrojanWSPort   = constant.TrojanWSPort
	TrojanHUPort   = constant.TrojanHUPort
	TrojanGRPCPort = constant.TrojanGRPCPort
	VMessTCPPort   = constant.VMessTCPPort
	VMessWSPort    = constant.VMessWSPort
	VMessHUPort    = constant.VMessHUPort
	VMessGRPCPort  = constant.VMessGRPCPort
	VLessTCPPort   = constant.VLessTCPPort
	VLessWSPort    = constant.VLessWSPort
	VLessHUPort    = constant.VLessHUPort
	VLessGRPCPort  = constant.VLessGRPCPort

	SSTCPPort     = constant.SSTCPPort
	TrojanUDPPort = constant.TrojanUDPPort
	VLessUDPPort  = constant.VLessUDPPort
	VMessUDPPort  = constant.VMessUDPPort
	MixedInPort     = constant.MixedInPort
	SSTCPBrutalPort = constant.SSTCPBrutalPort
)

// External services
const (
	IPResolverDomain = constant.IPResolverDomain
	IPResolverPath   = constant.IPResolverPath
)

// System services
const (
	ServiceLatinaServer = constant.ServiceLatinaServer
)

// Backward compatibility aliases
const (
	VLESS_TCP_PORT_COMPAT  = constant.VLESS_TCP_PORT_COMPAT
	VLESS_WS_PORT_COMPAT   = constant.VLESS_WS_PORT_COMPAT
	VLESS_HU_PORT_COMPAT   = constant.VLESS_HU_PORT_COMPAT
	VLESS_GRPC_PORT_COMPAT = constant.VLESS_GRPC_PORT_COMPAT
)


