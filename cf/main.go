package cf

import (
	"context"
	"errors"
	"os"

	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/cloudflare/cloudflare-go"
)

type cloudflareStruct struct {
	client cloudflare.API
	cfCtx  context.Context
}

func MakeCloudflareAPIClient() *cloudflareStruct {
	client, err := cloudflare.New(os.Getenv("GLOBAL_CF_KEY"), os.Getenv("EMAIL"))
	if err != nil {
		panic(err.Error())
	}

	return &cloudflareStruct{
		client: *client,
		cfCtx:  context.Background(),
	}
}

func (cf *cloudflareStruct) GetAllZones() []cloudflare.Zone {
	zones, err := cf.client.ListZones(cf.cfCtx)
	if err != nil {
		panic(err)
	}

	return zones
}

func (cf *cloudflareStruct) GetAssignedDomain() []string {
	var (
		result = []string{}
		zones  = cf.GetAllZones()
		geoip  = helper.GetIpInfo()
	)

	for _, zone := range zones {
		domains, info, err := cf.client.ListDNSRecords(cf.cfCtx, cloudflare.ZoneIdentifier(zone.ID), cloudflare.ListDNSRecordsParams{
			Content: geoip.Ip,
		})
		if err != nil {
			continue
		}

		if info.Count > 0 {
			for _, domain := range domains {
				result = append(result, domain.Name)
			}
		}
	}

	if len(result) <= 0 {
		panic(errors.New("no dns records found for this server's ip"))
	}

	return result
}
