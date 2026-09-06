package proxy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/google/uuid"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

var supportedShadowsocksMethods = map[string]struct{}{
	"2022-blake3-aes-128-gcm":       {},
	"2022-blake3-aes-256-gcm":       {},
	"2022-blake3-chacha20-poly1305": {},
	"aes-128-gcm":                   {},
	"aes-192-gcm":                   {},
	"aes-256-gcm":                   {},
	"chacha20-ietf-poly1305":        {},
	"chacha20-poly1305":             {},
	"xchacha20-ietf-poly1305":       {},
	"none":                          {},
}

// ConverterSingbox converts ProxyNode to sing-box outbound options.
type ConverterSingbox struct{}

// NewConverterSingbox creates a new sing-box converter.
func NewConverterSingbox() *ConverterSingbox {
	return &ConverterSingbox{}
}

func validateServerAndPort(server string, port int) error {
	if strings.TrimSpace(server) == "" {
		return errors.New("server address cannot be empty")
	}
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be between 1 and 65535)", port)
	}
	return nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(strings.TrimSpace(u))
	return err == nil
}

// ConvertToOutbound converts a ProxyNode to sing-box Outbound option with strict validation.
func (c *ConverterSingbox) ConvertToOutbound(node model.ProxyNode, tag string) (option.Outbound, error) {
	if strings.TrimSpace(tag) == "" {
		return option.Outbound{}, errors.New("outbound tag cannot be empty")
	}

	switch strings.ToLower(strings.TrimSpace(node.VPN)) {
	case "shadowsocks":
		return c.convertShadowsocks(node, tag)
	case "vmess":
		return c.convertVMess(node, tag)
	case "vless":
		return c.convertVLESS(node, tag)
	case "trojan":
		return c.convertTrojan(node, tag)
	default:
		return option.Outbound{}, fmt.Errorf("unsupported VPN type: %s", node.VPN)
	}
}

func (c *ConverterSingbox) convertShadowsocks(node model.ProxyNode, tag string) (option.Outbound, error) {
	if err := validateServerAndPort(node.Server, node.ServerPort); err != nil {
		return option.Outbound{}, err
	}
	if strings.TrimSpace(node.Password) == "" {
		return option.Outbound{}, errors.New("shadowsocks password cannot be empty")
	}

	method := strings.ToLower(strings.TrimSpace(node.Method))
	if _, ok := supportedShadowsocksMethods[method]; !ok {
		return option.Outbound{}, fmt.Errorf("unsupported shadowsocks method: %s", node.Method)
	}

	return option.Outbound{
		Type: "shadowsocks",
		Tag:  tag,
		Options: option.ShadowsocksOutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     strings.TrimSpace(node.Server),
				ServerPort: uint16(node.ServerPort),
			},
			Method:   method,
			Password: node.Password,
		},
	}, nil
}

func (c *ConverterSingbox) convertVMess(node model.ProxyNode, tag string) (option.Outbound, error) {
	if err := validateServerAndPort(node.Server, node.ServerPort); err != nil {
		return option.Outbound{}, err
	}
	if !isValidUUID(node.UUID) {
		return option.Outbound{}, fmt.Errorf("invalid VMess UUID: %s", node.UUID)
	}

	vmessOpts := option.VMessOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     strings.TrimSpace(node.Server),
			ServerPort: uint16(node.ServerPort),
		},
		UUID:    strings.TrimSpace(node.UUID),
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
		path := node.Path
		if path == "" {
			path = "/"
		}
		vmessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    path,
				Headers: headers,
			},
		}
	}

	return option.Outbound{
		Type:    "vmess",
		Tag:     tag,
		Options: vmessOpts,
	}, nil
}

func (c *ConverterSingbox) convertVLESS(node model.ProxyNode, tag string) (option.Outbound, error) {
	if err := validateServerAndPort(node.Server, node.ServerPort); err != nil {
		return option.Outbound{}, err
	}
	if !isValidUUID(node.UUID) {
		return option.Outbound{}, fmt.Errorf("invalid VLESS UUID: %s", node.UUID)
	}

	vlessOpts := option.VLESSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     strings.TrimSpace(node.Server),
			ServerPort: uint16(node.ServerPort),
		},
		UUID: strings.TrimSpace(node.UUID),
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
		path := node.Path
		if path == "" {
			path = "/"
		}
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    path,
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
	}, nil
}

func (c *ConverterSingbox) convertTrojan(node model.ProxyNode, tag string) (option.Outbound, error) {
	if err := validateServerAndPort(node.Server, node.ServerPort); err != nil {
		return option.Outbound{}, err
	}
	if strings.TrimSpace(node.Password) == "" {
		return option.Outbound{}, errors.New("trojan password cannot be empty")
	}

	trojanOpts := option.TrojanOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     strings.TrimSpace(node.Server),
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
		path := node.Path
		if path == "" {
			path = "/"
		}
		trojanOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    path,
				Headers: headers,
			},
		}
	}

	return option.Outbound{
		Type:    "trojan",
		Tag:     tag,
		Options: trojanOpts,
	}, nil
}
