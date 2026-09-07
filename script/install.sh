#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
PROJECT="$(cd "$DIR/.." >/dev/null 2>&1 && pwd)"
GITHUB_TOKEN="${GITHUB_TOKEN}"
ARCHIVE_NAME="latinaserver.tar.gz"

if [ -n "$GITHUB_TOKEN" ] && [ -z "$GH_TOKEN" ]; then
    export GH_TOKEN="$GITHUB_TOKEN"
fi

TAG="$(gh release list --repo LalatinaHub/LatinaServer --json tagName -q ".[0].tagName")"
gh release download -O "$PROJECT/$ARCHIVE_NAME" --repo LalatinaHub/LatinaServer "$TAG" -p "*.gz"
tar -xzf "$PROJECT/$ARCHIVE_NAME" -C "$PROJECT/"

sudo rm -rf /var/www/mipa
sudo mkdir -p /var/www/mipa
if [ -d "$PROJECT/web/dist" ]; then
    sudo cp -r "$PROJECT/web/dist/." /var/www/mipa/
elif [ -d "./LatinaServer/web/dist" ]; then
    sudo cp -r "./LatinaServer/web/dist/." /var/www/mipa/
elif [ -d "/root/LatinaServer/web/dist" ]; then
    sudo cp -r "/root/LatinaServer/web/dist/." /var/www/mipa/
fi

sudo mkdir -p /usr/local/etc/latinaserver
if [ -f "$PROJECT/resources/caddy/caddy.json" ]; then
    sudo cp "$PROJECT/resources/caddy/caddy.json" /usr/local/etc/latinaserver/
elif [ -f "./LatinaServer/resources/caddy/caddy.json" ]; then
    sudo cp "./LatinaServer/resources/caddy/caddy.json" /usr/local/etc/latinaserver/
elif [ -f "/root/LatinaServer/resources/caddy/caddy.json" ]; then
    sudo cp "/root/LatinaServer/resources/caddy/caddy.json" /usr/local/etc/latinaserver/
fi

if [ -f "$PROJECT/resources/sing-box/config.json" ]; then
    sudo cp "$PROJECT/resources/sing-box/config.json" /usr/local/etc/latinaserver/
elif [ -f "./LatinaServer/resources/sing-box/config.json" ]; then
    sudo cp "./LatinaServer/resources/sing-box/config.json" /usr/local/etc/latinaserver/
elif [ -f "/root/LatinaServer/resources/sing-box/config.json" ]; then
    sudo cp "/root/LatinaServer/resources/sing-box/config.json" /usr/local/etc/latinaserver/
fi

if [ ! -f /etc/systemd/system/latinaserver.service ]; then
    if [ -f "$DIR/latinaserver.service" ]; then
        sudo cp "$DIR/latinaserver.service" /etc/systemd/system/
    elif [ -f "$PROJECT/script/latinaserver.service" ]; then
        sudo cp "$PROJECT/script/latinaserver.service" /etc/systemd/system/
    fi
else 
    sudo systemctl stop latinaserver
fi

sudo cp "$PROJECT/latinaserver" /usr/local/bin/
sudo chmod +x /usr/local/bin/latinaserver
sudo systemctl daemon-reload

sudo systemctl start latinaserver

sudo rm -f "$PROJECT/$ARCHIVE_NAME"
sudo rm -f "$PROJECT/latinaserver"