package helper

import (
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/geoip"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/v2ray"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/LalatinaHub/LatinaServer/pkg/util"
)

// Aliases for compatibility
type (
	MyIp        = geoip.MyIp
	Countries   = geoip.Countries
	GeoIpJson   = geoip.GeoIpJson
	ProxyInfo   = util.ProxyInfo
	ProxyIPInfo = util.ProxyIPInfo
	ServerStat  = util.ServerStat
)

var (
	CountryList   = geoip.CountryList
	Parse         = geoip.Parse
	GetIpInfo     = geoip.GetIpInfo
	GetUserStats  = v2ray.GetUserStats
	ReloadService = systemctl.Reload
	CatchError    = systemctl.CatchError
	CheckProxyIP  = util.CheckProxyIP
	GetServerStatus = util.GetServerStatus
)

func RemoveDuplicate[T comparable](sliceList []T) []T {
	return util.RemoveDuplicate(sliceList)
}

