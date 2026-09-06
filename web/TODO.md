# Web (Index Samaran & User Portal) — Roadmap & Implementation Plan

Dokumen ini mendefinisikan rencana kerja dan spesifikasi arsitektur untuk membangun kembali antarmuka website pada `./web` menggunakan **Hugo Static Site Generator** dengan inspirasi desain editorial/brutalist dari tema [OliverObst/no-fate](https://github.com/OliverObst/no-fate).

Sistem ini memisahkan dua peran utama:
1. **Halaman `index` (`/`)**: **Halaman Samaran (Camouflage / Decoy Site)** berupa publikasi jurnal riset sistem dan arsitektur jaringan terdistribusi yang sangat kredibel, bersih, dan elegan, untuk mengelabui sensor jaringan, ISP, serta crawler otomatis.
2. **Halaman `portal` (`/portal`)**: **Halaman Login & Self-Service Portal Pengguna** untuk pelanggan mengelola kredensial VPN, akun trial tamu, konfigurasi node, import subscription, dan telemetri server secara mandiri.

---

## 🏛️ Konsep Desain & Arsitektur Utama

### 1. Filosofi Tema No Fate ([OliverObst/no-fate](https://github.com/OliverObst/no-fate))
- **Editorial Voice & Structured Records**: Mengadopsi tata letak bergaya jurnal akademik/editorial independen dengan grid presisi, tipografi cetak (*print surface*), dan hirarki informasi yang ketat.
- **Zero Heavy JS / Zero AI-Slop**:
  - Tanpa Tailwind runtime atau bundler Node.js raksasa.
  - Tanpa elemen AI-slop generik (menghilangkan bulatan gradien blur mengambang, efek glassmorphism murahan, atau card seragam tanpa karakter).
  - Menggunakan CSS semantik native dengan dua lapisan visual:
    - `Clean Style`: Tenang, lapang, tipografi langsung dan keterbacaan tinggi.
    - `Wild / Editorial Layer`: Blok warna hitam pekat berpadu aksen *signal* (amber/emas hangat khas rambut Darkness `#f59e0b` dan aksen merah ruby jepit rambut `#e11d48`), garis pembatas tegas, dan format monospaced untuk data teknis.
- **High Anti-Censorship & Camouflage Value**: Inspeksi manual maupun deep-packet inspection (DPI) pada level HTTP akan mendeteksi website sebagai publikasi riset sistem terdistribusi ("*Lalatina Distributed Systems & Network Resilience Lab*").

### 2. User Portal (`/portal`)
- Mengadopsi bahasa visual editorial No Fate yang sama (monokromatik kontras tinggi, aksen signal amber/ruby, kotak terminal brutalist, monospaced data readouts).
- **Dual Authentication Access**:
  - **Member Token Login**: Input token anggota 8-karakter atau akses langsung via query URL `/portal?token={TOKEN}`.
  - **Surat Jalan Tamu (Guest Trial Mode)**: 1-Click tiket uji coba 24 jam dengan batas kecepatan 2 Mbps via `/api/v1/trial`.
- **Fitur Lengkap (Paritas LatinaBot)**:
  - Profil & Kuota Akun (visualisasi data bar segmen RPG/editorial).
  - 1-Click Import (Clash Meta YAML, Sing-box JSON, Shadowrocket iOS, v2rayNG/NekoBox).
  - Generator QR Code interaktif untuk scan dari aplikasi VPN smartphone.
  - VPN Protocol Switcher (`VMess`, `VLESS`, `Trojan`).
  - Edge Server Cluster & Multi-hop Relay Switcher.
  - Reset UUID Mandiri & Auto-reload sing-box core.
  - Rotasi Token Akses Baru.
  - Toggle Perisai AdBlock & Anti-Malware (Dual Layer DNS & Route Reject).
  - Live Telemetry Server (CPU, RAM, Disk, Uptime).
  - Daftar Wildcard Domain & Panduan Bug Host SNI.
  - Komitmen Privasi Zero-Log & Tautan Donasi Trakteer.

### 3. Pipeline Integrasi Hugo & Go Webserver
- `./web` dikonfigurasi sebagai root project Hugo.
- Output build Hugo (`hugo --minify`) disimpan di direktori `./web/dist`.
- Go Webserver (`internal/service/web`):
  - Root `/` menyajikan `web/dist/index.html` (halaman samaran) dan sub-rute artikel (`/posts/*`, `/records/*`, `/archive/*`).
  - Rute `/portal` menyajikan `web/dist/portal/index.html` (halaman login & portal ksatria).
  - Rute `/api/v1/*` tetap menangani endpoint API dinamis.
  - Asset statis (`/css/*`, `/js/*`, `/fonts/*`) disajikan secara optimal dengan gzip compression.

---

## 📋 Rencana Tahapan Kerja (Checklist TODO)

### Fase 1: Setup Lingkungan Hugo & Struktur Proyek `./web`
- [x] **1.1. Prasyarat Tooling Hugo**:
  - Dokumentasikan dan sediakan instruksi instalasi Hugo Extended (via `go install -tags extended github.com/gohugoio/hugo@latest` atau unduhan binary resmi).
- [x] **1.2. Inisialisasi Proyek Hugo di `./web`**:
  - Buat file konfigurasi `web/hugo.toml` dengan namespace parameter terinspirasi No Fate:
    ```toml
    baseURL = "/"
    languageCode = "id"
    title = "Lalatina Systems — Distributed Network Laboratory"
    
    [params.noFate]
      style = "wild"
      accent = "signal"
      showReadingTime = true
      themeMode = "dark"
    ```
- [x] **1.3. Struktur Direktori Hugo**:
  - `web/content/` (halaman samaran, artikel, dan halaman portal).
  - `web/layouts/` (template HTML Go template).
  - `web/assets/` (source CSS, modular stylesheets).
  - `web/static/` (font, icons, QRCode library, assets statis).

---

### Fase 2: Desain Sistem CSS & Layout Terinspirasi No Fate
- [x] **2.1. Arsitektur CSS Modular Native (Tanpa Node.js / NPM)**:
  - `web/assets/css/tokens.css`: Token warna obsidian base (`#0c0a08`), plat zirah gelap (`#14100c`), kertas cetak editorial (`#1c1611`), aksen signal amber Darkness (`#f59e0b`), aksen ruby cross clip (`#e11d48`), serta border bronze (`#3d3023`).
  - `web/assets/css/typography.css`: Tipografi editorial berkontras tinggi (Cinzel / Serif tajam untuk headline masthead, Plus Jakarta Sans / Inter untuk teks artikel, JetBrains Mono untuk metrik data).
  - `web/assets/css/grid.css`: Sistem grid editorial No Fate, masthead strip, signal blocks, dan layout cetak presisi.
  - `web/assets/css/portal.css`: Komponen interaktif terminal portal (kartu brutalist, tombol taktil, segmen meter bar status, tabs).
- [x] **2.2. Template Layout Hugo (`web/layouts/`)**:
  - `layouts/_default/baseof.html`: Kerangka HTML5 semantik utama dengan accessible dark theme.
  - `layouts/partials/header.html`: Masthead editorial dengan navigasi samaran (Research, Records, Lab Notes) dan tautan diskret ke Portal.
  - `layouts/partials/footer.html`: Footer bergaya arsip riset independen dengan lisensi terbuka dan jaminan integritas data.
  - `layouts/index.html`: Template beranda editorial samaran (Opening Proposition, Recent Research Notes, Systems Registry).
  - `layouts/_default/single.html` & `layouts/_default/list.html`: Template untuk artikel dan arsip riset samaran.

---

### Fase 3: Konten Samaran (Decoy Content & Archives)
- [x] **3.1. Artikel Riset Arsitektur Jaringan (Kredibel & Berkualitas Tinggi)**:
  - `web/content/posts/resilient-packet-routing-under-adversarial-conditions.md`: Catatan riset tentang perutean paket tahan gangguan dan mitigasi kehilangan data.
  - `web/content/posts/transport-layer-multipath-and-stealth-tunneling.md`: Analisis performa protokol modern (multipath, low-overhead encapsulation).
  - `web/content/posts/zero-trust-edge-gateway-topologies.md`: Dokumentasi arsitektur edge proxy dan pertahanan lapisan jaringan.
- [x] **3.2. Structured Records & About**:
  - `web/content/records/spec-v1.14.md`: Dokumentasi spesifikasi formal node cluster.
  - `web/content/about.md`: Halaman profil laboratorium riset sistem terdistribusi.

---

### Fase 4: Halaman Portal Pengguna (`/portal`)
- [x] **4.1. Layout & Konten Portal**:
  - `web/content/portal/_index.md` dan `web/layouts/portal/list.html`: Halaman portal mandiri yang serasi dengan tema No Fate.
- [x] **4.2. UI Autentikasi Ganda**:
  - Form input token member 8-karakter dengan tombol taktil "*Masuk ke Portal*".
  - Opsi "*Akun Uji Coba (Trial 24 Jam)*" untuk generate akun instan tanpa registrasi via `/api/v1/trial`.
  - Notifikasi toast banner sistem dengan status responsif.
- [x] **4.3. Dashboard Pengguna & Kontrol Mandiri**:
  - **Statistik Akun**: Kuota terpakai vs total (segmented meter bar), status masa aktif, protokol aktif, status keamanan zero-log.
  - **1-Click Invocations**: Kotak salin konfigurasi Clash Meta, Sing-box, Shadowrocket, dan v2rayNG URI.
  - **QR Code Engine**: Generator QR client-side mandiri (menggunakan pustaka ringan `web/static/js/qrcode.min.js`).
  - **Switcher Konfigurasi**: Pilihan protokol (VMess/VLESS/Trojan), server edge cluster, dan multi-hop relay.
  - **Keamanan Mandiri**: Reset UUID + reload sing-box, rotasi token login, dan toggle AdBlock.
  - **Telemetri Server**: Monitor beban CPU, RAM, Disk, Uptime dengan visual bar segmen editorial monospaced.
  - **Wildcard SNI & Bug Host**: Daftar domain dan panduan injeksi SNI.
- [x] **4.4. Client Controller Ringan**:
  - Controller modular Vanilla JS murni (`web/assets/js/portal.js`) untuk menghubungkan UI dengan API `/api/v1/portal/*` dan `/api/v1/trial`.

---

### Fase 5: Integrasi Go Backend & Caddy Reverse Proxy
- [x] **5.1. Build Script & Output Directory**:
  - Konfigurasikan Hugo agar melakukan render output statis langsung ke `web/dist`.
  - Sediakan perintah build di root atau Makefile (`hugo --minify -s web -d dist`).
- [x] **5.2. Penyesuaian Router Gin (`internal/service/web/server.go`)**:
  - Perbarui rute root `/` agar melayani `web/dist/index.html` (halaman samaran) secara statis jika tanpa token.
  - Jika query token diberikan pada root (`/?token=xxx`), lakukan redirect cerdas ke `/portal?token=xxx`.
  - Sajikan rute `/portal` dan `/portal/*` langsung ke `web/dist/portal/index.html` (atau fallback statis).
  - Daftarkan penyajian aset statis (`/css`, `/js`, `/fonts`, `/images`, `/posts`) dari direktori `web/dist`.
- [x] **5.3. Kompatibilitas Caddy (`resources/caddy/caddy.json`)**:
  - Pastikan semua request HTTP reguler untuk path samaran maupun `/portal` dialirkan ke port webserver internal.

---

### Fase 6: Pengujian, Optimasi & Verifikasi Menyeluruh
- [x] **6.1. Validasi Build Hugo**:
  - Pastikan perintah build Hugo menghasilkan HTML semantik valid tanpa broken link.
- [x] **6.2. Unit & Integration Testing Go**:
  - Perbarui dan jalankan unit test pada `internal/service/web/...` untuk memvalidasi:
    - Handler `/` melayani halaman samaran dengan status `200 OK`.
    - Handler `/portal` melayani halaman portal dengan status `200 OK`.
    - Seluruh API `/api/v1/portal/*` dan `/api/v1/trial` tetap merespons dengan benar.
- [x] **6.3. Verifikasi Bebas AI-Slop & Kualitas Visual**:
  - Evaluasi visual: Memastikan tipografi bersih, tajam, kontras tinggi, bernuansa editorial terhormat, dan bebas dari artefak visual generik.
- [x] **6.4. Verifikasi Binary Standalone**:
  - Jalankan `go build ./cmd/latinaserver` dan pastikan binary terkompilasi bersih tanpa dependensi runtime eksternal yang hilang.

---

## 🛠️ Perintah Eksekusi Singkat (Developer Reference)

```bash
# 1. Instalasi Hugo (jika belum tersedia)
go install -tags extended github.com/gohugoio/hugo@latest

# 2. Build situs web Hugo ke web/dist
cd web && hugo --minify -d dist && cd ..

# 3. Jalankan pengujian Go
go test ./internal/service/web/... -v

# 4. Kompilasi binary LatinaServer
go build ./cmd/latinaserver
```
