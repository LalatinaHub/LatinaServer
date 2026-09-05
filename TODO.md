# TODO: REFACTORING & INDEPENDENCE PROJECT

> **Tujuan Utama**: 
> 1. Modularisasi, restrukturisasi, dan optimalisasi performa kode
> 2. Memerdekakan project dari dependensi FoolVPN-ID (megalodon, megalodon-api, tool) demi kemudahan pengembangan dan decoupling sing-box version conflict.

**Status**: ✅ Phase 1-9 Complete & Production Ready (Leaf Module Successfully Extracted & Integrated)  
**Target Selesai**: Selesai (Fase 1-9)  
**Prioritas**: 🟢 Completed & Fully Verified


---

## 🎯 VISION: ARSITEKTUR AKHIR (Leaf Module Pattern)

### Ekosistem Setelah Refactor & Ekstraksi Selesai

```
┌─────────────────────────────────────────────────────────┐
│        latina-common (Future Leaf Module)               │
│  ┌───────────────────────────────────────────────────┐  │
│  │ • Zero heavy dependencies (hanya stdlib + libsql)│  │
│  │ • Entities: User, Server, KV, ProxyNode          │  │
│  │ • DB Client & Repository Pattern                  │  │
│  │ • Shared domain logic (validasi, helper)         │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────┬───────────────────────────────────┘
                      │ (go get)
       ┌──────────────┼──────────────┐
       │              │              │
       ▼              ▼              ▼
┌────────────┐ ┌────────────┐ ┌────────────┐
│ LatinaServer│ │ Telegram Bot│ │Web Dashboard│
│ + sing-box │ │ (pure CLI)  │ │ (Gin/Fiber)│
│ + caddy    │ │ No conflicts│ │ Lightweight│
└────────────┘ └────────────┘ └────────────┘
```

**Prinsip Utama**:
- **In-tree First, Extract Later**: Bangun kode bersih di dalam `LatinaServer` dulu (`internal/domain/`, `internal/database/`), pastikan stabil di production, baru diekstrak ke repo terpisah (`latina-common`).
- **Zero Heavy Deps di Leaf**: Modul database/entity TIDAK BOLEH mengimpor `sing-box`, `caddy`, atau `gin`.
- **Consumer-Specific Converter**: Mapping dari `ProxyNode` ke `sing-box.option.Outbound` berada di `LatinaServer`, bukan di library database.

---

## 📋 FASE 1: AUDIT & DEPENDENCY MAPPING (Estimasi: 2 Hari) ✅ COMPLETE

### 1.1 Mapping Dependensi FoolVPN-ID ✅
- [x] Audit menyeluruh struktur data & fungsi dari:
  - `github.com/FoolVPN-ID/megalodon-api/modules/db`
  - `github.com/FoolVPN-ID/megalodon-api/modules/db/servers`
  - `github.com/FoolVPN-ID/megalodon-api/modules/db/users`
  - `github.com/FoolVPN-ID/megalodon-api/modules/proxy`
  - `github.com/FoolVPN-ID/megalodon/db` (`ProxyFieldStruct`)
  - `github.com/FoolVPN-ID/tool/modules/subconverter`
- [x] Catat konflik versi dependency transitive antara FoolVPN-ID dan `sagernet/sing-box`.
- [x] **Hasil**: Dokumentasi lengkap di `docs/phase1-dependency-audit.md`

### 1.2 Rancangan Kontrak Baru (Interfaces & Structs) ✅
- [x] Definisikan domain models independen untuk:
  - [x] `User` → `internal/domain/model/user.go`
  - [x] `Server` → `internal/domain/model/server.go`
  - [x] `KeyValue` → `internal/domain/model/kv.go`
  - [x] `ProxyNode` → `internal/domain/model/proxy.go`
- [x] Definisikan interface repository & service layer → `internal/repository/interfaces.go`
- [x] Setup database client singleton → `internal/database/client.go`

---

## 📋 FASE 2: ELIMINASI DEPENDENSI FOOLVPN-ID (Estimasi: 5-7 Hari) ✅ COMPLETE

### 2.1 Database Independence (`db/` -> `internal/database/`) ✅
- [x] Buat native database client wrapper langsung menggunakan `tursodatabase/libsql-client-go` (atau standard `database/sql` driver).
- [x] Implementasikan internal model structs:
  - `internal/domain/model/user.go` (menggantikan `users.UserStruct`)
  - `internal/domain/model/server.go` (menggantikan `servers.ServerStruct`)
  - `internal/domain/model/kv.go`
- [x] Implementasikan Repository Pattern:
  - `internal/repository/user_repository.go`
  - `internal/repository/server_repository.go`
  - `internal/repository/kv_repository.go`
  - `internal/repository/proxy_repository.go`
- [x] Migrasikan logic `UpdateAndCheckPremiumQuota` dan query users tanpa ketergantungan paket luar.

### 2.2 Subconverter & Proxy Parsing Independence ✅
- [x] Buat native proxy parser & converter:
  - `internal/proxy/parser.go` (parsing protocol URL shadowsocks, trojan, vmess, vless)
  - `internal/proxy/converter_singbox.go` (mapping langsung ke `option.Outbound` sing-box)
  - Menggantikan ketergantungan `FoolVPN-ID/tool/modules/subconverter` dan `mgpr.ConvertDBToURL`.
- [x] Buat native relay proxy model internal (`internal/domain/model/proxy.go`)
- [x] Implementasikan fetcher query relay internal (`GatherRelays`) langsung ke database target.
- [x] Refactor `config/relay/relay.go` dan `db/db.go` sebagai wrapper ke repository baru.
- [x] Jalankan `go mod tidy` dan pastikan ketiga library `FoolVPN-ID/*` terhapus bersih dari `go.mod`.
- [x] Unit tests untuk parser (shadowsocks, vmess, vless, trojan)
- [x] Unit tests untuk converter (sing-box outbound generation)
- [x] Unit tests untuk domain models (User.IsActive)

---

## 📋 FASE 3: REORGANISASI & MODULARISASI STRUKTUR (Estimasi: 4-5 Hari) ✅ COMPLETE

### 3.1 Standard Go Project Layout ✅
Restrukturisasi codebase ke pola standar Go (`cmd/`, `internal/`, `pkg/`):

```
LatinaServer/
├── cmd/
│   └── latinaserver/
│       └── main.go                  # Bootstrap & Lifecycle orchestration
├── internal/
│   ├── config/                      # Config generator (Caddy & Sing-box)
│   │   ├── caddy.go
│   │   ├── sing.go
│   │   ├── util.go
│   │   ├── constants.go
│   │   └── relay/
│   │       └── relay.go
│   ├── domain/                      # Domain logic & entity models
│   │   └── model/
│   │       ├── user.go
│   │       ├── server.go
│   │       ├── kv.go
│   │       └── proxy.go
│   ├── infrastructure/              # External clients & integrations
│   │   ├── cloudflare/              # CF DNS client
│   │   │   └── client.go
│   │   ├── database/                # Turso connection pool
│   │   │   └── client.go
│   │   ├── geoip/                   # IP & location resolver
│   │   │   ├── geoip.go
│   │   │   └── countries.go
│   │   └── v2ray/                   # V2Ray gRPC client
│   │       └── stats.go
│   ├── proxy/                       # Proxy parser & converter
│   │   ├── parser.go
│   │   └── converter_singbox.go
│   ├── repository/                  # Data access layer
│   │   ├── interfaces.go
│   │   ├── user_repository.go
│   │   ├── server_repository.go
│   │   ├── kv_repository.go
│   │   └── proxy_repository.go
│   ├── service/                     # Service runners
│   │   ├── caddy/                   # Caddy runner & context management
│   │   │   └── runner.go
│   │   ├── singbox/                 # Sing-box runner & registration
│   │   │   └── runner.go
│   │   └── web/                     # Gin router, handlers, middlewares
│   │       ├── server.go
│   │       ├── runner.go
│   │       ├── health_handler.go
│   │       ├── status_handler.go
│   │       ├── proxy_handler.go
│   │       └── admin_handler.go
│   └── services/                    # Compat wrapper
│       └── compat.go
├── pkg/                             # Shared utility packages
│   ├── constant/                    # Centralized constants
│   │   └── constant.go
│   ├── systemctl/                   # Systemctl wrapper
│   │   └── systemctl.go
│   └── util/                        # Utility helpers
│       ├── slice.go
│       ├── proxyip.go
│       └── psutil.go
├── resources/                       # Static configs & assets
│   ├── caddy/
│   └── sing-box/
├── scripts/                         # Systemd & deployment scripts
├── cf/                              # Compat wrapper
│   └── compat.go
├── config/                          # Compat wrapper
│   ├── compat.go
│   └── relay/
│       └── compat.go
├── constant/                        # Compat wrapper
│   └── compat.go
├── db/                              # Compat wrapper
│   └── db.go
├── helper/                          # Compat wrapper
│   └── compat.go
└── web/                             # Compat wrapper
    └── web.go
```

### 3.2 Pemecahan Tanggung Jawab (Decoupling) ✅
- [x] Pindahkan konstanta dari `constant/` ke `pkg/constant/constant.go`.
- [x] Buat compat wrapper di `constant/compat.go` untuk backward compatibility.
- [x] Pecah Web Server (`web/web.go`) menjadi modular handlers di `internal/service/web/`:
  - `health_handler.go` (`/ping`, `/info`)
  - `status_handler.go` (`/status`)
  - `proxy_handler.go` (`/check`, `/relay`)
  - `admin_handler.go` (`/reload` via secure header token)
  - `server.go` (Router setup)
  - `runner.go` (Context runner)
- [x] Pindahkan `helper/` ke struktur yang lebih terorganisir:
  - `internal/infrastructure/geoip/` (geoip.go, countries.go)
  - `internal/infrastructure/v2ray/` (stats.go)
  - `pkg/systemctl/` (systemctl.go)
  - `pkg/util/` (slice.go, proxyip.go, psutil.go)
- [x] Pindahkan `cf/` ke `internal/infrastructure/cloudflare/`
- [x] Pindahkan `config/` ke `internal/config/`
- [x] Pindahkan `internal/services/` ke `internal/service/` dengan struktur folder terpisah
- [x] Update `cmd/latinaserver/main.go` menggunakan package baru
- [x] Buat backward compatibility wrappers di folder lama
- [x] Verifikasi build: `go build ./cmd/latinaserver` ✅
- [x] Verifikasi tests: `go test -v ./...` ✅ All tests pass


---

## 📋 FASE 4: OPTIMALISASI PERFORMA (Estimasi: 3-4 Hari) ✅ COMPLETE

### 4.1 Database Connection Pooling & Query Optimization ✅
- [x] Setup connection pool singleton dengan `sync.Once`:
  - MaxOpenConns (25), MaxIdleConns (5), ConnMaxLifetime (5m), ConnMaxIdleTime (1m).
  - Health check (`Ping`) dengan context timeout.
- [x] Implementasikan prepared statements untuk query repetitif (`DeductQuotaBatch`).
- [x] Tambahkan batch update quota (bulk deduction dalam satu database transaction dengan prepared statements).
- [x] Buat database indexes via `database.InitIndexes`:
  - `CREATE INDEX IF NOT EXISTS idx_users_expired ON users(expired);`
  - `CREATE INDEX IF NOT EXISTS idx_users_quota ON users(quota);`
  - `CREATE INDEX IF NOT EXISTS idx_users_server_code ON users(server_code);`
  - `CREATE INDEX IF NOT EXISTS idx_users_vpn ON users(vpn);`

### 4.2 Config Generation Caching & Parallel Execution ✅
- [x] Implementasikan content-based hash untuk config:
  - Hitung SHA256 dari content config JSON via `SaveJsonToFileWithCache`.
  - Cache config jika hash sama, skip rewriting disk.
- [x] Generate Caddy & Sing-box config secara paralel (`GenerateConfigsParallel` dengan goroutine & WaitGroup).
- [x] Tambahkan config diff logging saat hash berubah/tidak berubah.

### 4.3 Non-Blocking Relay Fetching ✅
- [x] Ubah `relay.GatherRelays()` menjadi non-blocking / background fetching:
  - Background periodic refresh setiap 15 menit (`StartBackgroundRelayFetcher`).
  - Exponential backoff pada failure (1s -> 2s -> 4s -> max 1m).
- [x] Implementasikan in-memory cache relay thread-safe (`sync.RWMutex`) dengan `GetRelays()`.
- [x] Fallback ke stale relay data jika fetch gagal tanpa mengosongkan cache.

### 4.4 Memory & Goroutine Management ✅
- [x] Rate-limit `runtime.FreeOSMemory()` (dibatasi max 1x per 5 menit, bukan setiap loop restart).
- [x] Audit goroutine launch dengan context cancellation dan graceful shutdown (`cancelServices()`, `cancel()`).
- [x] Gunakan `context.WithTimeout` untuk external calls (DB ping/init, quota batch deduction, relay fetching).
- [x] Setup pprof endpoints (`/debug/pprof/*`) untuk production profiling di Gin.

### 4.5 API Optimizations ✅
- [x] Cache `/api/v1/info` response (TTL 5 menit) menggunakan thread-safe `MemoryCache`.
- [x] Cache `/api/v1/status` response (TTL 30 detik) menggunakan thread-safe `MemoryCache`.
- [x] Implementasikan rate limiting per-IP (token bucket 100 req/min) via `RateLimiter` middleware.
- [x] Tambahkan gzip compression middleware (`GzipMiddleware`) dengan header `Vary: Accept-Encoding`.



---

## 📋 FASE 5: ERROR HANDLING & LOGGING (Estimasi: 2-3 Hari) ✅ COMPLETE

### 5.1 Replace Panic dengan Proper Error Handling ✅
- [x] Buat custom error types di `pkg/errors/`:
  - `DatabaseError`, `ConfigError`, `NetworkError`, `ValidationError`, `InternalError`.
- [x] Audit semua `panic()` dan ganti dengan error return + wrapping (`appErrors.New*` / `%w`).
- [x] Implementasikan error middleware di Gin (`ErrorMiddleware` untuk centralized error response).

### 5.2 Structured Logging (Zerolog) ✅
- [x] Setup di `pkg/logger/logger.go`:
  - Log level via env var (`LOG_LEVEL`).
  - JSON format untuk production.
  - Console pretty-print untuk non-production.
- [x] Migrate semua `fmt.Println`, `log.Println` ke structured logger.
- [x] Tambahkan contextual fields: `request_id`, `trace_id`, `client_ip` via `RequestLoggerMiddleware`.

---

## 📋 FASE 6: TESTING & QUALITY ASSURANCE (Estimasi: 3-4 Hari) ✅ COMPLETE

### 6.1 Unit Testing (Target: >70% Coverage di Core Packages) ✅
- [x] Setup `testify/assert`, `testify/mock`, dan `go-sqlmock`.
- [x] Tulis unit tests untuk:
  - `internal/repository/*_test.go` (mock database SQL queries via sqlmock: user, server, kv, proxy)
  - `internal/proxy/parser_test.go` (test protocol parsing SS, VMess, VLESS, Trojan + edge cases)
  - `internal/proxy/converter_singbox_test.go` (validate sing-box outbound converter)
  - `internal/config/*_test.go` (config generation logic & caching)
  - `internal/domain/model/*_test.go` (User.IsActive & Server.IsFull)
  - `internal/service/web/handlers_test.go` (Ping, Relays, Check proxy)
  - `pkg/util/slice_test.go` (RemoveDuplicate, CheckProxyIP)
- [x] Setup code coverage reporting (`go test -coverprofile="coverage.out" ./...`). Target >70% tercapai di critical paths (`domain`: 100%, `proxy`: 93.7%, `repository`: 83.3%, `errors`: 92.3%, `logger`: 77.8%).

### 6.2 Integration Testing & Linting ✅
- [x] Test database repository dengan real Turso connection (`internal/repository/integration_test.go` dengan `//go:build integration`).
- [x] Test Caddy & Sing-box config generation integration (`internal/config/integration_test.go` dengan `//go:build integration`).
- [x] Linting dan sanity checks via `go vet ./...`.

---

## 📋 FASE 7: SECURITY HARDENING (Estimasi: 1-2 Hari) ✅ COMPLETE

### 7.1 Credential & API Security ✅
- [x] Remove hardcoded credentials (`resources/sing-box/config.json` dibersihkan dari default static passwords/UUIDs; `internal/config/secret.go` membaca `SS_PASSWORD`, `CLASH_SECRET` atau generate cryptographically secure random token saat runtime).
- [x] Pindahkan API token dari URL path ke Authorization header (`Bearer <token>` dengan middleware `AuthMiddleware()` memakai constant-time comparison `subtle.ConstantTimeCompare`; route `/api/v1/admin/reload` aktif; URL path lama diberi deprecation notice).
- [x] Replace `grpc.WithInsecure()` dengan TLS / transport credentials support (`internal/infrastructure/v2ray/stats.go` mendukung `V2RAY_API_TLS`, `V2RAY_API_CERT`, dan `insecure.NewCredentials()` fallback).


---

### 🌿 FASE 9: EKSTRAKSI LEAF MODULE (`github.com/LalatinaHub/common`) (Estimasi: 1-2 Hari) ✅ COMPLETE
 
 > **Catatan**: Fase ini dieksekusi SETELAH Fase 2 dan 3 stabil di production.
 
 ### 9.1 Pemisahan Target Ekstraksi ✅
 - [x] **Komponen yang DI-EKSTRAK ke Leaf Module**:
   - ✅ `internal/domain/model/` (`User`, `Server`, `KV`, `ProxyNode`) -> `model/`
   - ✅ `internal/infrastructure/database/` (Turso LibSQL client & connection pool) -> `database/`
   - ✅ `internal/repository/` (CRUD SQL query interfaces & implementations) -> `repository/`
   - ✅ `pkg/proxy/parser.go` (Generic URL parser: SS, VMess, VLESS, Trojan) -> `proxy/`
 - [x] **Komponen yang TETAP di LatinaServer**:
   - ❌ `converter.go` (Mapping `ProxyNode` -> `sing-box.option.Outbound`)
   - ❌ Caddy & Sing-box runners & config builders
   - ❌ Gin web server & HTTP handlers
 
 ### 9.2 Pembuatan Repositori Mandiri (`github.com/LalatinaHub/common`) ✅
 - [x] Buat repo baru: `github.com/LalatinaHub/common` (Public repo di GitHub).
 - [x] Setup `go.mod` murni (hanya depend ke Go stdlib dan `libsql-client-go`).
 - [x] Pindahkan kode yang sudah terisolasi ke repo baru (`model/`, `database/`, `repository/`, `proxy/`).
 - [x] Tambahkan unit tests dan dokumentasi godoc (100% tests passing, high coverage).
 - [x] Tag versi rilis pertama `v0.1.0` dan push ke GitHub origin.
 
 ### 9.3 Integrasi Balik ke Konsumen ✅
 - [x] Update `go.mod` di `LatinaServer`:
   - `github.com/LalatinaHub/common v0.1.0`
   - `replace github.com/LalatinaHub/common => ../common` (untuk local development workflow)
 - [x] Ganti implementasi internal di `LatinaServer` ke paket baru dengan backward compatibility forwarders.
 - [x] Pakai modul yang sama untuk project lain (Telegram Bot, CLI, Dashboard) tanpa khawatir bentrok versi sing-box.

## 📋 FASE 8: DOCUMENTATION & CI/CD (Estimasi: 2 Hari) ✅ COMPLETE

### 8.1 Dokumentasi & CI/CD ✅
- [x] Tambahkan godoc comments untuk semua exported symbols.
- [x] Buat **README.md** (arsitektur, panduan instalasi, konfigurasi, API).
- [x] Buat **MIGRATION.md** (breaking changes & upgrade guide).
- [x] Update GitHub Actions: add unit test stage, lint stage, upload coverage.


---

## 🎯 SUCCESS CRITERIA (Updated dengan Leaf Module Goal)

### Phase 1-2: Independence Achievement
- ✅ `go list -m all | grep FoolVPN` returns **empty**
- ✅ All tests passing tanpa FoolVPN-ID dependencies
- ✅ Production deployment successful dengan kode native
- ✅ Sing-box dapat di-upgrade tanpa broken build

### Phase 3-8: Code Quality & Performance
- ✅ Clean architecture fully implemented
- ✅ >70% test coverage pada critical paths
- ✅ Zero `panic()` dalam production code paths
- ✅ Structured logging (zerolog/zap) di semua layer
- ✅ Performance benchmarks improved (startup, config gen, API latency)

### Phase 9: Leaf Module Extraction
- ✅ `github.com/LalatinaHub/common` repo public & documented
- ✅ LatinaServer successfully imports `github.com/LalatinaHub/common@v0.1.0`
- ✅ Telegram Bot/CLI tool dapat pakai `github.com/LalatinaHub/common` tanpa conflict
- ✅ Future projects: `go get github.com/LalatinaHub/common` langsung akses shared DB logic
- ✅ Decoupled: LatinaServer upgrade sing-box v2.0, project lain tidak terpengaruh

---

## 📊 PROGRESS TRACKING

### Milestone Summary
- **Week 1**: Fase 1-2 (Audit, Isolasi DB & Parser Native)
- **Week 2**: Fase 3 (Restrukturisasi Folder)
- **Week 3**: Fase 4-6 (Optimasi, Error Handling, Testing)
- **Week 4**: Fase 7-8 (Security, Documentation)
- **Week 4+**: Fase 9 (Ekstraksi Leaf Module) ✅

### Risk Mitigation Strategy
- **Risk**: Breaking production saat migration
  - **Mitigation**: Feature branch `refactor/independence`, incremental PR per fase, staging deployment sebelum production
- **Risk**: Relay fetching gagal
  - **Mitigation**: Fallback ke cache lama, monitoring & alerting
- **Risk**: Performance regression
  - **Mitigation**: Benchmark test sebelum/sesudah, load testing 1000+ users

---

**Dibuat**: 2026-09-05  
**Versi**: 3.0 (Fase 1-9 Complete & Production Ready)  
**Next Action**: 🎉 **All Phases 1 to 9 Complete!** System stabilized and leaf module extracted.
