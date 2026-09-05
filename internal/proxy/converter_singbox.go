package proxy

import (
	"fmt"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

// ConverterSingbox converts ProxyNode to sing-box outbound options.
type ConverterSingbox struct{}

// NewConverterSingbox creates a new sing-box converter.
func NewConverterSingbox() *ConverterSingbox {
	return &ConverterSingbox{}
}

// ConvertToOutbound converts a ProxyNode to sing-box Outbound option.
func (c *ConverterSingbox) ConvertToOutbound(node model.ProxyNode, tag string) (option.Outbound, error) {
	switch node.VPN {
	case "shadowsocks":
		return c.convertShadowsocks(node, tag), nil
	case "vmess":
		return c.convertVMess(node, tag), nil
	case "vless":
		return c.convertVLESS(node, tag), nil
	case "trojan":
		return c.convertTrojan(node, tag), nil
	default:
		return option.Outbound{}, fmt.Errorf("unsupported VPN type: %s", node.VPN)
	}
}

func (c *ConverterSingbox) convertShadowsocks(node model.ProxyNode, tag string) option.Outbound {
	return option.Outbound{
		Type: "shadowsocks",
		Tag:  tag,
		Options: option.ShadowsocksOutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     node.Server,
				ServerPort: uint16(node.ServerPort),
			},
			Method:   node.Method,
			Password: node.Password,
		},
	}
}

func (c *ConverterSingbox) convertVMess(node model.ProxyNode, tag string) option.Outbound {
	vmessOpts := option.VMessOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     node.Server,
			ServerPort: uint16(node.ServerPort),
		},
		UUID:    node.UUID,
		AlterId: node.AlterID,
	}

	if node.TLS {
		vmessOpts.TLS = &option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: node.SNI,
		}
	}

	if node.Transport == "ws" {
		headers := badoption.HTTPHeader{}
		if node.Host != "" {
			headers = badoption.HTTPHeader{"Host": badoption.Listable[string]{node.Host}}
		}
		vmessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    node.Path,
				Headers: headers,
			},
		}
	}

	return option.Outbound{
		Type:    "vmess",
		Tag:     tag,
		Options: vmessOpts,
	}
}

func (c *ConverterSingbox) convertVLESS(node model.ProxyNode, tag string) option.Outbound {
	vlessOpts := option.VLESSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     node.Server,
			ServerPort: uint16(node.ServerPort),
		},
		UUID: node.UUID,
	}

	if node.TLS {
		vlessOpts.TLS = &option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: node.SNI,
		}
	}

	if node.Transport == "ws" {
		headers := badoption.HTTPHeader{}
		if node.Host != "" {
			headers = badoption.HTTPHeader{"Host": badoption.Listable[string]{node.Host}}
		}
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    node.Path,
				Headers: headers,
			},
		}
	} else if node.Transport == "grpc" {
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "grpc",
			GRPCOptions: option.V2RayGRPCOptions{
				ServiceName: node.ServiceName,
			},
		}
	}

	return option.Outbound{
		Type:    "vless",
		Tag:     tag,
		Options: vlessOpts,
	}
}

func (c *ConverterSingbox) convertTrojan(node model.ProxyNode, tag string) option.Outbound {
	trojanOpts := option.TrojanOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     node.Server,
			ServerPort: uint16(node.ServerPort),
		},
		Password: node.Password,
	}

	if node.TLS {
		trojanOpts.TLS = &option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: node.SNI,
		}
	}

	if node.Transport == "ws" {
		headers := badoption.HTTPHeader{}
		if node.Host != "" {
			headers = badoption.HTTPHeader{"Host": badoption.Listable[string]{node.Host}}
		}
		trojanOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    node.Path,
				Headers: headers,
			},
		}
	}

	return option.Outbound{
		Type:    "trojan",
		Tag:     tag,
		Options: trojanOpts,
	}
}
