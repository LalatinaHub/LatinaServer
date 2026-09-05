# FASE 1: DEPENDENCY AUDIT REPORT

**Tanggal**: 2026-09-05  
**Status**: ✅ Complete  
**Tujuan**: Mapping lengkap dependensi FoolVPN-ID dan analisis konflik transitive

---

## 1. DEPENDENCY MAPPING

### 1.1 Direct Dependencies FoolVPN-ID

```
LatinaServer directly imports:
├── github.com/FoolVPN-ID/megalodon@v0.0.0-20250629031537-324538133063
├── github.com/FoolVPN-ID/megalodon-api@v0.0.0-20250409180851-1a83b842cd39
└── github.com/FoolVPN-ID/tool@v0.0.0-20250629030518-3bf65b314abf
```

### 1.2 Usage Analysis per File

#### 📁 `db/db.go`
**Imports:**
- `database "github.com/FoolVPN-ID/megalodon-api/modules/db"`
- `"github.com/FoolVPN-ID/megalodon-api/modules/db/servers"`
- `"github.com/FoolVPN-ID/megalodon-api/modules/db/users"`

**Functions & Types Used:**
- `database.MakeDatabase().GetClient()` -> returns `*sql.DB` (Turso LibSQL client)
- `servers.ServerStruct` (fields: ID, Code, Domain, IP, Country, UsersCount, UsersMax)
- `users.UserStruct` (fields: ID, Token, Password, Expired, ServerCode, Quota, Relay, Adblock, VPN)

**Dependent Functions:**
1. `GetKVList() map[string]any`
2. `GetServerList() []servers.ServerStruct`
3. `GetPremiumList() map[string][]users.UserStruct`
4. `UpdateAndCheckPremiumQuota(name string) bool`

#### 📁 `config/relay/relay.go`
**Imports:**
- `database "github.com/FoolVPN-ID/megalodon-api/modules/db"`
- `mgpr "github.com/FoolVPN-ID/megalodon-api/modules/proxy"`
- `mgdb "github.com/FoolVPN-ID/megalodon/db"`
- `"github.com/FoolVPN-ID/tool/modules/subconverter"`

**Functions & Types Used:**
- `mgdb.ProxyFieldStruct` (25+ fields proxy schema)
- `mgpr.ConvertDBToURL(&proxy)` (converts proxy field struct to URL string)
- `subconverter.MakeSubconverterFromConfig(url)` (URL string -> sing-box outbounds)

**Dependent Functions:**
1. `GatherRelays()`
2. `GetRelayOutbounds() []option.Outbound`

---

## 2. TRANSITIVE DEPENDENCY & VERSION CONFLICT

### 2.1 Sing-box Version Lock
```
LatinaServer directly requires:
  github.com/sagernet/sing-box -> v1.12.19


---

## 3. RANCANGAN KONTRAK BARU & DOMAIN INTERFACES (Phase 1.2)

### 3.1 Domain Models (`internal/domain/model/`)

```go
package model

import "time"

type User struct {
	ID         int64
	Token      string
	Password   string
	Expired    time.Time
	ServerCode string
	Quota      int64 // bytes
	Relay      string
	Adblock    bool
	VPN        string
}

func (u *User) IsActive(now time.Time) bool {
	return u.Quota > 0 && !now.After(u.Expired) && u.ServerCode != "" && u.VPN != ""
}

type Server struct {
	ID         int64
	Code       string
	Domain     string
	IP         string
	Country    string
	UsersCount int64
	UsersMax   int64
}

type KeyValue struct {
	ID    int64
	Key   string
	Value any
}

type ProxyNode struct {
	ID          int64
	Server      string
	IP          string
	ServerPort  int
	UUID        string
	Password    string
	Security    string
	AlterID     int
	Method      string
	Plugin      string
	PluginOpts  string
	Host        string
	TLS         bool
	Transport   string
	Path        string
	ServiceName string
	Insecure    bool
	SNI         string
	Remark      string
	ConnMode    string
	CountryCode string
	Region      string
	Org         string
	VPN         string
	Raw         string
}
```

### 3.2 Repository Interfaces (`internal/repository/`)

```go
package repository

import (
	"context"
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
)

type UserRepository interface {
	GetActiveUsersGroupedByVPN(ctx context.Context) (map[string][]model.User, error)
	DeductQuota(ctx context.Context, userID int64, usedBytes int64) (remaining int64, isDepleted bool, err error)
}

type ServerRepository interface {
	GetAll(ctx context.Context) ([]model.Server, error)
}

type KVRepository interface {
	GetAll(ctx context.Context) (map[string]any, error)
}

type ProxyRepository interface {
	GetRelays(ctx context.Context, excludeCountryCodes []string, maxPerCountry int) ([]model.ProxyNode, error)
}
```

### 3.3 Proxy Converter Interface (Consumer-Specific)

```go
package proxy

import (
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/sagernet/sing-box/option"
)

type SingboxConverter interface {


---

## 4. REPLACEMENT STRATEGY

### 4.1 Database Layer (`db/db.go` -> `internal/database/`)

**Current Flow:**
```
db.go -> FoolVPN/megalodon-api -> MakeDatabase() -> sql.DB
```

**Target Flow:**
```
internal/database/client.go -> tursodatabase/libsql directly -> sql.DB (singleton pool)
internal/repository/user_repo.go -> Use client, return model.User (native struct)
```

**Steps:**
1. Create `internal/database/client.go` with connection pool singleton
2. Create `internal/repository/user_repo.go`, `server_repo.go`, `kv_repo.go`
3. Implement CRUD methods using prepared statements
4. Add backward-compatible adapter in `db/db.go` to delegate to new repository
5. Test equivalence, migrate callers, remove old code

### 4.2 Proxy Parser Layer (`config/relay/` -> `internal/proxy/`)

**Current Flow:**
```
ProxyFieldStruct -> mgpr.ConvertDBToURL() -> "ss://..." -> subconverter.MakeSubconverter() -> option.Outbound
```

**Target Flow:**
```
model.ProxyNode -> native parser/converter -> option.Outbound (direct, zero external deps)
```

**Steps:**
1. Create `internal/proxy/parser.go` with protocol-specific parsers (SS, VMess, VLESS, Trojan)
2. Create `internal/proxy/converter_singbox.go` to map ProxyNode -> option.Outbound
3. Update `config/relay/relay.go` to use native converter
4. Test output equivalence with existing relay outbounds

---

## 5. DEPENDENCY GRAPH COMPARISON

### Before (Current)
```
LatinaServer
├── FoolVPN-ID/megalodon-api
│   ├── sing-box@v1.12.0-beta ❌
│   └── 50+ transitive deps
├── FoolVPN-ID/megalodon
│   ├── sing-box@v1.12.0-beta ❌
│   └── echotron, azuretls (unused!)
├── FoolVPN-ID/tool
│   └── sing-box@v1.12.0-beta ❌
└── sing-box@v1.12.19 (blocked!) 🔒
```

### After (Phase 2 Complete)
```
LatinaServer
├── tursodatabase/libsql-client-go ✅
└── sagernet/sing-box@v1.12.19+ ✅ (free to upgrade!)
```

---

## 6. SUCCESS CRITERIA

- [x] **Phase 1.1**: Mapping selesai (current file)
- [x] **Phase 1.2**: Interface design complete (section 3)
- [ ] **Phase 2 Target**:
  - [ ] `go list -m all | grep FoolVPN` returns empty
  - [ ] All existing tests pass
  - [ ] Config output identical (byte-for-byte comparison)
  - [ ] Production deploy successful
  - [ ] Sing-box upgradable to v1.13+ without conflicts

---

## 7. NEXT ACTIONS

### Immediate (Day 1-3)
1. Create folder structure: `internal/domain/model/`, `internal/database/`, `internal/repository/`, `internal/proxy/`
2. Implement entity models (User, Server, KV, ProxyNode)
3. Implement database client singleton with connection pooling
4. Implement UserRepository with prepared statements

### Short-term (Day 4-7)
5. Implement ProxyRepository for relay fetching
6. Implement native proxy parser (SS, VMess, VLESS, Trojan)
7. Implement sing-box converter
8. Write unit tests with sqlmock and table-driven tests
9. Integration test dengan real Turso connection

### Final (Day 8-10)
10. Migrate `db/db.go` to use new repositories
11. Migrate `config/relay/relay.go` to use native parser
12. Run full integration test
13. Remove FoolVPN-ID imports, run `go mod tidy`
14. Deploy to staging, smoke test, then production

---

**Generated**: 2026-09-05  
**Author**: Kiro AI + golang-backend-development skill  
**Status**: Phase 1 Complete ✅ Ready for Phase 2 Implementation

	ToOutbound(node model.ProxyNode) (option.Outbound, error)
	BuildURLTest(countryCode string, tags []string) option.Outbound
}
```

FoolVPN-ID libraries force:
  github.com/sagernet/sing-box -> v1.12.0-beta.28.0.20250623124837-e5c6d3c080ca
```

**Dampak:**
- LatinaServer terkunci versi transitive sing-box beta dari FoolVPN-ID.
- Upgrade versi sing-box baru terblokir oleh transitive dependency tree.
- Modul DB dan parsing membawa ratusan dependensi yang tidak relevan (Echotron Telegram bot, AzureTLS, Caddy modules).
