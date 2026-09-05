package cf

import (
	"github.com/LalatinaHub/LatinaServer/internal/infrastructure/cloudflare"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	cf "github.com/cloudflare/cloudflare-go"
)

// Compatibility wrapper
type cloudflareStruct = cloudflare.Client

func MakeCloudflareAPIClient() *cloudflare.Client {
	client, err := cloudflare.NewClient()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create cloudflare client")
	}
	return client
}

// Re-export cloudflare types for backward compatibility
type Zone = cf.Zone

