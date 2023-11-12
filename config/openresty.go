package config

import (
	"os"
	"strconv"
	"strings"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	C "github.com/sagernet/sing-box/constant"
)

var locationTemplace = []string{
	`		location /PATH {`,
	`			proxy_redirect off;`,
	`			proxy_set_header X-Real-IP $remote_addr;`,
	`			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;`,
	`			proxy_set_header Upgrade $http_upgrade;`,
	`			proxy_set_header Connection $http_connection;`,
	`			proxy_set_header Host $http_host;`,
	``,
	`			set_by_lua $hasWebSocketKey 'if ngx.var.http_sec_websocket_key ~= nil then return "true" else return "false" end';`,
	`			set_by_lua $isGrpc 'if ngx.req.get_headers()["content-type"] == "application/grpc" then return "true" else return "false" end';`,
	``,
	`			if ($isGrpc = "true") {`,
	`				grpc_pass "grpc://127.0.0.1:GRPC_PORT";`,
	`				break;`,
	`			}`,
	`			if ($http_upgrade != "websocket") {`,
	`				proxy_pass "http://127.0.0.1:` + strconv.Itoa(CS.WebServerPort) + `";`,
	`				break;`,
	`			}`,
	`			if ($hasWebSocketKey = "true") {`,
	`				proxy_pass "http://127.0.0.1:WS_PORT";`,
	`				break;`,
	`			}`,
	`			if ($hasWebSocketKey = "false") {`,
	`				proxy_pass "http://127.0.0.1:HU_PORT";`,
	`				break;`,
	`			}`,
	`			ARG`,
	`		}`,
}

func WriteOpenrestyConfig() {
	var (
		locations = map[string]string{}
		config    = ReadSingConfig()
	)

	r, err := os.ReadFile("/usr/local/etc/latinaserver/nginx.conf")
	if err != nil {
		panic(err)
	}

	openrestyConfig := string(r)

	for _, inbound := range config.Inbounds {
		switch inbound.Type {
		case C.TypeTrojan:
			trojan := inbound.TrojanOptions
			if locations[C.TypeTrojan] == "" {
				locations[C.TypeTrojan] = strings.Join(locationTemplace[:], "\n")
			}

			if trojan.Transport != nil {
				location := locations[C.TypeTrojan]
				location = strings.Replace(location, "PATH", C.TypeTrojan, 1)

				switch trojan.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					location = strings.Replace(location, "WS_PORT", strconv.Itoa(int(trojan.ListenPort)), 1)
				case C.V2RayTransportTypeHTTPUpgrade:
					location = strings.Replace(location, "HU_PORT", strconv.Itoa(int(trojan.ListenPort)), 1)
				case C.V2RayTransportTypeGRPC:
					location = strings.Replace(location, "GRPC_PORT", strconv.Itoa(int(trojan.ListenPort)), 1)
				}

				locations[C.TypeTrojan] = location
			} else {
				openrestyConfig = strings.Replace(openrestyConfig, "TROJAN_TCP_PORT", strconv.Itoa(int(trojan.ListenPort)), 1)
			}
		case C.TypeVLESS:
			vless := inbound.VLESSOptions
			if locations[C.TypeVLESS] == "" {
				locations[C.TypeVLESS] = strings.Join(locationTemplace[:], "\n")
			}

			if vless.Transport != nil {
				location := locations[C.TypeVLESS]
				location = strings.Replace(location, "PATH", C.TypeVLESS, 1)

				switch vless.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					location = strings.Replace(location, "WS_PORT", strconv.Itoa(int(vless.ListenPort)), 1)
				case C.V2RayTransportTypeHTTPUpgrade:
					location = strings.Replace(location, "HU_PORT", strconv.Itoa(int(vless.ListenPort)), 1)
				case C.V2RayTransportTypeGRPC:
					location = strings.Replace(location, "GRPC_PORT", strconv.Itoa(int(vless.ListenPort)), 1)
				}

				locations[C.TypeVLESS] = location
			} else {
				// WIP
				// Do nothing but maybe someday will be used
				openrestyConfig = strings.Replace(openrestyConfig, "VLESS_TCP_PORT", strconv.Itoa(int(vless.ListenPort)), 1)
			}
		case C.TypeVMess:
			vmess := inbound.VMessOptions
			if locations[C.TypeVMess] == "" {
				locations[C.TypeVMess] = strings.Join(locationTemplace[:], "\n")
			}

			if vmess.Transport != nil {
				location := locations[C.TypeVMess]
				location = strings.Replace(location, "PATH", C.TypeVMess, 1)

				switch vmess.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					location = strings.Replace(location, "WS_PORT", strconv.Itoa(int(vmess.ListenPort)), 1)
				case C.V2RayTransportTypeHTTPUpgrade:
					location = strings.Replace(location, "HU_PORT", strconv.Itoa(int(vmess.ListenPort)), 1)
				case C.V2RayTransportTypeGRPC:
					location = strings.Replace(location, "GRPC_PORT", strconv.Itoa(int(vmess.ListenPort)), 1)
				}

				locations[C.TypeVMess] = location
			} else {
				openrestyConfig = strings.Replace(openrestyConfig, "VMESS_TCP_PORT", strconv.Itoa(int(vmess.ListenPort)), 1)
			}
		}
	}

	var ll []string
	for _, loc := range locations {
		ll = append(ll, loc)
	}

	loc := []string{
		`		location / {`,
		`			if ($http_upgrade = "websocket") {`,
		`				proxy_pass "http://127.0.0.1:` + strconv.Itoa(CS.WSTunnelPort) + `";`,
		`				break;`,
		`			}`,
		``,
		`			proxy_pass "http://127.0.0.1:` + strconv.Itoa(CS.WebServerPort) + `";`,
		`		}`,
	}
	ll = append(ll, strings.Join(loc[:], "\n"))

	openrestyConfig = strings.Replace(openrestyConfig, "DOMAIN", os.Getenv("DOMAIN"), -1)
	openrestyConfig = strings.Replace(openrestyConfig, "LOCATION_PLACEHOLDER", strings.Join(ll[:], "\n\n"), -1)
	openrestyConfig = strings.ReplaceAll(openrestyConfig, "ARG", "")

	f, err := os.Create("/etc/openresty/nginx.conf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(openrestyConfig)
}
