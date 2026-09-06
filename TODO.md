# LatinaServer — Product Roadmap & Development TODO

Dokumen ini memuat roadmap pengembangan fitur unggulan untuk meningkatkan daya tarik produk, kemudahan penggunaan (*user experience*), dan nilai jual (*unique selling proposition*) bagi pelanggan **LatinaServer**.

---

## 🎯 Prioritas Pengembangan (Roadmap)

```
[Phase 1] [DONE] Smart Subscription Engine (P1) — Selesai di ../LatinaApi
   └── User-Agent Auto Detect (Clash, sing-box, v2rayNG, Shadowrocket)
[Phase 2] [DONE] Streaming & AI Unlocker Routing (P2)
   └── Built-in WARP/WireGuard Outbound + Smart Geosite Routing
[Phase 3] [DONE] AdBlock & Malware Shield Per-User (P3)
   └── Sinkronisasi field User.Adblock ke sing-box Rule-Sets
[Phase 4] Turbo Mode — TCP Brutal Acceleration (P4)
   └── Anti-Packet Loss Profiles untuk Gaming & Streaming Jam Sibuk
[Phase 5] [DONE] Customer Self-Service Portal (P5)
   └── Web Portal per-token: Monitor Kuota, Reset UUID, 1-Click Import, Parity dengan LatinaBot
[Phase 6] Multi-Hop Ghost Routing & Relay Failover (P6)
   └── Automated Relay Health-Check & Circuit Breaker
```

---

## 📋 Detail Spesifikasi Tiap Fase

### 🚀 Phase 1: Universal Smart Subscription Engine (`/sub`)
> **Status**: ✅ **Selesai (Completed)** — Diimplementasikan pada microservice terpisah di `../LatinaApi`.
> **Target**: Menghilangkan hambatan setup bagi pengguna awam di berbagai platform (Android, iOS, Windows, macOS, Linux).

- [x] **Rute Endpoint**:
  - Implementasikan handler `GET /sub` dan `GET /api/v1/sub`.
  - Menerima parameter query token: `/sub?token={USER_TOKEN}`.
  - Tambahkan rute Caddy reverse proxy untuk `/sub` ke Gin `WEBSERVER_PORT`.
- [x] **Sistem Deteksi User-Agent Dinamis**:
  - `Clash.Meta` / `mihomo` / `ClashVerge` / `Flclash` $\rightarrow$ Kirim format **Clash Meta YAML** lengkap dengan *Proxy Groups* (Auto Fallback, Load Balance, Direct Indo) dan *Rule Providers*.
  - `sing-box` / `SFI` / `SFM` $\rightarrow$ Kirim format **sing-box v1.14+ JSON** lengkap dengan *inbounds*, *outbounds*, *route.rule_set*, dan *DNS*.
  - `v2rayNG` / `Shadowrocket` / `v2rayN` / `NekoBox` $\rightarrow$ Kirim format **Base64-encoded lines** (`vmess://...`, `vless://...`, `trojan://...`).
- [x] **Standard Subscription Headers**:
  - Sertakan header `Subscription-Userinfo: upload={bytes}; download={bytes}; total={quota_bytes}; expire={unix_timestamp}` agar aplikasi client langsung menampilkan sisa kuota dan masa aktif di UI mereka.
  - Header `Profile-Update-Interval: 24` (dalam jam).
  - Header `Content-Disposition: attachment; filename="{profile_name}"`.
- [x] **Validasi Token & Keamanan**:
  - Verifikasi token ke basis data: pastikan akun belum expired dan kuota belum habis.
  - Kirim HTTP 403 Forbidden / HTTP 404 jika token tidak valid atau expired.

---

### 🌐 Phase 2: Built-in Streaming & AI Unlocker (WARP Smart Routing)
> **Status**: ✅ **Selesai (Completed)** — Built-in Cloudflare WARP WireGuard endpoint, Smart Domain Routing, dan Health-Checker dengan graceful fallback.
> **Target**: Membuka blokir IP datacenter untuk ChatGPT, Claude, Netflix, Disney+, Reddit, dan bank lokal tanpa mengurangi kecepatan browsing umum.

- [x] **Outbound Cloudflare WARP / WireGuard**:
  - Tambahkan template outbound `wireguard` lokal di `resources/sing-box/config.json`.
  - Dukungan auto-register / generate WARP private key & peer endpoint via Go background service.
- [x] **Smart Rule Routing di sing-box v1.14+**:
  - Konfigurasikan `route.rules`:
    - Geosite AI (`category-ai-chat-!cn`, `openai`, `anthropic`) $\rightarrow$ diarahkan ke outbound `warp-out`.
    - Geosite Media Streaming (`netflix`, `disney`, `primevideo`, `hulu`) $\rightarrow$ diarahkan ke outbound `warp-out`.
    - Traffic umum & direct Indonesia $\rightarrow$ tetap melalui direct IP server untuk latensi terendah dan kecepatan maksimal.
- [x] **WARP Health-Checker Goroutine**:
  - Background goroutine di Go untuk ping/curl berkala ke endpoint WARP; otomatis bypass ke direct jika tunnel WARP mengalami kendala (*graceful fallback*).

---

### 🛡️ Phase 3: AdBlock & Malware Shield Per-User
> **Status**: ✅ **Selesai (Completed)** — Dual-layer DNS (AdGuard DNS) & Route Rule Reject per-user, database integration, dan instant toggle REST API (`/api/v1/user/settings`).
> **Target**: Menghemat kuota hingga 30% dan melindungi pengguna dari iklan spam, situs judi online, dan malware.

- [x] **Integrasi Field Basis Data**:
  - Manfaatkan field `Adblock bool` yang telah tersedia di [`model.User`](internal/domain/model/user.go).
- [x] **Injeksi DNS & Routing Rules di sing-box**:
  - Jika akun memiliki `Adblock: true`, asosiasikan user tersebut ke routing rule yang menolak domain ads, tracker, malware, dan judi online (`action: reject`).
  - Arahkan query DNS user tersebut ke upstream secure AdGuard DNS (`94.140.14.14:53`).
- [x] **API Pengaturan User**:
  - Endpoint `GET /api/v1/user/settings` dan `POST /api/v1/user/settings` untuk toggle fitur AdBlock secara instan dari dashboard/portal pelanggan.

---

### ⚡ Phase 4: "Turbo Mode" — Akselerasi Anti-Packet Loss (TCP Brutal)
> **Target**: Menjamin streaming 4K dan unduhan stabil tanpa patah-patah di jam sibuk ISP seluler yang sering mengalami packet loss tinggi.

- [ ] **Tiering Bandwidth Brutal**:
  - Profil Standard: TCP / BBR reguler (efisien bandwidth server).
  - Profil Turbo / Gaming: Outbound sing-box dengan parameter `multiplex.brutal`:
    ```json
    "brutal": {
      "enabled": true,
      "up_mbps": 25,
      "down_mbps": 50
    }
    ```
- [ ] **Diferensiasi Kuota/Paket**:
  - Pilihan aktivasi Turbo Mode langsung disuntikkan ke konfigurasi subscription client sing-box.

---

### 💻 Phase 5: Customer Self-Service Portal (`/portal`)
> **Status**: ✅ **Selesai (Completed)** — Single-Page App berbasis Vue 3 (Composition API) + Tailwind CSS + Lucide Icons + QRCode.js dengan fitur setara penuh `../LatinaBot`, reverse proxy Caddy ke webserver, dan REST API lengkap.
> **Target**: Portal mandiri profesional bagi pelanggan untuk memantau status langganan dan mengelola akun mereka secara mandiri.

- [x] **Halaman Web `/portal?token={USER_TOKEN}` (Vue 3 Single Page Application)**:
  - Visualisasi kuota: Terpakai vs Total kuota (Progress bar interaktif).
  - Status masa aktif & countdown hari tersisa ("28 hari lagi"), status donator badge.
  - Switcher protokol interaktif (VMess, VLESS, Trojan).
  - Pemilihan Edge Server dengan pemantauan load pengguna (`users_count / users_max`).
  - Pemilihan jalur relay negara multi-hop (ID, SG, JP, US, atau Tanpa Relay).
- [x] **1-Click Import Actions & Mobile QR Scanner**:
  - Tombol instan yang membuka aplikasi dengan scheme URL:
    - *Import to Clash* (`clash://install-config?url=...`)
    - *Import to Shadowrocket* (`shadowrocket://add/sub://...`)
    - *Import to Sing-box* (`sing-box://import-remote-profile?url=...`)
    - *Import to v2rayNG / NekoBox*
  - Tampilan QR Code subscription & raw node untuk scan kamera HP.
- [x] **Fitur Reset Kredensial & Keamanan Mandiri**:
  - Tombol **"Reset UUID / Password"**: Otomatis generate UUID baru dan reload sing-box untuk memutus perangkat tak dikenal.
  - Tombol **"Ganti Token / Password"**: Generate token 8-karakter baru.
- [x] **AdBlock & Malware Shield Toggle**:
  - 1-click toggle perlindungan AdBlock per-user dengan auto-reload sing-box.
- [x] **Hardware Telemetry & Live Stats**:
  - Pemantauan CPU%, RAM%, Disk%, Uptime sistem, dan Network I/O.
- [x] **Wildcard Domains Guide & Bug Host**:
  - Daftar domain wildcard beserta panduan format SNI (`bug.domain.com`).
- [x] **Donasi & Trust Policy**:
  - Trakteer support link, panduan donasi 29-hari premium, dan jaminan Zero-Log.
- [x] **Modifikasi Caddy Layer-4/HTTP Reverse Proxy**:
  - Meneruskan request HTTP normal (termasuk `/portal`, `/portal/*`, static assets) langsung ke `localhost:WEBSERVER_PORT`.

---

### 🔄 Phase 6: Multi-Hop Ghost Routing & Relay Circuit-Breaker
> **Target**: Keamanan tingkat lanjut untuk pengguna privasi tinggi (fintech/trader) dengan rute relay dinamis yang tahan gangguan.

- [ ] **Relay Ping & Latency Benchmark**:
  - Optimasi worker pool di [`internal/proxy`](internal/proxy) untuk menguji relay secara non-blocking dengan context timeout.
- [ ] **Automated Circuit Breaker**:
  - Jika relay node tujuan mengalami timeout > 3x berturut-turut, otomatis ganti rute ke relay cadangan tanpa memutus sesi client.
- [ ] **Multi-Hop Tunneling Option**:
  - Inbound Client $\rightarrow$ LatinaServer (Node A) $\rightarrow$ Relay Luar Negeri (Node B) $\rightarrow$ Internet.

---

## 🛠️ Standar Kualitas Kode (Go Best Practices)

- **Concurrency Safety**: Selalu gunakan `context.Context` dengan timeout untuk setiap request eksternal / network call.
- **Graceful Shutdown**: Pastikan semua background worker/goroutine bersih saat service dimatikan via systemd.
- **Unit & Mock Testing**: Setiap handler dan repository baru wajib dilengkapi unit test dengan mock interface (`testify/mock`).
- **Zero-Allocation Logging**: Gunakan logger terstruktur `zerolog` yang telah terpasang di [`pkg/logger`](pkg/logger).
