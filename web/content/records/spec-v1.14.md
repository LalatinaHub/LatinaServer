---
title: "Spesifikasi Formal Arsitektur Node Cluster v1.14 (Edge Ingress, Egress Routing & Telemetry Registry)"
date: 2026-08-28T09:00:00+07:00
author: "Lalatina Architecture Review Board"
description: "Dokumen spesifikasi formal node cluster v1.14 yang mendefinisikan matriks alokasi port, parameter kriptografi TLS 1.3, skema autentikasi dual-tier, dan format ekspor langganan."
tags: ["SPECIFICATION", "NODE-CLUSTER", "STANDARDS", "CRYPTOGRAPHY", "SING-BOX", "RFC"]
spec_id: "SPEC-V1.14"
draft: false
---

## 1. Ikhtisar & Ruang Lingkup Dokumen

Dokumen spesifikasi ini menetapkan standar arsitektur formal untuk seluruh node operasional yang tergabung dalam kluster jaringan terdistribusi Lalatina Systems. Standar ini mencakup konvensi perutean lapisan masuk (*ingress routing*), negosiasi parameter keamanan *Transport Layer Security* (TLS), isolasi kompartemen memori, dan format interoperabilitas profil konfigurasi klien.

---

## 2. Matriks Alokasi Port & Arsitektur Ingress

Setiap node cluster menjalankan stack terisolasi dengan pembagian tanggung jawab layanan yang tegas:

| Port Jaringan | Protokol Masuk | Daemon Penanggung Jawab | Deskripsi Fungsi |
| :--- | :--- | :--- | :--- |
| **Port 443** (TCP/UDP) | HTTPS / HTTP/3 | **Caddy Reverse Proxy v2.8+** | Gerbang utama internet publik, terminasi TLS resmi, routing SNI dan URL path |
| **Port 80** (TCP) | HTTP/1.1 | **Caddy Reverse Proxy** | Otomatisasi sertifikasi ACME (Let's Encrypt / ZeroSSL) & pengalihan 301 ke HTTPS |
| **Port 8080** (Internal) | HTTP/2 Cleartext | **LatinaServer Go Backend** | Layanan API `/api/v1/*`, User Settings Handler, Telemetri, dan Portal Pengguna |
| **Port 10086** (Loopback) | VMess-WebSocket | **sing-box Routing Core** | Soket internal untuk tunnel VMess terenkapsulasi WebSocket |
| **Port 10087** (Loopback) | VLESS-gRPC | **sing-box Routing Core** | Soket internal untuk tunnel VLESS terenkapsulasi gRPC |
| **Port 10088** (Loopback) | Trojan-TCP | **sing-box Routing Core** | Soket internal untuk tunnel Trojan dengan SNI camouflage |
| **Port 20001** (WireGuard)| UDP Mesh | **Cloudflare WARP Interface** | Terowongan antarmuka egress untuk masking IP keluar node |

```
                              [INTERNET PUBLIK]
                                      │
                               (Port 443 / HTTPS)
                                      ▼
                        ┌───────────────────────────┐
                        │   Caddy Webserver Edge    │
                        │    (TLS 1.3 Auto ACME)    │
                        └─────────────┬─────────────┘
                                      │
          ┌───────────────────────────┼───────────────────────────┐
          │ (Path: /api/v1/*, /portal)│ (Path: /ws-vmess)         │ (Path: /grpc-vless)
          ▼                           ▼                           ▼
┌──────────────────┐        ┌──────────────────┐        ┌──────────────────┐
│ LatinaServer Go  │        │  sing-box:10086  │        │  sing-box:10087  │
│ Web Service API  │        │  (VMess Inbound) │        │  (VLESS Inbound) │
└──────────────────┘        └──────────────────┘        └──────────────────┘
```

---

## 3. Profil Kriptografis & Standar Keamanan

### 3.1. TLS 1.3 Handshake Configuration
Untuk mencegah mitigasi penurunan keamanan (*downgrade attacks*) dan menjamin *Forward Secrecy* sempurna:
- **Protokol Minimum**: TLS 1.3 (Dukungan fallback ke TLS 1.2 dibatasi hanya dengan ciphersuite AEAD).
- **Ciphersuites Resmi**:
  - `TLS_AES_128_GCM_SHA256` (Direkomendasikan untuk perangkat keras dengan akselerasi AES-NI).
  - `TLS_CHACHA20_POLY1305_SHA256` (Dioptimalkan untuk perangkat seluler berbasis ARM/Android).
- **Kurva Elliptic Curve (ECC)**: X25519, Secp256r1.
- **Application-Layer Protocol Negotiation (ALPN)**: `h2`, `http/1.1`.

---

## 4. Skema Akun & Kebijakan Autentikasi Pengguna

Sistem membedakan dua kelas hak akses dengan spesifikasi parameter terukur:

### 4.1. Akun Uji Coba (*Guest Trial Mode*)
- **Inisiasi**: Dihasilkan secara instan tanpa registrasi via endpoint `POST /api/v1/trial`.
- **Masa Berlaku**: Tepat 86.400 detik (24 Jam) sejak waktu pembangkitan.
- **Batas Kecepatan (*Bandwidth Throttling*)**: Maksimum 2 Mbps (250 KB/detik) simetris unduh dan unggah.
- **Protokol**: VMess-WebSocket terenkapsulasi TLS.

### 4.2. Anggota Terdaftar (*Member Tier*)
- **Inisiasi**: Token otentikasi unik 8-karakter (misal `LATINA01`, `CRUSADER`).
- **Masa Berlaku**: Sesuai masa aktif keanggotaan guild / lisensi.
- **Batas Kecepatan**: Tidak dibatasi (*Uncapped Uplink/Downlink* hingga 10 Gbps).
- **Akses Protokol**: Fleksibel beralih antara VMess, VLESS, dan Trojan.
- **Fitur Khusus**:
  - Reset UUID mandiri via Portal dengan pembaruan otomatis ke file konfigurasi core.
  - Rotasi token login baru sewaktu-waktu.
  - Pengaturan Perisai Ganda DNS & AdBlock (`toggle-adblock`).

---

## 5. Standar Format Ekspor Interoperabilitas

Sistem menyediakan antarmuka 1-Click Invocation untuk platform proksi terkemuka dunia:

1. **Clash Meta / Mihomo (YAML)**: Memuat definisi lengkap `proxies`, `proxy-groups`, dan `rules` dengan format nama node terstandar: `[SG] Lalatina Direct`, `[JP] Lalatina Relay`.
2. **sing-box Platform (JSON)**: Format universal modern yang mendukung *inbounds*, *outbounds* multiplexing, dan rule-set geosite native.
3. **Shadowrocket (iOS)**: Skema URI terenkripsi dengan parameter obfs WebSocket dan SNI camouflage.
4. **v2rayNG / NekoBox (Android & Desktop)**: Format Base64 link `vmess://`, `vless://`, dan `trojan://` siap pindai kode QR.

---

## 6. Jaminan Integritas & Auditabilitas (Zero-Log Retention)

Secara kontraktual dan teknis, node kluster v1.14 menjamin:
- Direktori `/var/log` pada server produksi dialihkan ke perangkat virtual `/dev/null` untuk log akses koneksi.
- Database lokal (SQLite/Turso) hanya menyimpan data metadata relasional (Token hash, UUID, kuota byte agregat, timestamp kedaluwarsa).
- Tidak ada kueri nama domain, riwayat IP pengguna, atau payload komunikasi yang pernah disimpan di media penyimpanan non-volatil.
