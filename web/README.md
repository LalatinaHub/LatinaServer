# LatinaServer Web Frontend (Hugo Static Site & User Portal)

Direktori `./web` berisi kode sumber antarmuka website LatinaServer yang dibangun menggunakan **Hugo Static Site Generator** dengan inspirasi desain editorial/brutalist dari tema [OliverObst/no-fate](https://github.com/OliverObst/no-fate).

Situs ini menjalankan dua peran:
1. **Halaman `index` (`/`)**: Halaman Samaran (*Camouflage / Decoy Site*) bertema jurnal riset sistem terdistribusi (*Lalatina Distributed Systems & Network Resilience Lab*) untuk mengelabui sensor jaringan, ISP, dan crawler otomatis.
2. **Halaman `portal` (`/portal`)**: Halaman Login & Self-Service Portal Pengguna untuk mengelola kredensial VPN, trial tamu, konfigurasi node, import subscription, dan telemetri server.

---

## 🛠️ Prasyarat Tooling (Hugo Extended)

Website ini membutuhkan **Hugo Extended Edition** (v0.120.0 atau lebih baru) untuk kompilasi stylesheet SCSS/CSS dan pipeline aset Hugo Pipes.

### 1. Instalasi via Go Toolchain
```bash
go install -tags extended github.com/gohugoio/hugo@latest
```
*Catatan: Pastikan direktori `$GOPATH/bin` atau `%USERPROFILE%\go\bin` telah ditambahkan ke variabel lingkungan `PATH`.*

### 2. Instalasi di Windows
- **Via Windows Package Manager (winget)**:
  ```powershell
  winget install Hugo.Hugo.Extended
  ```
- **Via Chocolatey**:
  ```powershell
  choco install hugo-extended
  ```
- **Via GitHub Releases**:
  Unduh binary `hugo_extended_*_windows-amd64.zip` dari [Hugo Releases](https://github.com/gohugoio/hugo/releases), ekstrak file `hugo.exe`, lalu tambahkan ke `PATH`.

### 3. Instalasi di Linux
- **Debian / Ubuntu**:
  Unduh paket `.deb` edisi extended dari [Hugo Releases](https://github.com/gohugoio/hugo/releases):
  ```bash
  wget https://github.com/gohugoio/hugo/releases/download/v0.165.0/hugo_extended_0.165.0_linux-amd64.deb
  sudo dpkg -i hugo_extended_0.165.0_linux-amd64.deb
  ```
- **Arch Linux**:
  ```bash
  sudo pacman -S hugo
  ```
- **Via Snap**:
  ```bash
  snap install hugo --channel=extended
  ```

### 4. Instalasi di macOS
```bash
brew install hugo
```

### Verifikasi Instalasi
Jalankan perintah berikut untuk memastikan edisi extended telah terpasang:
```bash
hugo version
# Output diharapkan memuat flag "+extended", contoh:
# hugo v0.165.0-...-extended ...
```

---

## 📁 Struktur Direktori

```
web/
├── hugo.toml            # Konfigurasi utama proyek Hugo
├── README.md            # Dokumentasi & panduan developer
├── TODO.md              # Roadmap & implementasi bertahap
├── web.go               # Wrapper handler Go webserver
├── content/             # Konten Markdown
│   ├── posts/           # Catatan riset & artikel samaran
│   ├── records/         # Spesifikasi teknis cluster & arsip
│   └── portal/          # Konten halaman self-service portal
├── layouts/             # Template HTML Go
│   ├── _default/        # Base layout, single, dan list
│   └── partials/        # Komponen header, footer, masthead, nav
├── assets/              # Source modular CSS & JS
│   ├── css/             # Stylesheet native (tokens, typography, grid, portal)
│   └── js/              # Client controller modular
├── static/              # Asset statis langsung disalin ke output
│   ├── fonts/           # Web fonts
│   ├── images/          # Diagram, ikon, favicon
│   └── js/              # Pustaka mandiri (qrcode.min.js)
└── dist/                # Output build statis (dihasilkan oleh Hugo)
```

---

## 🚀 Alur Kerja Pengembangan & Build

### Menjalankan Development Server (Live Reload)
```bash
cd web
hugo server -D
```
Akses di browser melalui `http://localhost:1313`.

### Menghasilkan Output Statis (Production Build)
```bash
cd web
hugo --minify -d dist
```
Hasil kompilasi akan disimpan di `./web/dist` dan siap disajikan oleh Go webserver LatinaServer (`internal/service/web`).
