#!/usr/bin/env bash

DIR=$(dirname "$0")
PROJECT=$DIR/..
GITHUB_TOKEN="${GITHUB_TOKEN}"
ARCHIVE_NAME="latinaserver.tar.gz"

echo $GITHUB_TOKEN | gh auth login --with-token
gh release download -O $PROJECT/$ARCHIVE_NAME --repo LalatinaHub/LatinaServer $(gh release list --repo LalatinaHub/LatinaServer --json tagName -q ".[0].tagName") -p "*.gz"
tar -xzf $PROJECT/$ARCHIVE_NAME -C $PROJECT/

sudo mkdir -p /usr/local/etc/latinaserver
sudo cp $PROJECT/resources/caddy/caddy.json /usr/local/etc/latinaserver/
sudo cp $PROJECT/config.json /usr/local/etc/latinaserver/

if [ ! -f /etc/systemd/system/latinaserver.service ]; then
    sudo cp ./latinaserver.service /etc/systemd/system/
else 
    sudo systemctl stop latinaserver
fi
sudo cp $PROJECT/latinaserver /usr/local/bin/
sudo systemctl daemon-reload

sudo systemctl start latinaserver

sudo rm -rf $PROJECT/*.gz
sudo rm -rf $PROJECT/latinaserver