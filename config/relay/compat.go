package relay

import (
	"github.com/LalatinaHub/LatinaServer/internal/config/relay"
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/sagernet/sing-box/option"
)

func GetRelays() []model.ProxyNode {
	return relay.GetRelays()
}

func GatherRelays() {
	relay.GatherRelays()
}

func GetRelayOutbounds() []option.Outbound {
	return relay.GetRelayOutbounds()
}

// Re-export for backward compatibility
type ProxyNode = model.ProxyNode
