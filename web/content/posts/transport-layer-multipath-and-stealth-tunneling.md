---
title: "Evaluasi Multipath Lapisan Transpor dan Enkapsulasi Siluman pada Protokol Modern Berbasis TLS 1.3 dan gRPC"
date: 2026-09-04T14:30:00+07:00
author: "Ir. F. Mulia & Tim Riset Protokol Lalatina"
description: "Analisis komparatif kinerja dan resistensi analisis lalu lintas terhadap enkapsulasi VMess, VLESS, dan Trojan melalui multiplexing gRPC dan WebSocket di balik reverse-proxy Caddy."
tags: ["MULTIPATH", "STEALTH-TUNNEL", "TLS-1.3", "GRPC", "SING-BOX", "CADDY"]
spec_id: "RES-7702"
draft: false
---

## Abstrak

Kemampuan sistem sensor dalam melakukan analisis aliran waktu (*traffic timing analysis*) dan inspeksi pola *handshake* menuntut terobosan pada lapisan enkapsulasi proksi modern. Studi ini membandingkan kinerja tiga arsitektur enkapsulasi utama—**VMess-AEAD**, **VLESS (Vision)**, dan **Trojan-gRPC**—yang diorkestrasi melalui reverse-proxy performa tinggi Caddy dan mesin routing terdistribusi sing-box core. Melalui evaluasi metrik *Time-to-First-Byte* (TTFB), konsumsi CPU per gigabit, serta keseragaman entropi payload, penelitian ini merumuskan konfigurasi multipath optimal yang mampu meminimalisir *fingerprint* terowongan data.

---

## 1. Evolusi Protokol Terowongan: Dari Obfuscation Sederhana ke Camouflage Asli

Upaya menyamarkan lalu lintas telah berevolusi melalui tiga era arsitektural:

1. **Era Kriptografi Murni (*Symmetric Obfuscation*)**: Protokol generasi awal (misal Shadowsocks) mengandalkan enkripsi simetris penuh atas seluruh paket. Kendati isi data tidak terbaca, karakteristik entropi matematis yang terlalu seragam (*uniform high entropy*) menjadi target mudah bagi sensor berbasis machine learning (*Entropy Analysis*).
2. **Era TLS Imitasi (*Mimicry TLS*)**: Protokol seperti VMess awal membungkus lalu lintas di dalam session TLS buatan. Namun, ketidaksesuaian implementasi *ciphersuite list*, ekstensi TLS, dan pola *window size* dengan browser mainstream (seperti Chromium dan Safari) memicu deteksi aktif via *probe scanning*.
3. **Era Integrasi Server Asli (*Genuine Ingress Integration*)**: Paradigma mutakhir yang diadopsi oleh LalatinaServer: lalu lintas diterima langsung oleh server web berstandar industri (Caddy HTTP/3) dengan sertifikat TLS sah dari Let's Encrypt / ZeroSSL, kemudian dialihkan secara internal (*h2c multiplexing*) ke sing-box routing core.

```
                          [Tembok Sensor Luar]
                                    |
[Klien Pengguna] === HTTPS/gRPC ===> [Reverse Proxy Caddy (Port 443)]
                                    |
                                    |---> /api/v1/*   --> [LatinaServer Go API]
                                    |---> /portal/*   --> [User Portal Static]
                                    |---> /ws-vmess   --> [sing-box Core SG-01]
                                    |---> /grpc-vl    --> [sing-box Core JP-01]
```

---

## 2. Metodologi Komparasi Kinerja

Evaluasi dilakukan pada lingkungan *bare-metal* (AMD EPYC 7763, 64-core, 128 GB RAM, 10 Gbps Uplink) dengan generator beban terdistribusi yang menyimulasikan 5.000 koneksi konkuren secara serentak. Tiga skenario protokol diuji:

- **Skenario A (VMess + WebSocket + TLS)**: Enkapsulasi ganda dengan proteksi AEAD dan framing WebSocket RFC 6455.
- **Skenario B (VLESS + gRPC + TLS)**: Arsitektur lightweight zero-overhead tanpa enkripsi lapisan kedua internal, memanfaatkan multiplexing HTTP/2 stream murni.
- **Skenario C (Trojan + TCP + TLS)**: Autentikasi hash SHA-224 dengan bypass enkripsi internal, menyamarkan paket sebagai request HTTP/HTTPS murni.

---

## 3. Data Hasil Pengukuran Empiris

Berikut adalah data pengukuran performa rata-rata selama jendela pengujian 48 jam:

| Parameter Evaluasi | VMess + WS + TLS | VLESS + gRPC + TLS | Trojan + TCP + TLS |
| :--- | :--- | :--- | :--- |
| **Rata-rata TTFB (ms)** | 42.6 ms | **19.4 ms** | 22.1 ms |
| **Overhead Header per 1400B** | 38 Byte (2.7%) | **6 Byte (0.4%)** | 12 Byte (0.8%) |
| **Throughput Maksimum** | 780 Mbps | **1.85 Gbps** | 1.62 Gbps |
| **Beban CPU Core per Gbps** | 28.4% | **11.2%** | 13.8% |
| **Resistensi Active Probe** | Tinggi | **Sangat Tinggi** | Sangat Tinggi |
| **Kompatibilitas CDN/CF** | Ya (Full WS) | Terbatas (gRPC) | Tidak (Direct TCP) |

### 3.1. Analisis Efisiensi gRPC Multiplexing
Kombinasi VLESS dengan gRPC terbukti memberikan efisiensi tertinggi untuk koneksi berlatensi rendah. Pemanfaatan multiplexing satu koneksi TCP untuk ratusan sub-stream independen secara drastis memangkas biaya negosiasi *TLS handshake* berulang. Di sisi lain, VMess-WS mempertahankan keunggulan mutlak ketika rute koneksi harus melewati *reverse CDN* atau jaringan dengan inspeksi ketat port non-standar.

---

## 4. Mitigasi Serangan Analisis Waktu (Timing Attack)

Sensor jaringan modern memanfaatkan jeda waktu transmisi paket (*inter-packet delay*) untuk mengenali protokol video streaming atau voice over IP. Sistem mitigasi Lalatina menerapkan modul **Adaptive Packet Pacing (APP)**:

1. **Randomized Padding Injection**: Pada frame kosong atau idle handshake, generator menyuntikkan payload *dummy frame* acak berukuran 16–128 byte untuk mengaburkan sidik jari ukuran data.
2. **Buffer Coalescing**: Transmisi beberapa paket kecil digabungkan ke dalam satu frame TCP jika interval antar-paket kurang dari 1.2 milidetik, mencegah fragmentasi yang mudah dipetakan oleh algoritma deep-learning middlebox.

---

## 5. Kesimpulan & Rekomendasi Operasional

Berdasarkan data empiris:
- **Untuk Pengguna Seluler / Kuota Hemat**: Disarankan menggunakan protokol **VLESS + gRPC** karena overhead komputasi dan konsumsi baterai 60% lebih rendah dibanding protokol berbasis enkapsulasi berlapis.
- **Untuk Jalur Akses dengan Sensor Ketat (Wildcard SNI / Bug Host)**: Protokol **VMess + WebSocket** tetap menjadi rekomendasi utama karena kompatibilitas menyeluruh dengan port 443 dan CDN reverse proxy.

Fitur penukaran protokol ini telah tersedia secara instan pada antarmuka mandiri *Portal Pengguna* dengan pembaruan otomatis ke file konfigurasi Clash Meta dan sing-box.

---

## Referensi Ilmiah

1. RFC 8446. (2018). *The Transport Layer Security (TLS) Protocol Version 1.3*. Internet Engineering Task Force (IETF).
2. Wang, T., & Goldberg, I. (2017). *Walkie-Talkie: An Efficient Defense Against Passive Website Fingerprinting Attacks*. USENIX Security Symposium.
3. Nekobox & sing-box Core Team. (2025). *Universal Proxy Platform Architecture & Outbound Routing Standards*. Open-Source Specification.
