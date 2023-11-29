package relay

import (
	"github.com/LalatinaHub/LatinaApi/common/account/converter"
	supabase "github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaSub-go/account"
	db "github.com/LalatinaHub/LatinaSub-go/db"
	"github.com/LalatinaHub/LatinaSub-go/provider"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	Relays            = []db.DBScheme{}
	excludedRelayCode = []string{helper.GetIpInfo().CountryCode}
)

func GatherRelays() {
	var (
		proxies        []db.DBScheme
		relayCodeCount = map[string]int{}
	)

	supabase.Connect().DB.From("proxies").Select("*").Neq("vpn", "shadowsocks").Execute(&proxies)

	for _, proxy := range proxies {
		isExcluded := func() bool {
			for _, cc := range excludedRelayCode {
				if cc == proxy.CountryCode {
					return true
				}
			}
			return false
		}()

		if relayCodeCount[proxy.CountryCode] < 10 && !isExcluded {
			Relays = append(Relays, proxy)
			relayCodeCount[proxy.CountryCode]++
		}
	}
}

func GetRelayOutbounds() []option.Outbound {
	var (
		proxies      = Relays
		outbounds    = []option.Outbound{}
		outboundsMap = map[string][]option.Outbound{}
	)

	if len(proxies) == 0 {
		return outbounds
	}

	for _, proxy := range proxies {
		if len(outboundsMap[proxy.CountryCode]) < 5 {
			node := converter.ToRaw([]db.DBScheme{proxy})
			out, err := provider.Parse(node)
			if err != nil {
				continue
			}

			for _, o := range out {
				outboundsMap[proxy.CountryCode] = append(outboundsMap[proxy.CountryCode], account.New(o).Outbound)
			}
		}
	}

	for cc, out := range outboundsMap {
		urltest := option.Outbound{
			Tag:  cc,
			Type: C.TypeURLTest,
			URLTestOptions: option.URLTestOutboundOptions{
				Outbounds: []string{},
			},
		}

		for _, outbound := range out {
			urltest.URLTestOptions.Outbounds = append(urltest.URLTestOptions.Outbounds, outbound.Tag)
		}
		outbounds = append(outbounds, urltest)
		outbounds = append(outbounds, out...)
	}

	return outbounds
}
