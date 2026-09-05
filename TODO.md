# LatinaServer — Product Roadmap & Development TODO

Dokumen ini memuat roadmap pengembangan fitur unggulan untuk meningkatkan daya tarik produk, kemudahan penggunaan (*user experience*), dan nilai jual (*unique selling proposition*) bagi pelanggan **LatinaServer**.

---

## 🎯 Prioritas Pengembangan (Roadmap)

```
[Phase 1] Smart Subscription Engine (P1)
   └── User-Agent Auto Detect (Clash, sing-box, v2rayNG, Shadowrocket)
[Phase 2] Streaming & AI Unlocker Routing (P2)
   └── Built-in WARP/WireGuard Outbound + Smart Geosite Routing
[Phase 3] AdBlock & Malware Shield Per-User (P3)
   └── Sinkronisasi field User.Adblock ke sing-box Rule-Sets
[Phase 4] Turbo Mode — TCP Brutal Acceleration (P4)
   └── Anti-Packet Loss Profiles untuk Gaming & Streaming Jam Sibuk
[Phase 5] Customer Self-Service Portal (P5)
   └── Web Portal per-token: Monitor Kuota, Reset UUID, 1-Click Import
[Phase 6] Multi-Hop Ghost Routing & Relay Failover (P6)
   └── Automated Relay Health-Check & Circuit Breaker
```

---

## 📋 Detail Spesifikasi Tiap Fase

### 🚀 Phase 1: Universal Smart Subscription Engine (`/sub`)
> **Target**: Menghilangkan hambatan setup bagi pengguna awam di berbagai platform (Android, iOS, Windows, macOS, Linux).

- [ ] **Rute Endpoint**:
  - Implementasikan handler `GET /sub` dan `GET /api/v1/sub`.
  - Menerima parameter query token: `/sub?token={USER_TOKEN}`.
  - Tambahkan rute Caddy reverse proxy untuk `/sub` ke Gin `WEBSERVER_PORT`.
- [ ] **Sistem Deteksi User-Agent Dinamis**:
  - `Clash.Meta` / `mihomo` / `ClashVerge` / `Flclash` $\rightarrow$ Kirim format **Clash Meta YAML** lengkap dengan *Proxy Groups* (Auto Fallback, Load Balance, Direct Indo) dan *Rule Providers*.
  - `sing-box` / `SFI` / `SFM` $\rightarrow$ Kirim format **sing-box v1.14+ JSON** lengkap dengan *inbounds*, *outbounds*, *route.rule_set*, dan *DNS*.
  - `v2rayNG` / `Shadowrocket` / `v2rayN` / `NekoBox` $\rightarrow$ Kirim format **Base64-encoded lines** (`vmess://...`, `vless://...`, `trojan://...`).
- [ ] **Standard Subscription Headers**:
  - Sertakan header `Subscription-Userinfo: upload={bytes}; download={bytes}; total={quota_bytes}; expire={unix_timestamp}` agar aplikasi client langsung menampilkan sisa kuota dan masa aktif di UI mereka.
  - Header `Profile-Update-Interval: 24` (dalam jam).
  - Header `Content-Disposition: attachment; filename="{profile_name}"`.
- [ ] **Validasi Token & Keamanan**:
  - Verifikasi token ke basis data: pastikan akun belum expired dan kuota belum habis.
  - Kirim HTTP 403 Forbidden / HTTP 404 jika token tidak valid atau expired.

---

### 🌐 Phase 2: Built-in Streaming & AI Unlocker (WARP Smart Routing)
> **Target**: Membuka blokir IP datacenter untuk ChatGPT, Claude, Netflix, Disney+, Reddit, dan bank lokal tanpa mengurangi kecepatan browsing umum.

- [ ] **Outbound Cloudflare WARP / WireGuard**:
  - Tambahkan template outbound `wireguard` lokal di `resources/sing-box/config.json`.
  - Dukungan auto-register / generate WARP private key & peer endpoint via Go background service.
- [ ] **Smart Rule Routing di sing-box v1.14+**:
  - Konfigurasikan `route.rules`:
    - Geosite AI (`category-ai-chat-!cn`, `openai`, `anthropic`) $\rightarrow$ diarahkan ke outbound `warp-out`.
    - Geosite Media Streaming (`netflix`, `disney`, `primevideo`, `hulu`) $\rightarrow$ diarahkan ke outbound `warp-out`.
    - Traffic umum & direct Indonesia $\rightarrow$ tetap melalui direct IP server untuk latensi terendah dan kecepatan maksimal.
- [ ] **WARP Health-Checker Goroutine**:
  - Background goroutine di Go untuk ping/curl berkala ke endpoint WARP; otomatis bypass ke direct jika tunnel WARP mengalami kendala (*graceful fallback*).

---

### 🛡️ Phase 3: AdBlock & Malware Shield Per-User
> **Target**: Menghemat kuota hingga 30% dan melindungi pengguna dari iklan spam, situs judi online, dan malware.

- [ ] **Integrasi Field Basis Data**:
  - Manfaatkan field `Adblock bool` yang telah tersedia di [`model.User`](internal/domain/model/user.go).
- [ ] **Injeksi DNS & Routing Rules di sing-box**:
  - Jika akun memiliki `Adblock: true`, asosiasikan user tersebut ke routing rule yang menolak domain `geosite:category-ads-all` dan malware (`action: reject` atau resolve ke `0.0.0.0`).
  - Opsional: Arahkan query DNS user tersebut ke upstream secure AdGuard DNS / NextDNS.
- [ ] **API Pengaturan User**:
  - Endpoint `POST /api/v1/user/settings` untuk toggle fitur AdBlock secara instan dari dashboard/portal pelanggan.

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
> **Target**: Portal mandiri profesional bagi pelanggan untuk memantau status langganan dan mengelola akun mereka.

- [ ] **Halaman Web `/portal?token={USER_TOKEN}`**:
  - Visualisasi kuota: Terpakai vs Total kuota (Progress bar interaktif).
  - Status masa aktif & countdown hari tersisa.
  - Riwayat protokol yang tersedia (VMess, VLESS, Trojan).
- [ ] **1-Click Import Actions**:
  - Tombol instan yang membuka aplikasi dengan scheme URL:
    - *Import to Clash* (`clash://install-config?url=...`)
    - *Import to Shadowrocket* (`shadowrocket://add/sub://...`)
    - *Import to Sing-box*
  - Tampilan QR Code subscription untuk scan kamera HP.
- [ ] **Fitur Reset Kredensial**:
  - Tombol **"Reset UUID / Password"**: Jika link pengguna tidak sengaja tersebar, user dapat menggenerasi UUID baru secara mandiri tanpa bantuan admin.

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
