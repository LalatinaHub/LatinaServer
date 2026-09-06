package proxy_test

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/proxy"
)

func TestConverterSingbox_ConvertToOutbound(t *testing.T) {
	converter := proxy.NewConverterSingbox()

	validUUID := "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d"

	tests := []struct {
		name      string
		node      model.ProxyNode
		tag       string
		wantErr   bool
		checkType string
	}{
		{
			name: "shadowsocks outbound - valid",
			node: model.ProxyNode{
				VPN:        "shadowsocks",
				Server:     "example.com",
				ServerPort: 8388,
				Method:     "aes-128-gcm",
				Password:   "password",
			},
			tag:       "ss-test",
			wantErr:   false,
			checkType: "shadowsocks",
		},
		{
			name: "shadowsocks outbound - 2022 method valid",
			node: model.ProxyNode{
				VPN:        "shadowsocks",
				Server:     "example.com",
				ServerPort: 8388,
				Method:     "2022-blake3-aes-128-gcm",
				Password:   "password",
			},
			tag:       "ss-2022-test",
			wantErr:   false,
			checkType: "shadowsocks",
		},
		{
			name: "shadowsocks outbound - unsupported method",
			node: model.ProxyNode{
				VPN:        "shadowsocks",
				Server:     "example.com",
				ServerPort: 8388,
				Method:     "rc4-md5",
				Password:   "password",
			},
			tag:     "ss-bad-method",
			wantErr: true,
		},
		{
			name: "shadowsocks outbound - empty password",
			node: model.ProxyNode{
				VPN:        "shadowsocks",
				Server:     "example.com",
				ServerPort: 8388,
				Method:     "aes-128-gcm",
				Password:   "",
			},
			tag:     "ss-no-pass",
			wantErr: true,
		},
		{
			name: "vmess outbound - valid",
			node: model.ProxyNode{
				VPN:        "vmess",
				Server:     "example.com",
				ServerPort: 443,
				UUID:       validUUID,
				AlterID:    0,
				TLS:        true,
				Transport:  "ws",
				Host:       "example.com",
				Path:       "/ws",
			},
			tag:       "vmess-test",
			wantErr:   false,
			checkType: "vmess",
		},
		{
			name: "vmess outbound - invalid UUID",
			node: model.ProxyNode{
				VPN:        "vmess",
				Server:     "example.com",
				ServerPort: 443,
				UUID:       "not-a-valid-uuid",
			},
			tag:     "vmess-bad-uuid",
			wantErr: true,
		},
		{
			name: "vless outbound - valid",
			node: model.ProxyNode{
				VPN:        "vless",
				Server:     "example.com",
				ServerPort: 443,
				UUID:       validUUID,
				TLS:        true,
				Transport:  "ws",
				Host:       "example.com",
				Path:       "/ws",
			},
			tag:       "vless-test",
			wantErr:   false,
			checkType: "vless",
		},
		{
			name: "vless outbound - invalid UUID",
			node: model.ProxyNode{
				VPN:        "vless",
				Server:     "example.com",
				ServerPort: 443,
				UUID:       "",
			},
			tag:     "vless-empty-uuid",
			wantErr: true,
		},
		{
			name: "trojan outbound - valid",
			node: model.ProxyNode{
				VPN:        "trojan",
				Server:     "example.com",
				ServerPort: 443,
				Password:   "trojan-password",
				TLS:        true,
				Transport:  "ws",
				Host:       "example.com",
				Path:       "/trojan",
			},
			tag:       "trojan-test",
			wantErr:   false,
			checkType: "trojan",
		},
		{
			name: "trojan outbound - empty password",
			node: model.ProxyNode{
				VPN:        "trojan",
				Server:     "example.com",
				ServerPort: 443,
				Password:   "",
			},
			tag:     "trojan-no-pass",
			wantErr: true,
		},
		{
			name: "invalid server address",
			node: model.ProxyNode{
				VPN:        "trojan",
				Server:     "   ",
				ServerPort: 443,
				Password:   "pass",
			},
			tag:     "trojan-no-server",
			wantErr: true,
		},
		{
			name: "invalid port (out of range)",
			node: model.ProxyNode{
				VPN:        "trojan",
				Server:     "example.com",
				ServerPort: 70000,
				Password:   "pass",
			},
			tag:     "trojan-bad-port",
			wantErr: true,
		},
		{
			name: "empty tag",
			node: model.ProxyNode{
				VPN:        "trojan",
				Server:     "example.com",
				ServerPort: 443,
				Password:   "pass",
			},
			tag:     "   ",
			wantErr: true,
		},
		{
			name: "unsupported outbound",
			node: model.ProxyNode{
				VPN: "unsupported",
			},
			tag:     "bad-test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outbound, err := converter.ConvertToOutbound(tt.node, tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToOutbound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if outbound.Type != tt.checkType {
					t.Errorf("Outbound.Type = %v, want %v", outbound.Type, tt.checkType)
				}
				if outbound.Tag != tt.tag {
					t.Errorf("Outbound.Tag = %v, want %v", outbound.Tag, tt.tag)
				}
				if outbound.Options == nil {
					t.Errorf("Outbound.Options is nil")
				}
			}
		})
	}
}
