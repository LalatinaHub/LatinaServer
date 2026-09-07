#!/usr/bin/env bash
set -e

# ==============================================================================
# LatinaServer Deployment & Management Script
# ==============================================================================

export GITHUB_TOKEN="${GITHUB_TOKEN:-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}"
export GH_TOKEN="${GH_TOKEN:-$GITHUB_TOKEN}"

REPO_URL="https://github.com/LalatinaHub/LatinaServer"
PROJECT_DIR="${LATINA_DIR:-$(pwd)/LatinaServer}"

install_dependencies() {
    echo "[+] Memeriksa dependensi sistem..."
    local pkgs=(wget gnupg ca-certificates build-essential openssl cron socat curl vnstat cmake make g++ git jq tar)
    local to_install=()

    for pkg in "${pkgs[@]}"; do
        if ! dpkg -s "$pkg" >/dev/null 2>&1; then
            to_install+=("$pkg")
        fi
    done

    if [ ${#to_install[@]} -gt 0 ]; then
        echo "[+] Menginstal dependensi: ${to_install[*]}..."
        apt-get update -y
        apt-get install -y "${to_install[@]}"
    else
        echo "[✓] Semua dependensi sistem sudah terpasang."
    fi

    # Tambahkan path Go ke .bashrc jika belum ada
    if ! grep -q '/usr/local/go/bin' ~/.bashrc 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
}

install_gh() {
    if command -v gh >/dev/null 2>&1; then
        echo "[✓] GitHub CLI (gh) sudah terpasang: $(gh --version | head -n1)"
        return 0
    fi

    echo "[+] Memasang GitHub CLI (gh)..."
    local GH_VERSION
    GH_VERSION=$(curl -s "https://api.github.com/repos/cli/cli/releases/latest" | jq -r .tag_name | sed 's/v//')

    if [ -n "$GH_VERSION" ]; then
        wget -qO /tmp/gh.tar.gz "https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_amd64.tar.gz"
        tar -xzf /tmp/gh.tar.gz -C /tmp/
        mv "/tmp/gh_${GH_VERSION}_linux_amd64/bin/gh" /usr/local/bin/
        chmod +x /usr/local/bin/gh
        rm -rf "/tmp/gh_${GH_VERSION}_linux_amd64" /tmp/gh.tar.gz
        echo "[✓] GitHub CLI berhasil dipasang."
    else
        echo "[!] Gagal mendeteksi versi gh via API. Menggunakan apt fallback..."
        apt-get update -y && apt-get install -y gh
    fi
}

install_tcp_brutal() {
    if command -v brutalctl >/dev/null 2>&1; then
        echo "[✓] TCP Brutal sudah terpasang."
        return 0
    fi

    echo "[+] Memasang TCP Brutal..."
    bash <(curl -fsSL https://tcp.hy2.sh/)
}

deploy_latinaserver() {
    echo "[+] Menyiapkan LatinaServer..."

    # Hentikan service berjalan jika ada
    systemctl stop latinaserver 2>/dev/null || true

    # Siapkan direktori repo
    if [ -d "$PROJECT_DIR/.git" ]; then
        echo "[+] Memperbarui repository yang ada di $PROJECT_DIR..."
        cd "$PROJECT_DIR"
        git fetch --all --prune
        git reset --hard origin/main || git pull origin main
    else
        echo "[+] Meng-clone repository dari $REPO_URL..."
        rm -rf "$PROJECT_DIR"
        git clone "$REPO_URL" "$PROJECT_DIR"
        cd "$PROJECT_DIR"
    fi

    # Pastikan script install executable dan jalankan
    if [ -f "$PROJECT_DIR/script/install.sh" ]; then
        chmod +x "$PROJECT_DIR/script/install.sh"
        echo "[+] Menjalankan script install..."
        cd "$PROJECT_DIR/script"
        bash ./install.sh
        cd "$PROJECT_DIR"
    else
        echo "[x] Error: $PROJECT_DIR/script/install.sh tidak ditemukan!"
        exit 1
    fi

    # Pastikan izin akses direktori runtime sudah tepat
    sudo mkdir -p /var/www/mipa /usr/local/etc/latinaserver
    sudo chmod -R 755 /var/www/mipa /usr/local/etc/latinaserver

    # Validasi dan reload service systemd
    echo "[+] Memulai ulang service latinaserver..."
    systemctl daemon-reload
    systemctl enable latinaserver
    systemctl restart latinaserver

    echo ""
    echo "==================================================================="
    echo " [✓] LatinaServer berhasil dipasang / diperbarui!"
    echo "==================================================================="
    echo " Status Service: $(systemctl is-active latinaserver 2>/dev/null || echo 'unknown')"
    echo ""
    echo " Catatan penting:"
    echo " Jika ini instalasi baru, pastikan kredensial database sudah terisi:"
    echo " sudo nano /etc/systemd/system/latinaserver.service"
    echo " (Isi TURSO_DATABASE_URL dan TURSO_AUTH_TOKEN lalu reload)"
    echo "==================================================================="
}

menu() {
    echo "=========================================="
    echo "     LatinaServer Deployment Utility      "
    echo "=========================================="
    echo "1. Install (Dependensi Lengkap + LatinaServer)"
    echo "2. Update (Perbarui LatinaServer Cepat)"
    echo "3. Cek Status Service"
    echo "4. Lihat Log Live (journalctl)"
    echo "5. Keluar"
    echo "=========================================="
    read -t 10 -p "Pilihan Anda (1/2/3/4/5) [default: 1 dalam 10s]: " choice || true

    choice=${choice:-1}

    case $choice in
        1)
            install_dependencies
            install_gh
            install_tcp_brutal
            deploy_latinaserver
            ;;
        2)
            # Pastikan dependensi minimal terpasang tanpa update apt lama
            install_dependencies
            install_gh
            deploy_latinaserver
            ;;
        3)
            systemctl status latinaserver --no-pager
            ;;
        4)
            journalctl -u latinaserver --output cat -f
            ;;
        5)
            echo "Keluar..."
            exit 0
            ;;
        *)
            echo "Pilihan tidak valid. Silakan pilih 1 - 5."
            sleep 1
            menu
            ;;
    esac
}

menu
