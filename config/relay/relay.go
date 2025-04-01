package relay

import (
	"encoding/json"
	"fmt"
	"slices"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
	mgpr "github.com/FoolVPN-ID/megalodon-api/modules/proxy"
	mgdb "github.com/FoolVPN-ID/megalodon/db"
	"github.com/FoolVPN-ID/tool/modules/subconverter"
	"github.com/LalatinaHub/LatinaServer/helper"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	Relays            = []mgdb.ProxyFieldStruct{}
	excludedRelayCode = []string{helper.GetIpInfo().CountryCode}
)

func GatherRelays() {
	Relays = []mgdb.ProxyFieldStruct{}
	var (
		proxies        []mgdb.ProxyFieldStruct
		relayCodeCount = map[string]int{}
	)

	client := database.MakeDatabase().GetClient()
	rows, err := client.Query("SELECT * FROM proxies WHERE vpn = 'shadowsocks'")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var (
			result = mgdb.ProxyFieldStruct{}
			id     int
		)

		err := rows.Scan(
			&id,
			&result.Server,
			&result.Ip,
			&result.ServerPort,
			&result.UUID,
			&result.Password,
			&result.Security,
			&result.AlterId,
			&result.Method,
			&result.Plugin,
			&result.PluginOpts,
			&result.Host,
			&result.TLS,
			&result.Transport,
			&result.Path,
			&result.ServiceName,
			&result.Insecure,
			&result.SNI,
			&result.Remark,
			&result.ConnMode,
			&result.CountryCode,
			&result.Region,
			&result.Org,
			&result.VPN,
			&result.Raw,
		)

		if err != nil {
			continue
		}

		proxies = append(proxies, result)
	}

	for _, proxy := range proxies {
		isExcluded := func() bool {
			return slices.Contains(excludedRelayCode, proxy.CountryCode)
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
			node := mgpr.ConvertDBToURL(&proxy)
			subOut, err := subconverter.MakeSubconverterFromConfig(node.String())
			if err != nil {
				continue
			}

			for _, outbound := range subOut.Outbounds {
				if _, err := json.MarshalIndent(outbound, "", "\t"); err == nil {
					outboundsMap[proxy.CountryCode] = append(outboundsMap[proxy.CountryCode], outbound)
				} else {
					fmt.Printf("Error parsing: %s\n", node)
				}
			}
		}
	}

	for cc, out := range outboundsMap {
		var outboundTags = []string{}
		for _, outbound := range out {
			outboundTags = append(outboundTags, outbound.Tag)
		}

		urltest := option.Outbound{
			Tag:  cc,
			Type: C.TypeURLTest,
			Options: option.URLTestOutboundOptions{
				Outbounds: outboundTags,
			},
		}

		outbounds = append(outbounds, urltest)
		outbounds = append(outbounds, out...)
	}

	return outbounds
}
