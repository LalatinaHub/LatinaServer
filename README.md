# LatinaServer

LatinaServer is a high-performance, unified VPN edge server gateway built with Go, integrating [sing-box](https://github.com/sagernet/sing-box), [Caddy](https://github.com/caddyserver/caddy), and [Turso LibSQL](https://turso.tech).

It automates TLS termination, user credential provisioning, dynamic proxy relay routing, and real-time quota deduction via V2Ray gRPC stats.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                      Client Traffic                     │
│         (TLS / WebSocket / gRPC / HTTPUpgrade / TCP)    │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                   Caddy L4 / Reverse Proxy              │
│      (SNI Demux, TLS Termination, HTTP / REST API)      │
└──────────────┬───────────────────────────┬──────────────┘
               │                           │
   [Proxy Traffic Loopback]     [Admin / REST API]
               │                           │
               ▼                           ▼
┌────────────────────────────┐  ┌─────────────────────────┐
│      sing-box Core         │  │     Gin Web Service     │
│ - Inbounds: SS, VMess,     │  │ - /api/v1/ping          │
│   VLESS, Trojan            │  │ - /api/v1/status        │
│ - Dynamic Relay Outbounds  │  │ - /api/v1/check         │
│ - V2Ray API Stats Server   │  │ - /api/v1/relay         │
└──────────────┬─────────────┘  │ - /api/v1/admin/reload  │
               │ (gRPC Stats)   └───────────┬─────────────┘
               ▼                            │
┌────────────────────────────┐              │
│    Background Daemon       │              │
│ - Quota deduction sync     │              │
│ - Dynamic relay fetcher    │              │
└──────────────┬─────────────┘              │
               │                            │
               ▼                            ▼
┌─────────────────────────────────────────────────────────┐
│              Turso Database (LibSQL SQLite)             │
│            Users, Servers, Proxies, Key-Values          │
└─────────────────────────────────────────────────────────┘
```

---

## Features

- **Decoupled & Native**: 100% independent from proprietary legacy wrappers. Pure Go standard library + official drivers.
- **Multi-Protocol Support**: Shadowsocks, VMess (WS/GRPC/TCP), VLESS (WS/GRPC/TCP), Trojan (WS/GRPC/TCP).
- **Dynamic Relay Routing**: Background worker queries verified relays and routes user traffic through designated outbound nodes.
- **Real-time Quota Deduction**: Automated V2Ray gRPC downlink stats collection with batch database updates in single atomic transactions.
- **Config Hashing & Caching**: Avoids unnecessary disk writes when configuration has not changed.
- **Security Hardened**:
  - No hardcoded secrets or static tokens.
  - Constant-time token verification for administrative endpoints.
  - Rate limiting, Gzip compression, and structured contextual logging.

---

## Getting Started

### Prerequisites

- **Go**: 1.24+
- **Linux** (Production target with systemd) or Windows (Development)

### Configuration (Environment Variables)

| Variable | Description | Default / Example |
| :--- | :--- | :--- |
| `TURSO_DATABASE_URL` | Turso LibSQL database connection URL (**required**) | `libsql://your-db.turso.io?authToken=...` |
| `ADMIN_API_TOKEN` | Bearer token for protected `/api/v1/admin/*` routes | Random / configured string |
| `SS_PASSWORD` | Shadowsocks inbound/outbound master password | Auto-generated 32-hex chars |
| `CLASH_SECRET` | Secret token for Clash API dashboard | Optional |
| `LOG_LEVEL` | Application logging verbosity (`debug`, `info`, `warn`, `error`) | `info` |
| `SERVER_ENV` | Environment mode (`production`, `development`) | `production` |
| `V2RAY_API_TLS` | Enable TLS for internal V2Ray gRPC stats connection | `false` |
| `V2RAY_API_CERT` | Custom TLS certificate file for V2Ray stats | Optional |

---

## API Endpoints

### Public Endpoints

- `GET /api/v1/ping`: Health check probe, returns `Pong`.
- `GET /api/v1/info`: Returns server IP and GeoIP information (cached 5 min).
- `GET /api/v1/status`: Returns system CPU, memory, and uptime metrics (cached 30 sec).
- `GET /api/v1/check?ip=<proxy_ip>`: Verifies proxy IP reachability and latency.
- `GET /api/v1/relay`: Returns list of currently active outbound proxy relays.

### Admin Endpoints

- `POST /api/v1/admin/reload`: Triggers graceful configuration reload and service restart.
  - **Header Required**: `Authorization: Bearer <ADMIN_API_TOKEN>`
  - Response: `{"status":"success","message":"Service reload triggered"}`

---

## Build & Run

### Development

```bash
# Run tests
go test -v ./...

# Run with test coverage
go test -coverprofile="coverage.out" ./...
go tool cover -func="coverage.out"

# Vet and lint
go vet ./...
```

### Production Build

```bash
go build -tags with_grpc,with_clash_api,with_v2ray_api,with_utls,with_gvisor,with_quic -o latinaserver ./cmd/latinaserver
```

---

## Testing & Quality

- **Domain Models**: 100% code coverage.
- **Proxy Parser & Converters**: >93% coverage.
- **Data Repositories**: >83% coverage using `DATA-DOG/go-sqlmock`.
- **Integration Tests**: Supported with `//go:build integration` tags for Turso and live config builders.
