# DESIGN.md — Panduan & Aturan Desain Antarmuka LatinaServer
**Dokumen Spesifikasi: Premium Utilitarian Minimalism & Editorial UI Architecture (Tailwind CSS Edition)**  
*Versi: 2.1.0 | Status: Active Specification | Framework: Tailwind CSS + Hugo | Target: `./web`*

---

## 1. Filosofi & Manifesto Desain

### 1.1 Latar Belakang & Pemisahan dari Desain Lama
Desain visual sebelumnya ditinggalkan secara total. Pendekatan lama yang menggunakan warna gelap pekat (*obsidian*), efek kilau neon (*cyan/amber glow*), bayangan tebal brutalist, dan tumpukan kotak terminal menghasilkan antarmuka yang membingungkan, melelahkan mata, dan mengaburkan fungsi utama.

Website ini didesain ulang dari fondasi pertama dengan mengadopsi protokol **Premium Utilitarian Minimalism & Editorial UI** yang diimplementasikan menggunakan **Tailwind CSS**. Antarmuka baru mengedepankan:
- Ketenangan visual (*quiet sophistication*).
- Struktur *macro-whitespace* yang lapang (`py-20` hingga `py-32`, kontainer terbatas `max-w-5xl`).
- Palet monokrom hangat (*warm monochrome*).
- Aksen warna pastel pucat (*washed-out pastels*) yang sangat terkontrol.
- Tipografi dokumen bernilai estetika tinggi seperti media publikasi ilmiah dan workspace modern kelas atas.

### 1.2 Prinsip Inti Desain
1. **Utilitarian & Task-Oriented**: Pengguna datang untuk menyelesaikan tugas (menyalin konfigurasi VPN, menguji koneksi, memilih server, memantau kuota). Setiap utilitas Tailwind disusun untuk mempermudah tugas tersebut tanpa distraksi visual.
2. **Extreme Typographic Contrast**: Hierarki informasi dibangun menggunakan perbedaan jenis huruf (Editorial Serif berkarakter untuk judul, Geometric Sans-Serif untuk navigasi dan kontrol, serta Monospace untuk token, UUID, dan data teknis).
3. **Warna Sebagai Sumber Daya Langka**: Warna dilarang digunakan untuk dekorasi semata. Warna hanya hadir untuk memberi sinyal semantik (status aktif, peringatan kedaluwarsa, atau kategori protokol).
4. **Bento Grid Bergaris Tunggal 1px**: Menggunakan tata letak kotak bento asimetris dengan garis tepi tunggal `border border-[#EAEAEA]` dan sudut melengkung tegas `rounded-lg` atau `rounded-xl`.
5. **Dualitas Fungsi yang Selaras**:
   - **Lapisan Penyamaran Publik (*Decoy / Camouflage*)**: Tampil sebagai jurnal publikasi riset sistem terdistribusi independen yang kredibel, tenang, dan natural untuk mitigasi sensor ISP dan Deep Packet Inspection (DPI).
   - **Lapisan Portal Mandiri (*User Self-Service Enclave*)**: Dasbor fungsional yang langsung menyajikan alat, konfigurasi 1-klik, dan telemetri tanpa jargon yang membingungkan.

---

## 2. Matriks Pemetaan Fungsi Sistem (Functional Matrix)

Seluruh kemampuan fungsional dari arsitektur backend (`internal/service/web`) dan antarmuka web dipertahankan secara utuh tanpa ada fitur yang dipangkas:

| No | Fitur / Tool | Lokasi Endpoint / Handler | Fungsi & Peran Pengguna |
|---|---|---|---|
| 1 | **Autentikasi Token Akses** | `GET /portal?token=...` & `GET /api/v1/portal/profile` | Login pengguna via token akses 8-karakter; ekstraksi otomatis dari parameter URL; penyimpanan aman di `localStorage`. |
| 2 | **Akun Uji Coba 1-Klik (Trial)** | `POST /api/v1/trial` & `GET /api/v1/trial` | Pembuatan akun tamu instan 24 jam dengan batas kecepatan 2 Mbps (VMess) tanpa syarat pendaftaran. |
| 3 | **Ringkasan Akun & Kuota** | `GET /api/v1/portal/profile` | Menampilkan masa berlaku (tanggal & sisa hari), status aktif/kedaluwarsa, kuota bandwidth terpakai, dan UUID aktif. |
| 4 | **Pemilih Protokol & Server** | `POST /api/v1/portal/update-config` | Penggantian protokol (VMess WS, VLESS gRPC, Trojan TLS), pemilihan node server (Singapura, Tokyo, San Jose, Jakarta), dan gerbang relay multi-hop. |
| 5 | **Ekspor Profil Klien (1-Click Invocations)** | `GET /sub?token=...` | Penyediaan tautan langganan universal dan skema aplikasi langsung untuk Clash Meta (`clash://`), Sing-box (`sing-box://`), Shadowrocket (`shadowrocket://`), dan v2rayNG / NekoBox (Raw URI). |
| 6 | **QR Code Interaktif** | Integrasi `qrcode.min.js` | Menghasilkan QR Code dinamis dari profil koneksi aktif untuk pemindaian langsung dari smartphone. |
| 7 | **Proteksi DNS AdBlock & Malware** | `POST /api/v1/portal/toggle-adblock` | Sakelar satu sentuhan untuk mengaktifkan pemblokiran iklan dan pelacak di tingkat core sing-box. |
| 8 | **Rotasi Kunci & Kredensial** | `POST /api/v1/portal/reset-uuid` & `/change-token` | Reset UUID koneksi acak baru dan rotasi token akses login baru dengan dialog konfirmasi aman. |
| 9 | **Telemetri Server Real-Time** | `GET /api/v1/portal/status` | Pemantauan berkala (polling 12 detik) untuk utilisasi CPU, memori RAM, kapasitas penyimpanan disk, dan waktu aktif sistem (uptime). |
| 10 | **Direktori SNI Bug Host / Wildcard** | `GET /api/v1/portal/wildcards` | Daftar domain wildcard untuk injeksi paket kuota edukasi/khusus beserta tombol salin per item. |
| 11 | **Publikasi Riset & Arsip (Camouflage)** | `/posts/`, `/records/`, `/about/` | Artikel riset sistem terdistribusi, spesifikasi cluster v1.14, dan profil etika laboratorium untuk penyamaran lalu lintas. |
| 12 | **Tautan Komunitas & Dukungan** | Link Telegram & Trakteer | Tautan eksternal yang rapi menuju forum komunitas dan apresiasi donasi operasional server. |

---

## 3. Aturan Negatif Mutlak Khusus Tailwind (Tailwind Negative Constraints)

Penggunaan Tailwind CSS **DILARANG MENGIKUTI DEFAULT SAAS GENERIK**. Terapkan aturan pembatasan ketat berikut:

```
+-------------------------------------------------------------------------------+
|                      DAFTAR LARANGAN PENGGUNAAN TAILWIND                      |
+------------------------------------+------------------------------------------+
| DILARANG (BANNED)                  | PENGGANTI YANG DITETAPKAN (REQUIRED)     |
+------------------------------------+------------------------------------------+
| font-sans bawaan (Inter/Roboto)    | font-sans ('Geist Sans', 'SF Pro')       |
|                                    | font-serif ('Newsreader', 'Instrument')  |
|                                    | font-mono ('Geist Mono', 'SF Mono')      |
+------------------------------------+------------------------------------------+
| shadow-md, shadow-lg, shadow-xl    | shadow-none atau shadow-diffuse          |
|                                    | (box-shadow: 0 2px 8px rgba(0,0,0,0.04)) |
+------------------------------------+------------------------------------------+
| bg-blue-600, bg-indigo-500,        | bg-[#111111] untuk primary dark button   |
| bg-emerald-500 (warna saturasi)    | bg-pastel-* (hanya 4 warna pastel pucat) |
+------------------------------------+------------------------------------------+
| bg-gradient-to-r (gradien neon)    | bg-[#FFFFFF] murni atau bg-[#FBFBFA]     |
+------------------------------------+------------------------------------------+
| rounded-full pada kartu / tombol   | rounded-md (6px) s.d. rounded-lg (8px)   |
| (pill container besar)             | rounded-full HANYA untuk status badge    |
+------------------------------------+------------------------------------------+
| Ikon Lucide / Feather tipis        | Phosphor Icons (Bold/Fill) / Radix Icons |
+------------------------------------+------------------------------------------+
| Emoji di dalam teks atau markup    | Ikon SVG monokrom atau teks deskriptif   |
+------------------------------------+------------------------------------------+
| Teks klise AI ("Elevate", dll.)    | Bahasa lugas, teknis, dan manusiawi      |
+------------------------------------+------------------------------------------+
```

---

## 4. Konfigurasi Tailwind CSS (`tailwind.config.js`)

Seluruh desain direkayasa melalui konfigurasi resmi berikut:

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./layouts/**/*.html",
    "./content/**/*.md",
    "./assets/js/**/*.js",
    "./static/**/*.html",
  ],
  theme: {
    extend: {
      colors: {
        // Kanvas & Permukaan
        canvas: {
          DEFAULT: '#FBFBFA',
          pure: '#FFFFFF',
        },
        surface: {
          DEFAULT: '#FFFFFF',
          subtle: '#F7F6F3',
          card: '#FFFFFF',
        },
        // Garis batas struktural (Hairline borders)
        border: {
          DEFAULT: '#EAEAEA',
          subtle: 'rgba(0, 0, 0, 0.06)',
          hover: '#CCCCCC',
        },
        // Teks & Tinta
        charcoal: {
          DEFAULT: '#111111',
          muted: '#2F3437',
        },
        secondary: '#787774',
        ghost: '#999999',

        // Muted Pastels (Aksen semantik terkontrol)
        pastel: {
          green: {
            bg: '#EDF3EC',
            text: '#346538',
            border: '#D3E5D1',
          },
          blue: {
            bg: '#E1F3FE',
            text: '#1F6C9F',
            border: '#C3E4FC',
          },
          yellow: {
            bg: '#FBF3DB',
            text: '#956400',
            border: '#F4E6B8',
          },
          red: {
            bg: '#FDEBEC',
            text: '#9F2F2D',
            border: '#F8D2D4',
          },
        },
      },
      fontFamily: {
        serif: ['Newsreader', 'Instrument Serif', 'Playfair Display', 'Georgia', 'serif'],
        sans: ['Geist Sans', 'SF Pro Display', 'Helvetica Neue', 'sans-serif'],
        mono: ['Geist Mono', 'SF Mono', 'JetBrains Mono', 'Menlo', 'monospace'],
      },
      letterSpacing: {
        'tight-title': '-0.03em',
        'subtle-tight': '-0.015em',
        'badge': '0.05em',
      },
      lineHeight: {
        'title': '1.15',
        'body': '1.6',
      },
      boxShadow: {
        'diffuse': '0 2px 8px rgba(0, 0, 0, 0.04)',
        'dropdown': '0 4px 16px rgba(0, 0, 0, 0.06)',
      },
      borderRadius: {
        'crisp': '4px',
        'card': '8px',
        'card-lg': '12px',
      },
    },
  },
  plugins: [],
}
```

---

## 5. Kamus Utilitas Komponen Tailwind (Tailwind Component Utility Dictionary)

Gunakan pola kelas Tailwind berikut secara konsisten untuk membangun setiap komponen:

### 5.1 Kartu Bento (Bento Grid Card)
```html
<div class="bg-surface border border-border rounded-card p-6 md:p-8 transition-all duration-200 hover:-translate-y-0.5 hover:border-border-hover hover:shadow-diffuse">
  <div class="flex items-center justify-between mb-4">
    <span class="text-xs font-mono font-semibold uppercase tracking-badge text-secondary">Label Kartu</span>
    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono font-medium bg-pastel-green-bg text-pastel-green-text">AKTIF</span>
  </div>
  <h3 class="font-sans text-base md:text-lg font-semibold text-charcoal mb-2">Judul Komponen</h3>
  <p class="text-sm text-secondary leading-body mb-6">Deskripsi fungsional singkat dan terarah tanpa basa-basi.</p>
  <!-- Konten / Kontrol -->
</div>
```

### 5.2 Tombol Aksi (Buttons)
- **Tombol Utama (Primary CTA)**:
  ```html
  <button type="submit" class="inline-flex items-center justify-center bg-charcoal text-white text-sm font-medium px-4 py-2.5 rounded-crisp border border-charcoal hover:bg-charcoal-muted active:scale-[0.98] transition-all duration-150">
    Masuk ke Dasbor &rarr;
  </button>
  ```
- **Tombol Sekunder / Garis Tepi (Secondary Outline)**:
  ```html
  <button type="button" class="inline-flex items-center justify-center bg-transparent text-charcoal border border-border text-sm font-medium px-4 py-2 rounded-crisp hover:bg-surface-subtle hover:border-border-hover active:scale-[0.98] transition-all duration-150">
    Salin Tautan
  </button>
  ```
- **Tombol Bahaya / Rotasi Kunci (Danger Action)**:
  ```html
  <button type="button" class="inline-flex items-center justify-center bg-pastel-red-bg text-pastel-red-text border border-pastel-red-border text-xs font-medium px-3 py-1.5 rounded-crisp hover:opacity-90 active:scale-[0.98] transition-all duration-150">
    Reset UUID Baru
  </button>
  ```

### 5.3 Badges Status & Tag Semantik
- **Status Hijau (Aktif / Online / Tersimpan)**:
  `inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-mono font-semibold tracking-badge uppercase bg-pastel-green-bg text-pastel-green-text border border-pastel-green-border`
- **Status Kuning (Akun Uji Coba / Peringatan)**:
  `inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-mono font-semibold tracking-badge uppercase bg-pastel-yellow-bg text-pastel-yellow-text border border-pastel-yellow-border`
- **Status Merah (Kedaluwarsa / Bahaya)**:
  `inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-mono font-semibold tracking-badge uppercase bg-pastel-red-bg text-pastel-red-text border border-pastel-red-border`
- **Status Biru (Protokol / Info Enklave)**:
  `inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-mono font-semibold tracking-badge uppercase bg-pastel-blue-bg text-pastel-blue-text border border-pastel-blue-border`

### 5.4 Form Input & Dropdown Selector
```html
<div class="space-y-1.5">
  <label for="server-select" class="block text-xs font-mono uppercase tracking-badge text-secondary font-medium">Pilihan Node Server</label>
  <select id="server-select" class="w-full bg-surface border border-border rounded-crisp px-3.5 py-2.5 text-sm text-charcoal font-sans focus:outline-none focus:border-charcoal focus:ring-1 focus:ring-charcoal transition-colors">
    <option value="sg">SG-01 — Singapura (Direct)</option>
    <option value="jp">JP-01 — Tokyo (Relay)</option>
  </select>
</div>
```

### 5.5 Tab Selector Aplikasi Klien
```html
<div class="flex border-b border-border gap-6 mb-6 overflow-x-auto" role="tablist">
  <button type="button" class="pb-2.5 text-sm font-sans font-semibold text-charcoal border-b-2 border-charcoal transition-colors">
    Clash Meta / Mihomo
  </button>
  <button type="button" class="pb-2.5 text-sm font-sans text-secondary hover:text-charcoal border-b-2 border-transparent transition-colors">
    Sing-box (JSON)
  </button>
  <button type="button" class="pb-2.5 text-sm font-sans text-secondary hover:text-charcoal border-b-2 border-transparent transition-colors">
    Shadowrocket (iOS)
  </button>
  <button type="button" class="pb-2.5 text-sm font-sans text-secondary hover:text-charcoal border-b-2 border-transparent transition-colors">
    v2rayNG / NekoBox
  </button>
</div>
```

### 5.6 Keystroke & Token Display (`<kbd>`)
```html
<div class="flex items-center justify-between bg-surface-subtle border border-border rounded-crisp p-3 font-mono text-xs">
  <span class="text-secondary">UUID KONEKSI:</span>
  <span class="font-semibold text-charcoal break-all select-all">e1f83420-7cd3-4e1b-944e-3f9a72b0a192</span>
</div>
```

### 5.7 Sakelar AdBlock (Minimal Toggle Switch)
```html
<label class="relative inline-flex items-center cursor-pointer">
  <input type="checkbox" id="adblock-toggle" class="sr-only peer">
  <div class="w-10 h-6 bg-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-border after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-charcoal"></div>
</label>
```

### 5.8 Telemetry Meters (CPU, RAM, Disk)
```html
<div class="space-y-1.5">
  <div class="flex justify-between text-xs font-mono">
    <span class="text-secondary">BEBAN CPU</span>
    <span class="font-semibold text-charcoal" id="cpu-val">28.4%</span>
  </div>
  <div class="w-full h-1.5 bg-[#F0EFEA] rounded-full overflow-hidden">
    <div class="h-full bg-charcoal transition-all duration-500" style="width: 28.4%;"></div>
  </div>
</div>
```

---

## 6. Desain Ulang Arsitektur Informasi & Tata Letak Antarmuka

### 6.1 Halaman Beranda Penyamaran (`/` - Decoy Research Portal)
- **Macro-spacing**: `py-20 md:py-32` dengan `max-w-5xl mx-auto px-4 sm:px-6`.
- **Editorial Hero**:
  - Tanggal/Volume: `text-xs font-mono uppercase tracking-badge text-secondary mb-3`.
  - Judul Display: `font-serif text-3xl md:text-5xl font-medium tracking-tight-title leading-title text-charcoal mb-4`.
  - Lead Text: `text-base md:text-lg text-secondary leading-body max-w-3xl mb-8`.
  - CTAs: Sejajar horisontal dengan tombol trial instan (`Masuk Portal &rarr;` dan `Coba Gratis 24 Jam`).
- **Bento Grid Penelitian & Telemetri**:
  - Grid 12 kolom (`grid grid-cols-1 md:grid-cols-12 gap-6`).
  - Kolom 1-7: Kartu artikel riset terbaru dengan format dokumen editorial.
  - Kolom 8-12: Registri topologi node aktif (SG, JP, US, ID) lengkap dengan indikator latensi dan status normal.

### 6.2 Halaman Autentikasi Mandiri (`/portal` - Guest / Member Gate)
- Terletak terpusat di layar (`min-h-[80vh] flex items-center justify-center py-12 px-4`).
- Kontainer kartu tunggal `max-w-md w-full bg-surface border border-border rounded-card p-8 shadow-diffuse`.
- Tab peralihan atau pembagian vertikal yang jernih:
  - Bagian Atas: Form Input Token Akses Member + Tombol Masuk.
  - Pembatas: Garis tipis dengan label tengah `border-t border-border my-6 relative flex justify-center text-xs text-secondary`.
  - Bagian Bawah: Opsi Akun Uji Coba Tamu 24 Jam (1-Klik) dengan batas kecepatan 2 Mbps.

### 6.3 Dasbor Pengguna Aktif (Bento Grid 4 Kuadran)
Setelah terotentikasi, tampilan dasbor mengutamakan alur kerja pengguna:
```
+-----------------------------------------------------------------------------------+
| HEADER: Sesi Aktif [ LATINA01 ]  |  Masa Aktif [ 28 Hari Lagi ]  |  [ Tombol Log Out ]|
+-----------------------------------------------------------------------------------+
| KUADRAN 1: SETUP KLIEN 1-KLIK (col-span-12)                                       |
| - Tab Navigasi Klien: Clash Meta | Sing-box | Shadowrocket | v2rayNG             |
| - Input Tautan Langganan + Tombol Salin Cepat                                     |
| - Tombol Buka di Aplikasi Langsung (Deep Link)                                    |
| - Generator QR Code Responsif untuk Scan Ponsel Cepat                             |
+-------------------------------------------------+---------------------------------+
| KUADRAN 2: RUTE & NODE (col-span-12 lg:col-6)   | KUADRAN 3: KEAMANAN (lg:col-6)  |
| - Dropdown Pilihan Protokol (VMess/VLESS/Trojan)| - Sakelar Toggle AdBlock DNS    |
| - Dropdown Node Edge Cluster                    | - Tombol Reset UUID Baru        |
| - Dropdown Relay Multi-Hop                      | - Tombol Rotasi Token Akses     |
| - Tombol Simpan Konfigurasi & Terapkan          |                                 |
+-------------------------------------------------+---------------------------------+
| KUADRAN 4: TELEMETRI & SNI BUG HOST (col-span-12)                                 |
| - 3 Batang Meter Linear Ramping: CPU (%), RAM (%), Disk (%)                       |
| - Tabel Registri Wildcard SNI Bug Host dengan Tombol Salin per Host               |
+-----------------------------------------------------------------------------------+
```

---

## 7. Pipeline Tooling & Integrasi Teknis (Tailwind + Hugo)

Pengembangan tidak memerlukan runtime server Node.js di lingkungan produksi. Tailwind CSS diintegrasikan secara bersih melalui salah satu dari dua pendekatan:

### 7.1 Pendekatan Standalone Tailwind CLI (Direkomendasikan)
Menggunakan binary resmi Tailwind CLI tanpa dependensi npm global:
```bash
# Watch mode saat development
./tailwindcss -i ./web/assets/css/input.css -o ./web/static/css/styles.css --watch

# Build minify untuk production
./tailwindcss -i ./web/assets/css/input.css -o ./web/static/css/styles.css --minify
```

### 7.2 Struktur Berkas Input CSS (`web/assets/css/input.css`)
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  body {
    background-color: theme('colors.canvas.DEFAULT');
    color: theme('colors.charcoal.DEFAULT');
    font-family: theme('fontFamily.sans');
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
}
```

### 7.3 Penautan pada Template Hugo (`web/layouts/_default/baseof.html`)
Berkas output CSS disajikan secara langsung oleh webserver atau melalui pipa Hugo Pipes:
```html
<link rel="stylesheet" href="/css/styles.css">
```

---

## 8. Daftar Periksa Kualitas Eksekusi (Tailwind Quality Checklist)

Sebelum rilis atau setiap pengubahan kode antarmuka, verifikasi hal-hal berikut:
- [ ] Tidak ada kelas font bawaan generik (`font-sans` telah diarahkan ke `Geist Sans`/`SF Pro`).
- [ ] Tidak ada kelas bayangan tebal (`shadow-md`, `shadow-lg`, `shadow-xl`).
- [ ] Tidak ada warna saturasi tinggi bawaan Tailwind (`bg-blue-600`, `bg-emerald-500`, dll.).
- [ ] Tidak ada sudut `rounded-full` pada kontainer kartu atau tombol utama.
- [ ] Tidak ada emoji di seluruh kode markup, teks halaman, judul, atau skrip.
- [ ] Semua kartu bento menggunakan kelas `border border-[#EAEAEA]` dengan latar belakang putih.
- [ ] Tombol utama menggunakan kelas `bg-[#111111] text-white hover:bg-[#2F3437]`.
- [ ] Seluruh fungsi interaktif (login token, buat trial 24 jam, salin link langganan, scan QR code, ganti protokol/node/relay, toggle adblock, reset UUID, rotasi token, dan telemetri) bekerja 100% tanpa kendala.
