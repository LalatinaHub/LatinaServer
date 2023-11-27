package relay

import (
	"strings"

	"github.com/LalatinaHub/LatinaApi/common/account/converter"
	supabase "github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaSub-go/account"
	db "github.com/LalatinaHub/LatinaSub-go/db"
	"github.com/LalatinaHub/LatinaSub-go/provider"
	"github.com/LalatinaHub/LatinaSub-go/sandbox"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	Relays            []db.DBScheme
	excludedRelayCode = []string{helper.GetIpInfo().CountryCode}
)

func GatherRelays() {
	var (
		proxies        []db.DBScheme
		relayCodeCount = map[string]int{}
		ipServerList   = []string{}
	)

	supabase.Connect().DB.From("proxies").Select("*").Eq("conn_mode", "sni").Execute(&proxies)

	p := proxies
	proxies = []db.DBScheme{}
	for _, proxy := range p {
		isExists := func() bool {
			for _, ip := range ipServerList {
				if ip == proxy.Ip || proxy.Ip == "" {
					return true
				}
			}
			return false
		}()
		isExcluded := func() bool {
			for _, cc := range excludedRelayCode {
				if cc == proxy.CountryCode {
					return true
				}
			}
			return false
		}()

		if relayCodeCount[proxy.CountryCode] < 10 && !isExcluded && !isExists {
			proxies = append(proxies, proxy)
			relayCodeCount[proxy.CountryCode]++

			ipServerList = append(ipServerList, proxy.Ip)
		}
	}

	Relays = []db.DBScheme{}
	for i, node := range strings.Split(converter.ToRaw(proxies), "\n") {
		go func(i int, node string) {
			out, err := provider.Parse(node)
			if err != nil {
				return
			}

			if len(out) > 0 {
				box := sandbox.Test(out[0])
				if len(box.ConnectMode) > 0 {
					Relays = append(Relays, proxies[i])
				}
			}
		}(i, node)
	}
}

func GetRelayOutbounds() []option.Outbound {
	var (
		proxies      = Relays
		outbounds    = []option.Outbound{}
		outboundsMap = map[string][]option.Outbound{}
	)

	for _, proxy := range proxies {
		if len(outboundsMap[proxy.CountryCode]) < 5 {
			node := converter.ToRaw([]db.DBScheme{proxy})
			out, err := provider.Parse(node)
			if err != nil {
				continue
			}

			outboundsMap[proxy.CountryCode] = append(outboundsMap[proxy.CountryCode], account.New(out[0]).Outbound)
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
