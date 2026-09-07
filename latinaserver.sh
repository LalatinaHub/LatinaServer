#!/usr/bin/env bash
set -e

# ==============================================================================
# LatinaServer Deployment & Management Script
# ==============================================================================

TOKEN_FILE="$HOME/.latinatoken"
REPO_NAME="LalatinaHub/LatinaServer"
REPO_URL="https://github.com/${REPO_NAME}"
PROJECT_DIR="${LATINA_DIR:-$(pwd)/LatinaServer}"

# 1. Muat token yang tersimpan jika ada
if [ -z "$GH_TOKEN" ] && [ -f "$TOKEN_FILE" ]; then
    GH_TOKEN=$(cat "$TOKEN_FILE" | tr -d '\r\n ')
fi
if [ -z "$GH_TOKEN" ] && [ -n "$GITHUB_TOKEN" ]; then
    GH_TOKEN="$GITHUB_TOKEN"
fi
export GH_TOKEN
export GITHUB_TOKEN="$GH_TOKEN"

# 2. Fungsi validasi kredensial GitHub CLI
check_token() {
    if [ -n "$GH_TOKEN" ]; then
        if gh api user >/dev/null 2>&1; then
            return 0
        fi
    fi
    return 1
}

# 3. Minta token jika belum ada atau tidak valid (401 Bad credentials)
ensure_token() {
    if check_token; then
        return 0
    fi

    echo ""
    echo "==================================================================="
    echo " [!] Autentikasi GitHub Diperlukan (Private Repository)"
    echo "==================================================================="
    echo " Repositori ${REPO_NAME} bersifat privat dan memerlukan token akses."
    echo " Buat Personal Access Token (PAT) baru di:"
    echo " -> https://github.com/settings/tokens"
    echo " (Pilih 'Generate new token (classic)', centang izin 'repo')"
    echo "==================================================================="
    echo ""

    while true; do
        read -p "Masukkan GitHub Token (ghp_...): " input_token
        input_token=$(echo "$input_token" | tr -d '\r\n ')

        if [ -z "$input_token" ]; then
            echo "[!] Token tidak boleh kosong. Coba lagi."
            continue
        fi

        export GH_TOKEN="$input_token"
        export GITHUB_TOKEN="$input_token"

        echo "[+] Memverifikasi token ke GitHub API..."
        if check_token; then
            echo "$input_token" > "$TOKEN_FILE"
            chmod 600 "$TOKEN_FILE"
            echo "[✓] Token valid dan berhasil disimpan di $TOKEN_FILE."
            break
        else
            echo "[x] Error: Token tidak valid atau kedaluwarsa (401 Unauthorized)."
            echo "    Pastikan token memiliki hak akses 'repo' dan belum direvoke."
        fi
    done
}

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
        echo "[!] Menggunakan apt fallback untuk gh..."
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
    ensure_token

    echo "[+] Menyiapkan LatinaServer..."
    systemctl stop latinaserver 2>/dev/null || true

    # Sinkronisasi repository privat menggunakan gh atau authenticated git
    local auth_repo_url="https://${GH_TOKEN}@github.com/${REPO_NAME}.git"

    if [ -d "$PROJECT_DIR/.git" ]; then
        echo "[+] Memperbarui repository di $PROJECT_DIR..."
        cd "$PROJECT_DIR"
        git remote set-url origin "$auth_repo_url"
        git fetch origin main --prune
        git reset --hard origin/main
    else
        echo "[+] Meng-clone private repository..."
        rm -rf "$PROJECT_DIR"
        git clone "$auth_repo_url" "$PROJECT_DIR"
        cd "$PROJECT_DIR"
    fi

    # Pastikan script install executable dan jalankan
    if [ -f "$PROJECT_DIR/script/install.sh" ]; then
        chmod +x "$PROJECT_DIR/script/install.sh"
        echo "[+] Menjalankan script instalasi & unduh release..."
        cd "$PROJECT_DIR/script"
        export GH_TOKEN
        export GITHUB_TOKEN="$GH_TOKEN"
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
    echo " Pastikan kredensial database sudah terisi di:"
    echo " sudo nano /etc/systemd/system/latinaserver.service"
    echo " (Isi TURSO_DATABASE_URL dan TURSO_AUTH_TOKEN lalu restart)"
    echo "==================================================================="
}

reset_token() {
    rm -f "$TOKEN_FILE"
    unset GH_TOKEN
    unset GITHUB_TOKEN
    echo "[✓] Token lama telah dihapus."
    ensure_token
}

menu() {
    echo "=========================================="
    echo "     LatinaServer Deployment Utility      "
    echo "=========================================="
    echo "1. Install (Dependensi Lengkap + LatinaServer)"
    echo "2. Update (Perbarui LatinaServer Cepat)"
    echo "3. Cek Status Service"
    echo "4. Lihat Log Live (journalctl)"
    echo "5. Ganti / Perbarui GitHub Token"
    echo "6. Keluar"
    echo "=========================================="
    read -t 15 -p "Pilihan Anda (1-6) [default: 1 dalam 15s]: " choice || true

    choice=${choice:-1}

    case $choice in
        1)
            install_dependencies
            install_gh
            install_tcp_brutal
            deploy_latinaserver
            ;;
        2)
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
            reset_token
            ;;
        6)
            echo "Keluar..."
            exit 0
            ;;
        *)
            echo "Pilihan tidak valid."
            sleep 1
            menu
            ;;
    esac
}

menu
