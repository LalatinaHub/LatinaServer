module github.com/LalatinaHub/LatinaServer

go 1.22.7

toolchain go1.23.2

require (
	github.com/LalatinaHub/LatinaApi v0.0.0-20231222030803-2f2f646bc586
	github.com/LalatinaHub/LatinaSub-go v0.0.0-20240702171100-97fdb5e1708c
	github.com/LalatinaHub/wstunnel v0.0.0-20231112003007-4c4b149076c8
	github.com/dickymuliafiqri/BenchBox v0.0.0-20240313095405-9efff46826a8
	github.com/gin-gonic/gin v1.10.0
	github.com/go-co-op/gocron v1.37.0
	github.com/nedpals/supabase-go v0.4.0
	github.com/sagernet/sing v0.6.0-alpha.17
	github.com/sagernet/sing-box v1.10.1
	github.com/v2fly/v2ray-core/v5 v5.22.0
	golang.org/x/sync v0.9.0
	google.golang.org/grpc v1.68.0
)

replace github.com/sagernet/sing-box => github.com/LalatinaHub/sing-box v0.0.0-20241118142254-348a76bcd618

// replace github.com/sagernet/sing-box => ../sing-box

require (
	berty.tech/go-libtor v1.0.385 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/bytedance/sonic v1.12.4 // indirect
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/caddyserver/certmagic v0.21.4 // indirect
	github.com/caddyserver/zerossl v0.1.3 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/cloudflare/circl v1.5.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cretz/bine v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
	github.com/gaukas/godicttls v0.0.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.23.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20241101162523-b92577c0c142 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/insomniacslk/dhcp v0.0.0-20240829085014-a3a4c1f04475 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/libdns/alidns v1.0.3 // indirect
	github.com/libdns/cloudflare v0.1.1 // indirect
	github.com/libdns/libdns v0.2.2 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/metacubex/tfo-go v0.0.0-20241006021335-daedaf0ca7aa // indirect
	github.com/mholt/acmez v1.2.0 // indirect
	github.com/mholt/acmez/v2 v2.0.3 // indirect
	github.com/miekg/dns v1.1.62 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nedpals/postgrest-go v0.1.3 // indirect
	github.com/onsi/ginkgo/v2 v2.21.0 // indirect
	github.com/ooni/go-libtor v1.1.8 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pires/go-proxyproto v0.8.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/qtls-go1-20 v0.4.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sagernet/bbolt v0.0.0-20231014093535-ea5cb2fe9f0a // indirect
	github.com/sagernet/cloudflare-tls v0.0.0-20231208171750-a4483c1b7cd1 // indirect
	github.com/sagernet/fswatch v0.1.1 // indirect
	github.com/sagernet/gvisor v0.0.0-20241021032506-a4324256e4a3 // indirect
	github.com/sagernet/netlink v0.0.0-20240916134442-83396419aa8b // indirect
	github.com/sagernet/nftables v0.3.0-beta.4 // indirect
	github.com/sagernet/quic-go v0.48.1-beta.1 // indirect
	github.com/sagernet/reality v0.0.0-20230406110435-ee17307e7691 // indirect
	github.com/sagernet/sing-dns v0.4.0-alpha.3 // indirect
	github.com/sagernet/sing-mux v0.3.0-alpha.2 // indirect
	github.com/sagernet/sing-quic v0.4.0-alpha.4 // indirect
	github.com/sagernet/sing-shadowsocks v0.2.7 // indirect
	github.com/sagernet/sing-shadowsocks2 v0.2.0 // indirect
	github.com/sagernet/sing-shadowtls v0.2.0-alpha.2 // indirect
	github.com/sagernet/sing-tun v0.6.0-alpha.9 // indirect
	github.com/sagernet/sing-vmess v0.1.12 // indirect
	github.com/sagernet/smux v0.0.0-20231208180855-7041f6ea79e7 // indirect
	github.com/sagernet/tfo-go v0.0.0-20231209031829-7b5343ac1dc6 // indirect
	github.com/sagernet/utls v1.6.7 // indirect
	github.com/sagernet/websocket v0.0.0-20220913015213-615516348b4e // indirect
	github.com/sagernet/wireguard-go v0.0.0-20231215174105-89dec3b2f3e8 // indirect
	github.com/scjalliance/comshim v0.0.0-20240712181150-e070933cb68e // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba // indirect
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241113202542-65e8d215514f // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.3.0 // indirect
)
