package cloudflare

import (
	"context"

	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/geoip"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	cf "github.com/cloudflare/cloudflare-go"
)

type Client struct {
	api   *cf.API
	cfCtx context.Context
}

func NewClient() (*Client, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, appErrors.NewDatabaseError("failed to get database connection", err)
	}

	kvRepo := repository.NewKVRepository(db)
	kvList, err := kvRepo.GetAll(context.Background())
	if err != nil {
		return nil, appErrors.NewDatabaseError("failed to get KV list", err)
	}

	api, err := cf.New(kvList["globalCfKey"].(string), kvList["email"].(string))
	if err != nil {
		return nil, appErrors.NewNetworkError("failed to create cloudflare API client", err)
	}

	return &Client{
		api:   api,
		cfCtx: context.Background(),
	}, nil
}

func (c *Client) GetAllZones() ([]cf.Zone, error) {
	zones, err := c.api.ListZones(c.cfCtx)
	if err != nil {
		return nil, appErrors.NewNetworkError("failed to list cloudflare zones", err)
	}

	return zones, nil
}

func (c *Client) GetAssignedDomain() ([]string, error) {
	var result = []string{}
	
	zones, err := c.GetAllZones()
	if err != nil {
		return nil, err
	}

	ipInfo := geoip.GetIpInfo()

	for _, zone := range zones {
		domains, info, err := c.api.ListDNSRecords(c.cfCtx, cf.ZoneIdentifier(zone.ID), cf.ListDNSRecordsParams{
			Content: ipInfo.Ip,
		})
		if err != nil {
			logger.Warn().Err(err).Str("zone", zone.Name).Msg("Failed to list DNS records for zone")
			continue
		}

		if info.Count > 0 {
			for _, domain := range domains {
				result = append(result, domain.Name)
			}
		}
	}

	if len(result) <= 0 {
		return nil, appErrors.NewValidationError("no DNS records found for this server's IP", nil)
	}

	return result, nil
}

