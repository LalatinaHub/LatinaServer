package proxy_test

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/proxy"
)

func TestParser_ParseShadowsocks(t *testing.T) {
	parser := proxy.NewParser()

	tests := []struct {
		name    string
		url     string
		wantErr bool
		check   func(*testing.T, *model.ProxyNode)
	}{
		{
			name:    "valid shadowsocks with base64 full URL",
			url:     "ss://YWVzLTEyOC1nY206dGVzdHBhc3N3b3JkQGV4YW1wbGUuY29tOjgzODg=#TestServer",
			wantErr: false,
			check: func(t *testing.T, node *model.ProxyNode) {
				if node.Method != "aes-128-gcm" {
					t.Errorf("Method = %v, want aes-128-gcm", node.Method)
				}
				if node.Password != "testpassword" {
					t.Errorf("Password = %v, want testpassword", node.Password)
				}
				if node.Server != "example.com" {
					t.Errorf("Server = %v, want example.com", node.Server)
				}
				if node.ServerPort != 8388 {
					t.Errorf("ServerPort = %v, want 8388", node.ServerPort)
				}
				if node.Remark != "TestServer" {
					t.Errorf("Remark = %v, want TestServer", node.Remark)
				}
				if node.VPN != "shadowsocks" {
					t.Errorf("VPN = %v, want shadowsocks", node.VPN)
				}
			},
		},
		{
			name:    "valid shadowsocks with base64 userinfo only",
			url:     "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#TestServer2",
			wantErr: false,
			check: func(t *testing.T, node *model.ProxyNode) {
				if node.Method != "aes-256-gcm" {
					t.Errorf("Method = %v, want aes-256-gcm", node.Method)
				}
				if node.Password != "password" {
					t.Errorf("Password = %v, want password", node.Password)
				}
				if node.Server != "example.com" {
					t.Errorf("Server = %v, want example.com", node.Server)
				}
				if node.ServerPort != 8388 {
					t.Errorf("ServerPort = %v, want 8388", node.ServerPort)
				}
				if node.Remark != "TestServer2" {
					t.Errorf("Remark = %v, want TestServer2", node.Remark)
				}
			},
		},
		{
			name:    "valid shadowsocks without base64",
			url:     "ss://aes-256-cfb:password123@192.168.1.1:8388",
			wantErr: false,
			check: func(t *testing.T, node *model.ProxyNode) {
				if node.Method != "aes-256-cfb" {
					t.Errorf("Method = %v, want aes-256-cfb", node.Method)
				}
				if node.Password != "password123" {
					t.Errorf("Password = %v, want password123", node.Password)
				}
				if node.Server != "192.168.1.1" {
					t.Errorf("Server = %v, want 192.168.1.1", node.Server)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parser.Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, node)
			}
		})
	}
}

func TestParser_ParseVMess(t *testing.T) {
	parser := proxy.NewParser()

	// VMess JSON: {"add":"example.com","port":"443","id":"uuid-test","aid":"0","net":"ws","type":"none","host":"example.com","path":"/path","tls":"tls","ps":"Test"}
	vmessURL := "vmess://eyJhZGQiOiJleGFtcGxlLmNvbSIsInBvcnQiOiI0NDMiLCJpZCI6InV1aWQtdGVzdCIsImFpZCI6IjAiLCJuZXQiOiJ3cyIsInR5cGUiOiJub25lIiwiaG9zdCI6ImV4YW1wbGUuY29tIiwicGF0aCI6Ii9wYXRoIiwidGxzIjoidGxzIiwicHMiOiJUZXN0In0="

	node, err := parser.Parse(vmessURL)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if node.VPN != "vmess" {
		t.Errorf("VPN = %v, want vmess", node.VPN)
	}
	if node.Server != "example.com" {
		t.Errorf("Server = %v, want example.com", node.Server)
	}
	if node.ServerPort != 443 {
		t.Errorf("ServerPort = %v, want 443", node.ServerPort)
	}
	if node.UUID != "uuid-test" {
		t.Errorf("UUID = %v, want uuid-test", node.UUID)
	}
	if node.Transport != "ws" {
		t.Errorf("Transport = %v, want ws", node.Transport)
	}
	if node.Path != "/path" {
		t.Errorf("Path = %v, want /path", node.Path)
	}
	if !node.TLS {
		t.Errorf("TLS = %v, want true", node.TLS)
	}
	if node.Remark != "Test" {
		t.Errorf("Remark = %v, want Test", node.Remark)
	}
}

func TestParser_ParseVLESS(t *testing.T) {
	parser := proxy.NewParser()

	vlessURL := "vless://uuid-test@example.com:443?type=ws&security=tls&path=/ws&host=example.com&sni=example.com#TestVLESS"

	node, err := parser.Parse(vlessURL)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if node.VPN != "vless" {
		t.Errorf("VPN = %v, want vless", node.VPN)
	}
	if node.Server != "example.com" {
		t.Errorf("Server = %v, want example.com", node.Server)
	}
	if node.ServerPort != 443 {
		t.Errorf("ServerPort = %v, want 443", node.ServerPort)
	}
	if node.UUID != "uuid-test" {
		t.Errorf("UUID = %v, want uuid-test", node.UUID)
	}
	if node.Transport != "ws" {
		t.Errorf("Transport = %v, want ws", node.Transport)
	}
	if node.Path != "/ws" {
		t.Errorf("Path = %v, want /ws", node.Path)
	}
	if node.Security != "tls" {
		t.Errorf("Security = %v, want tls", node.Security)
	}
	if !node.TLS {
		t.Errorf("TLS = %v, want true", node.TLS)
	}
	if node.Remark != "TestVLESS" {
		t.Errorf("Remark = %v, want TestVLESS", node.Remark)
	}
}

func TestParser_ParseTrojan(t *testing.T) {
	parser := proxy.NewParser()

	trojanURL := "trojan://password123@example.com:443?type=ws&security=tls&path=/trojan&host=example.com#TestTrojan"

	node, err := parser.Parse(trojanURL)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if node.VPN != "trojan" {
		t.Errorf("VPN = %v, want trojan", node.VPN)
	}
	if node.Server != "example.com" {
		t.Errorf("Server = %v, want example.com", node.Server)
	}
	if node.ServerPort != 443 {
		t.Errorf("ServerPort = %v, want 443", node.ServerPort)
	}
	if node.Password != "password123" {
		t.Errorf("Password = %v, want password123", node.Password)
	}
	if node.Transport != "ws" {
		t.Errorf("Transport = %v, want ws", node.Transport)
	}
	if node.Path != "/trojan" {
		t.Errorf("Path = %v, want /trojan", node.Path)
	}
	if !node.TLS {
		t.Errorf("TLS = %v, want true", node.TLS)
	}
	if node.Remark != "TestTrojan" {
		t.Errorf("Remark = %v, want TestTrojan", node.Remark)
	}
}

func TestParser_UnsupportedProtocol(t *testing.T) {
	parser := proxy.NewParser()

	_, err := parser.Parse("http://example.com")
	if err == nil {
		t.Error("Expected error for unsupported protocol, got nil")
	}
}

func TestParser_EdgeCases(t *testing.T) {
	parser := proxy.NewParser()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "invalid shadowsocks missing @",
			url:     "ss://bm9hdHNpZ24=",
			wantErr: true,
		},
		{
			name:    "invalid shadowsocks missing colon in method password",
			url:     "ss://bm9jb2xvbg==@1.1.1.1:8388",
			wantErr: true,
		},
		{
			name:    "invalid shadowsocks invalid port",
			url:     "ss://aes-128-gcm:pass@1.1.1.1:notaport",
			wantErr: true,
		},
		{
			name:    "invalid vmess bad base64",
			url:     "vmess://!!!invalidbase64",
			wantErr: true,
		},
		{
			name:    "invalid vmess invalid json",
			url:     "vmess://bm90LWpzb24=",
			wantErr: true,
		},
		{
			name:    "invalid vless missing @",
			url:     "vless://noatsign-example.com:443",
			wantErr: true,
		},
		{
			name:    "invalid vless missing port",
			url:     "vless://uuid@example.com",
			wantErr: true,
		},
		{
			name:    "invalid trojan missing @",
			url:     "trojan://noatsign-example.com:443",
			wantErr: true,
		},
		{
			name:    "invalid trojan invalid port",
			url:     "trojan://pass@example.com:abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%s) error = %v, wantErr = %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

