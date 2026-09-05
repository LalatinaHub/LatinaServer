package proxy_test

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/proxy"
)

func TestConverterSingbox_ConvertToOutbound(t *testing.T) {
	converter := proxy.NewConverterSingbox()

	tests := []struct {
		name      string
		node      model.ProxyNode
		tag       string
		wantErr   bool
		checkType string
	}{
		{
			name: "shadowsocks outbound",
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
			name: "vmess outbound",
			node: model.ProxyNode{
				VPN:        "vmess",
				Server:     "example.com",
				ServerPort: 443,
				UUID:       "test-uuid",
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
			name: "vless outbound",
			node: model.ProxyNode{
				VPN:        "vless",
				Server:     "example.com",
				ServerPort: 443,
				UUID:       "test-uuid",
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
			name: "trojan outbound",
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
