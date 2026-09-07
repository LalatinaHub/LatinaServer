//go:build !race

package caddy_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	caddy "github.com/caddyserver/caddy/v2"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/mholt/caddy-l4"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestCertificate generates a temporary self-signed TLS cert and key for localhost and 127.0.0.1.
func generateTestCertificate(certFile, keyFile string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"LatinaServer Test"},
			CommonName:   "127.0.0.1",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost", "127.0.0.1"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	return pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
}

func getFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return port
}

type VPNTestCase struct {
	Name         string
	OutboundJSON string
}

func TestVPN_EndToEnd_ThroughCaddy(t *testing.T) {
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "cert.pem")
	keyFile := filepath.Join(tempDir, "key.pem")
	require.NoError(t, generateTestCertificate(certFile, keyFile))

	// Target Echo Server
	echoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ECHO_LATINA_SUCCESS"))
	}))
	defer echoServer.Close()

	// Ports definition for isolated test
	gatePort := getFreePort(t)
	demuxPort := getFreePort(t)
	caddyHttpPort := getFreePort(t)
	webserverPort := getFreePort(t)

	trojanTcpPort := getFreePort(t)
	trojanWsPort := getFreePort(t)
	trojanHuPort := getFreePort(t)
	trojanGrpcPort := getFreePort(t)

	vmessTcpPort := getFreePort(t)
	vmessWsPort := getFreePort(t)
	vmessHuPort := getFreePort(t)
	vmessGrpcPort := getFreePort(t)

	vlessTcpPort := getFreePort(t)
	vlessWsPort := getFreePort(t)
	vlessHuPort := getFreePort(t)
	vlessGrpcPort := getFreePort(t)

	ssTcpPort := getFreePort(t)

	testPassword := "550e8400-e29b-41d4-a716-446655440000"

	// 1. Build Sing-box Server Config
	singServerJSON := fmt.Sprintf(`{
		"log": { "level": "warn" },
		"inbounds": [
			{ "type": "trojan", "tag": "trojan-tcp", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "password": "%s"}] },
			{ "type": "trojan", "tag": "trojan-ws", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "password": "%s"}], "transport": {"type": "ws", "path": "/trojan"} },
			{ "type": "trojan", "tag": "trojan-hu", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "password": "%s"}], "transport": {"type": "httpupgrade", "path": "/trojan"} },
			{ "type": "trojan", "tag": "trojan-grpc", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "password": "%s"}], "transport": {"type": "grpc", "service_name": "trojan"} },

			{ "type": "vmess", "tag": "vmess-tcp", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}] },
			{ "type": "vmess", "tag": "vmess-ws", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}], "transport": {"type": "ws", "path": "/vmess"} },
			{ "type": "vmess", "tag": "vmess-hu", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}], "transport": {"type": "httpupgrade", "path": "/vmess"} },
			{ "type": "vmess", "tag": "vmess-grpc", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}], "transport": {"type": "grpc", "service_name": "vmess"} },

			{ "type": "vless", "tag": "vless-tcp", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}] },
			{ "type": "vless", "tag": "vless-ws", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}], "transport": {"type": "ws", "path": "/vless"} },
			{ "type": "vless", "tag": "vless-hu", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}], "transport": {"type": "httpupgrade", "path": "/vless"} },
			{ "type": "vless", "tag": "vless-grpc", "listen": "127.0.0.1", "listen_port": %d, "users": [{"name": "u1", "uuid": "%s"}], "transport": {"type": "grpc", "service_name": "vless"} },

			{ "type": "shadowsocks", "tag": "ss-in", "listen": "127.0.0.1", "listen_port": %d, "method": "aes-256-gcm", "password": "%s" }
		],
		"outbounds": [
			{ "type": "direct", "tag": "direct" }
		],
		"route": { "final": "direct" }
	}`,
		trojanTcpPort, testPassword,
		trojanWsPort, testPassword,
		trojanHuPort, testPassword,
		trojanGrpcPort, testPassword,
		vmessTcpPort, testPassword,
		vmessWsPort, testPassword,
		vmessHuPort, testPassword,
		vmessGrpcPort, testPassword,
		vlessTcpPort, testPassword,
		vlessWsPort, testPassword,
		vlessHuPort, testPassword,
		vlessGrpcPort, testPassword,
		ssTcpPort, testPassword,
	)

	singCtx := config.SingContext(context.Background())
	var serverOpts option.Options
	require.NoError(t, serverOpts.UnmarshalJSONContext(singCtx, []byte(singServerJSON)))

	serverBox, err := box.New(box.Options{
		Context: singCtx,
		Options: serverOpts,
	})
	require.NoError(t, err)
	require.NoError(t, serverBox.Start())
	defer serverBox.Close()

	// 2. Build Caddy Config based on production template logic
	certFileEscaped := strings.ReplaceAll(certFile, `\`, `/`)
	keyFileEscaped := strings.ReplaceAll(keyFile, `\`, `/`)

	caddyJSON := fmt.Sprintf(`{
		"apps": {
			"tls": {
				"certificates": {
					"load_files": [
						{
							"certificate": "%s",
							"key": "%s"
						}
					]
				}
			},
			"layer4": {
				"servers": {
					"gate": {
						"listen": ["127.0.0.1:%d"],
						"routes": [
							{
								"handle": [
									{
										"handler": "tls",
										"connection_policies": [{ "default_sni": "127.0.0.1", "fallback_sni": "127.0.0.1" }]
									},
									{ "handler": "proxy", "upstreams": [{ "dial": ["127.0.0.1:%d"] }] }
								]
							}
						]
					},
					"srv0": {
						"listen": ["127.0.0.1:%d"],
						"routes": [
							{
								"match": [
									{
										"regexp": {
											"count": 58,
											"pattern": "^[0-9a-fA-F].*\\x0d\\x0a$"
										}
									}
								],
								"handle": [{ "handler": "proxy", "upstreams": [{ "dial": ["127.0.0.1:%d"] }] }]
							},
							{
								"match": [
									{
										"regexp": {
											"count": 1,
											"pattern": "\\x00"
										}
									}
								],
								"handle": [{ "handler": "proxy", "upstreams": [{ "dial": ["127.0.0.1:%d"] }] }]
							},
							{
								"match": [
									{
										"regexp": {
											"count": 4,
											"pattern": "^(GET |POST|HEAD|PUT |DELE|OPTI|PATC|PRI )"
										}
									}
								],
								"handle": [{ "handler": "proxy", "upstreams": [{ "dial": ["127.0.0.1:%d"] }] }]
							},
							{
								"handle": [{ "handler": "proxy", "upstreams": [{ "dial": ["127.0.0.1:%d"] }] }]
							}
						]
					}
				}
			},
			"http": {
				"servers": {
					"srv0": {
						"listen": [":%d"],
						"protocols": ["h1", "h2", "h2c", "h3"],
						"routes": [
							{
								"match": [{ "path": ["/trojan", "/trojan/"], "header": { "Sec-Websocket-Key": ["*"] } }],
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							},
							{
								"match": [{ "path": ["/trojan", "/trojan/"], "header": { "Upgrade": ["*"], "Sec-Websocket-Key": null } }],
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							},
							{
								"match": [{ "path": ["/trojan/Tun", "/trojan/Tun/*"] }],
								"handle": [
									{
										"handler": "reverse_proxy",
										"upstreams": [{ "dial": "127.0.0.1:%d" }],
										"transport": { "protocol": "http", "versions": ["h2c"] }
									}
								],
								"terminal": true
							},
							{
								"match": [{ "path": ["/vmess", "/vmess/"], "header": { "Sec-Websocket-Key": ["*"] } }],
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							},
							{
								"match": [{ "path": ["/vmess", "/vmess/"], "header": { "Upgrade": ["*"], "Sec-Websocket-Key": null } }],
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							},
							{
								"match": [{ "path": ["/vmess/Tun", "/vmess/Tun/*"] }],
								"handle": [
									{
										"handler": "reverse_proxy",
										"upstreams": [{ "dial": "127.0.0.1:%d" }],
										"transport": { "protocol": "http", "versions": ["h2c"] }
									}
								],
								"terminal": true
							},
							{
								"match": [{ "path": ["/vless", "/vless/"], "header": { "Sec-Websocket-Key": ["*"] } }],
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							},
							{
								"match": [{ "path": ["/vless", "/vless/"], "header": { "Upgrade": ["*"], "Sec-Websocket-Key": null } }],
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							},
							{
								"match": [{ "path": ["/vless/Tun", "/vless/Tun/*"] }],
								"handle": [
									{
										"handler": "reverse_proxy",
										"upstreams": [{ "dial": "127.0.0.1:%d" }],
										"transport": { "protocol": "http", "versions": ["h2c"] }
									}
								],
								"terminal": true
							},
							{
								"handle": [{ "handler": "reverse_proxy", "upstreams": [{ "dial": "127.0.0.1:%d" }] }],
								"terminal": true
							}
						]
					}
				}
			}
		}
	}`,
		certFileEscaped, keyFileEscaped,
		gatePort, demuxPort,
		demuxPort,
		trojanTcpPort,
		vlessTcpPort,
		caddyHttpPort,
		vmessTcpPort,
		caddyHttpPort,
		trojanWsPort,
		trojanHuPort,
		trojanGrpcPort,
		vmessWsPort,
		vmessHuPort,
		vmessGrpcPort,
		vlessWsPort,
		vlessHuPort,
		vlessGrpcPort,
		webserverPort,
	)

	var caddyConfig caddy.Config
	require.NoError(t, json.Unmarshal([]byte(caddyJSON), &caddyConfig))

	require.NoError(t, caddy.Run(&caddyConfig))
	defer caddy.Stop()

	// 3. Test Cases for all registered VPN protocols through Caddy Gate
	testCases := []VPNTestCase{
		{
			Name: "Trojan-TCP",
			OutboundJSON: fmt.Sprintf(`{
				"type": "trojan",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"password": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true }
			}`, gatePort, testPassword),
		},
		{
			Name: "Trojan-WS",
			OutboundJSON: fmt.Sprintf(`{
				"type": "trojan",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"password": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "ws", "path": "/trojan" }
			}`, gatePort, testPassword),
		},
		{
			Name: "Trojan-HTTPUpgrade",
			OutboundJSON: fmt.Sprintf(`{
				"type": "trojan",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"password": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "httpupgrade", "path": "/trojan" }
			}`, gatePort, testPassword),
		},
		{
			Name: "Trojan-gRPC",
			OutboundJSON: fmt.Sprintf(`{
				"type": "trojan",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"password": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "grpc", "service_name": "trojan" }
			}`, gatePort, testPassword),
		},
		{
			Name: "VMess-TCP",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vmess",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"security": "auto",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true }
			}`, gatePort, testPassword),
		},
		{
			Name: "VMess-WS",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vmess",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"security": "auto",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "ws", "path": "/vmess" }
			}`, gatePort, testPassword),
		},
		{
			Name: "VMess-HTTPUpgrade",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vmess",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"security": "auto",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "httpupgrade", "path": "/vmess" }
			}`, gatePort, testPassword),
		},
		{
			Name: "VMess-gRPC",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vmess",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"security": "auto",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "grpc", "service_name": "vmess" }
			}`, gatePort, testPassword),
		},
		{
			Name: "VLESS-TCP",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vless",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true }
			}`, gatePort, testPassword),
		},
		{
			Name: "VLESS-WS",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vless",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "ws", "path": "/vless" }
			}`, gatePort, testPassword),
		},
		{
			Name: "VLESS-HTTPUpgrade",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vless",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "httpupgrade", "path": "/vless" }
			}`, gatePort, testPassword),
		},
		{
			Name: "VLESS-gRPC",
			OutboundJSON: fmt.Sprintf(`{
				"type": "vless",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"uuid": "%s",
				"tls": { "enabled": true, "server_name": "127.0.0.1", "insecure": true },
				"transport": { "type": "grpc", "service_name": "vless" }
			}`, gatePort, testPassword),
		},
		{
			Name: "Shadowsocks-TCP",
			OutboundJSON: fmt.Sprintf(`{
				"type": "shadowsocks",
				"tag": "proxy",
				"server": "127.0.0.1",
				"server_port": %d,
				"method": "aes-256-gcm",
				"password": "%s"
			}`, ssTcpPort, testPassword),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			clientMixedPort := getFreePort(t)

			clientJSON := fmt.Sprintf(`{
				"log": { "level": "warn" },
				"inbounds": [
					{
						"type": "mixed",
						"tag": "mixed-in",
						"listen": "127.0.0.1",
						"listen_port": %d
					}
				],
				"outbounds": [
					%s,
					{ "type": "direct", "tag": "direct" }
				],
				"route": { "final": "proxy" }
			}`, clientMixedPort, tc.OutboundJSON)

			var clientOpts option.Options
			require.NoError(t, clientOpts.UnmarshalJSONContext(singCtx, []byte(clientJSON)))

			clientBox, err := box.New(box.Options{
				Context: singCtx,
				Options: clientOpts,
			})
			require.NoError(t, err)
			require.NoError(t, clientBox.Start())
			defer clientBox.Close()

			// Wait briefly for client listener to accept connections
			time.Sleep(50 * time.Millisecond)

			proxyURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", clientMixedPort))
			require.NoError(t, err)

			httpClient := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
				Timeout: 8 * time.Second,
			}

			req, err := http.NewRequest(http.MethodGet, echoServer.URL, nil)
			require.NoError(t, err)

			resp, err := httpClient.Do(req)
			require.NoError(t, err, "Traffic through Caddy and sing-box proxy tunnel failed")
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "ECHO_LATINA_SUCCESS", string(body))
		})
	}
}
