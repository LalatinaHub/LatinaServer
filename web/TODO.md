# TODO.md — Roadmap Eksekusi Desain Ulang LatinaServer Web
**Rencana Kerja Terstruktur Berdasarkan Spesifikasi `DESIGN.md` (Tailwind CSS Edition)**  
*Protokol: Deep Agents Planning & Todos | Status: Ready for Implementation*

---

## 🎯 Ringkasan Eksekutif & Tujuan Utama
Mengubah seluruh antarmuka frontend LatinaServer (`./web`) dari tampilan lama yang gelap, berat, dan membingungkan (*obsidian brutalist*) menjadi antarmuka **Premium Utilitarian Minimalism & Editorial UI** berbasis **Tailwind CSS**.

Seluruh kapabilitas fungsional backend (`internal/service/web`) dan alur pengguna dipertahankan 100%, meliputi penyamaran publik (*decoy research journal*), autentikasi token, akun uji coba instan 24 jam, ekspor profil proksi 1-klik (Clash Meta, Sing-box, Shadowrocket, Raw URI), generator QR code, konfigurasi server/protokol, sakelar proteksi AdBlock, rotasi kredensial keamanan, dan telemetri server real-time.

---

## 📊 Matriks Progres Milestones

| Milestone | Deskripsi Pekerjaan | Estimasi Subtask | Status |
|---|---|---|---|
| **M1** | Tooling Setup & Konfigurasi Tailwind CSS | 3 Subtask | `[x] Completed` |
| **M2** | Kerangka Dasar Hugo & Aset Tipografi Editorial | 4 Subtask | `[x] Completed` |
| **M3** | Desain Ulang Beranda Publik / Camouflage Journal | 4 Subtask | `[x] Completed` |
| **M4** | Desain Ulang Portal Mandiri Bento Grid 4 Kuadran | 5 Subtask | `[x] Completed` |
| **M5** | Integrasi JS Controller & Sinkronisasi Dashboard | 3 Subtask | `[x] Completed` |
| **M6** | Audit Kualitas, Testing Fungsional, & Build Verifikasi | 4 Subtask | `[x] Completed` |

---

## 📋 Rincian Tahapan Kerja (Actionable Checklist)

### Milestone 1: Tooling Setup & Konfigurasi Tailwind CSS
Fondasi konfigurasi tema, variabel warna, font stacks, dan pipeline build CSS.

- [x] **1.1. Pembuatan Konfigurasi Tailwind (`web/tailwind.config.js`)**
  - Daftarkan path konten: `./layouts/**/*.html`, `./content/**/*.md`, `./assets/js/**/*.js`, `./static/**/*.html`.
  - Konfigurasikan palet warna resmi: `canvas` (`#FBFBFA`/`#FFFFFF`), `surface` (`#FFFFFF`/`#F7F6F3`), `border` (`#EAEAEA`, hover `#CCCCCC`), `charcoal` (`#111111`, muted `#2F3437`), `secondary` (`#787774`).
  - Daftarkan 4 warna semantik *Muted Pastel* (Green, Blue, Yellow, Red) dengan pasangan warna latar, teks, dan batasnya.
  - Konfigurasikan font stacks: `font-serif` (Newsreader/Instrument Serif), `font-sans` (Geist Sans/SF Pro), `font-mono` (Geist Mono/SF Mono).
  - Tentukan ukuran sudut tegas: `rounded-crisp` (`4px`), `rounded-card` (`8px`), `rounded-card-lg` (`12px`).
  - Tentukan bayangan ultra-halus: `shadow-diffuse: 0 2px 8px rgba(0,0,0,0.04)`.

- [x] **1.2. Pembuatan Berkas Entry Stylesheet (`web/assets/css/input.css`)**
  - Masukkan direktif `@tailwind base; @tailwind components; @tailwind utilities;`.
  - Atur aturan dasar `@layer base` untuk kanvas latar belakang, warna teks utama, font default, serta perataan antialiasing peramban.

- [x] **1.3. Pipeline Build & Skrip Kompilasi CSS**
  - Siapkan skrip kompilasi Tailwind (menggunakan Tailwind CLI standalone binary atau npm command).
  - Target keluaran kompilasi diarahkan ke `web/static/css/styles.css` agar siap disajikan oleh Hugo dan Go webserver.

---

### Milestone 2: Kerangka Dasar Hugo & Aset Tipografi Editorial
Pembersihan artefak visual usang dan perombakan struktur layout utama.

- [x] **2.1. Perombakan Layout Utama (`web/layouts/_default/baseof.html`)**
  - Hapus pemuatan font lama (*Cinzel*, *Plus Jakarta Sans*).
  - Pasang link Google Fonts resmi untuk `Newsreader` (bobot 400, 500, 600 italic), `Geist Sans` (bobot 400, 500, 600), dan `Geist Mono` (bobot 400, 500).
  - Tautkan stylesheet hasil kompilasi Tailwind `/css/styles.css`.
  - Hapus bundle CSS lama yang memuat styling obsidian brutalist (`tokens.css`, `typography.css`, `grid.css`, `portal.css`).

- [x] **2.2. Desain Ulang Navigasi Header (`web/layouts/partials/header.html`)**
  - Terapkan navigasi minimalis bergaris batas bawah `1px solid #EAEAEA`.
  - Tampilkan identitas brand editorial bersih: `LALATINA SYSTEMS` (sans-serif uppercase tracking lega).
  - Tautan menu navigasi jernih: *Overview*, *Research*, *Records*, *Dossier*.
  - Sediakan tombol akses diskret menuju Portal Pengguna di pojok kanan dengan tombol outline minimalis.

- [x] **2.3. Desain Ulang Penutup Footer (`web/layouts/partials/footer.html`)**
  - Terapkan footer berpenampilan dokumen arsip riset independen.
  - Tampilkan ringkasan mandat laboratorium, garansi kebijakan zero-log, dan verifikasi status enkripsi TLS 1.3 tanpa teks klise.
  - Cantumkan hak cipta, lisensi terbuka (CC BY-NC 4.0 & MIT), dan tautan portal mandiri.

- [x] **2.4. Penyesuaian Template Artikel & Arsip (`web/layouts/_default/single.html` & `list.html`)**
  - Pastikan halaman artikel riset (`/posts/*`) dan catatan teknis (`/records/*`) dirender dengan tipografi editorial yang nyaman dibaca (`max-w-3xl`, `leading-relaxed`, heading serif).

---

### Milestone 3: Desain Ulang Beranda Publik / Camouflage Journal (`web/layouts/index.html`)
Membangun tampilan muka publik yang kredibel sebagai jurnal riset jaringan terdistribusi.

- [x] **3.1. Hero Section Editorial (Masthead Proposition)**
  - Terapkan *macro-whitespace* (`py-20 md:py-28`, `max-w-5xl mx-auto px-4 sm:px-6`).
  - Volume badge: `text-xs font-mono uppercase tracking-badge text-secondary`.
  - Headline utama bergaya editorial serif: `text-3xl md:text-5xl font-serif font-medium tracking-tight-title text-charcoal`.
  - Deskripsi ringkasan riset terdistribusi dan mitigasi inspeksi mendalam (DPI).
  - Tombol aksi ganda: Tombol utama `Masuk ke Portal Pengguna &rarr;` dan tombol sekunder `Coba Akses Gratis 24 Jam`.

- [x] **3.2. Callout Dokumen Protokol (Minimal Callout Box)**
  - Ganti *signal block* brutalist lama dengan kartu bento datar berbingkai `border border-[#EAEAEA]` dengan latar `#FFFFFF`.
  - Tampilkan spesifikasi kebijakan zero-log dan indikator status enklave aktif dengan tag pastel hijau.

- [x] **3.3. Bento Grid Catatan Riset Terbaru**
  - Tampilkan kartu artikel riset menggunakan layout grid asimetris.
  - Format kartu: metadata tanggal & estimasi waktu baca (`font-mono text-xs text-secondary`), judul artikel, ringkasan pendek, dan tautan baca.
  - Efek hover kartu: elevasi mikro `-translate-y-0.5` dan border `#CCCCCC` tanpa bayangan berat.

- [x] **3.4. Registri Topologi Node Aktif (Systems Registry Table)**
  - Rombak tabel node kluster (SG-01 Singapura, JP-01 Tokyo, US-01 San Jose, ID-01 Jakarta).
  - Tampilkan kolom identitas node, wilayah, protokol, mitigasi proteksi, latensi, dan status operasional.
  - Gunakan badge status pil pastel hijau (`bg-pastel-green-bg text-pastel-green-text`).

---

### Milestone 4: Desain Ulang Portal Mandiri Bento Grid 4 Kuadran (`web/layouts/partials/portal.html`)
Pusat manajemen mandiri pengguna yang intuitif, cepat, dan berfokus pada kemudahan akses.

- [x] **4.1. Tampilan Gerbang Autentikasi Pengguna & Akun Tamu (`#portal-login-view`)**
  - Kontainer kartu terpusat di layar (`max-w-md w-full bg-surface border border-border rounded-card p-8 shadow-diffuse mx-auto my-12`).
  - Form Login Member: Input token akses dengan styling dokumen bersih + tombol submit utama solid dark (`#111111`).
  - Garis pemisah minimalis dengan label *"ATAU"*.
  - Kotak Akses Tamu: Opsi instan aktivasi akun uji coba 24 jam / batas 2 Mbps (VMess) via 1-klik tombol outline.

- [x] **4.2. Header Dasbor & Status Sesi (`#portal-dashboard-view`)**
  - Baris informasi sesi aktif: Token pengguna (`font-mono font-semibold`), masa berlaku akun dengan sisa hari, status aktif/kedaluwarsa (badge pastel), dan tombol Keluar Sesi.

- [x] **4.3. Kuadran 1: Ekspor Profil Klien 1-Klik & Scanner QR Code (Prioritas Utama)**
  - Tab selector klien: *Clash Meta / Mihomo*, *Sing-box JSON*, *Shadowrocket iOS*, *v2rayNG / NekoBox*.
  - Kotak salin konfigurasi (`copy-box`) dengan tombol satu sentuhan ke clipboard.
  - Tombol aksi deep-link langsung (`Buka di Aplikasi &rarr;`).
  - Wadah QR Code responsif di sebelah kanan dengan bingkai bersih untuk pemindaian instan dari ponsel.

- [x] **4.4. Kuadran 2 & 3: Konfigurasi Rute & Kontrol Keamanan Kredensial**
  - **Kuadran 2 (Switcher)**: Dropdown Protokol (VMess/VLESS/Trojan), Dropdown Node Server, Dropdown Relay Multi-Hop, dan tombol *Simpan Perubahan*.
  - **Kuadran 3 (Keamanan)**:
    - Sakelar toggle datar untuk Proteksi DNS AdBlock & Anti-Malware.
    - Tombol aksi reset UUID baru dengan konfirmasi modal/dialog yang aman.
    - Tombol aksi rotasi token akses login baru (badge pastel merah).

- [x] **4.5. Kuadran 4: Telemetri Beban Node & Registri SNI Bug Host**
  - Tiga batang meter linier ramping (`h-1.5`) untuk utilisasi CPU, RAM, dan Disk beserta teks angka persentase.
  - Teks waktu aktif sistem (uptime).
  - Tabel direktori SNI Wildcard / Bug Host dengan tombol *Salin Host* per baris.

---

### Milestone 5: Integrasi JS Controller & Sinkronisasi Dashboard
Sinkronisasi interaktivitas antarmuka dan penanganan data real-time.

- [x] **5.1. Penyesuaian DOM Manipulator pada `web/assets/js/portal.js`**
  - Perbarui pemetaan elemen DOM dan selektor kelas agar sesuai dengan markup Tailwind baru.
  - Perbarui fungsi render meter telemetri linier (`renderMeterSegments` / pengatur lebar persentase CSS).
  - Perbarui fungsi badge status agar menerapkan kelas pastel semantik (`bg-pastel-green-bg` saat aktif, `bg-pastel-red-bg` saat kedaluwarsa).

- [x] **5.2. Desain Ulang Toast Notification & Loading Overlay**
  - Ganti container toast lama dengan notifikasi datar modern di sudut kanan bawah berbingkai `border border-[#EAEAEA] bg-white shadow-diffuse`.
  - Ganti overlay loading gelap menjadi backdrop blur lembut dengan indikator pemrosesan dokumen yang tenang.

- [x] **5.3. Sinkronisasi Halaman Mandiri `web/static/dashboard.html`**
  - Terapkan gaya tema minimalis Tailwind ke berkas `dashboard.html` agar tidak ada lagi tampilan lama yang tertinggal jika diakses melalui rute langsung `/dashboard`.

---

### Milestone 6: Audit Kualitas, Testing Fungsional, & Build Verifikasi
Pengujian menyeluruh untuk menjamin stabilitas fungsional dan kepatuhan desain.

- [x] **6.1. Audit Kepatuhan Desain (*Design System Compliance Audit*)**
  - [x] Tidak ada font `Inter`, `Roboto`, atau `Open Sans`.
  - [x] Tidak ada kelas bayangan tebal Tailwind (`shadow-md`, `shadow-lg`, `shadow-xl`).
  - [x] Tidak ada warna primer saturasi tinggi (`bg-blue-600`, `bg-emerald-500`, dll.).
  - [x] Tidak ada `rounded-full` pada kartu bento atau tombol utama.
  - [x] Tidak ada emoji di seluruh kode markup, judul, paragraf, atau tombol.
  - [x] Seluruh kartu menggunakan `border border-[#EAEAEA]` dengan latar `#FFFFFF`.
  - [x] Tombol CTA utama menggunakan `bg-[#111111] text-white hover:bg-[#2F3437]`.

- [x] **6.2. Pengujian Fungsional Seluruh Alur Pengguna**
  - [x] Login via token akses dan deteksi token dari URL (`/?token=...` atau `/portal?token=...`).
  - [x] Pembuatan akun uji coba gratis 24 jam (VMess 2 Mbps) via tombol trial instan.
  - [x] Penyalinan konfigurasi ke clipboard pada semua tab (Clash, Sing-box, Shadowrocket, Raw).
  - [x] Rendering QR Code saat berganti tab atau server.
  - [x] Penggantian protokol, server edge, dan relay dengan tombol Simpan.
  - [x] Sakelar toggle DNS AdBlock (POST `/api/v1/portal/toggle-adblock`).
  - [x] Reset UUID baru dan rotasi token baru.
  - [x] Polling berkala telemetri server (CPU, RAM, Disk, Uptime) setiap 12 detik.
  - [x] Penyalinan domain dari tabel SNI Bug Host.

- [x] **6.3. Verifikasi Build Hugo & Integrasi Go**
  - [x] Jalankan kompilasi Tailwind CSS ke `web/static/css/styles.css`.
  - [x] Konfigurasi build Hugo diselaraskan (`publishDir = "dist"`).
  - [x] Static assets disajikan langsung oleh route Gin (`r.Static("/css", ...)`).

- [x] **6.4. Eksekusi Unit Test Go Web Service**
  - [x] Jalankan `go test ./internal/service/web/... -v` dan pastikan seluruh test suite lolos (`PASS`).

---

## 🛠️ Perintah Eksekusi Utama (Developer Cheatsheet)

```bash
# 1. Kompilasi Tailwind CSS (Watch mode saat development)
npx tailwindcss -i ./web/assets/css/input.css -o ./web/static/css/styles.css --watch

# 2. Kompilasi Tailwind CSS (Production Minified)
npx tailwindcss -i ./web/assets/css/input.css -o ./web/static/css/styles.css --minify

# 3. Build Static Site Hugo ke direktori dist
cd web && hugo --minify -d dist && cd ..

# 4. Jalankan Unit Test Go Web Service
go test ./internal/service/web/... -v

# 5. Kompilasi Binary Utama LatinaServer
go build -o latinaserver.exe ./cmd/latinaserver
```
