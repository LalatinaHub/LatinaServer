#!/usr/bin/env bash

DIR=$(dirname "$0")
PROJECT=$DIR/..
GITHUB_TOKEN="${GITHUB_TOKEN}"
ARCHIVE_NAME="latinaserver.tar.gz"

LATINASERVER_DOWNLOAD_URL=$(curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/LalatinaHub/LatinaServer/releases | jq -r ".[0].assets[0].browser_download_url")

curl -o $PROJECT/$ARCHIVE_NAME  -H "Authorization: token $GITHUB_TOKEN" $LATINASERVER_DOWNLOAD_URL
tar -xzf $PROJECT/$ARCHIVE_NAME

sudo mkdir -p /usr/local/etc/latinaserver
sudo cp $PROJECT/resources/openresty/nginx.conf /usr/local/etc/latinaserver/
sudo cp $PROJECT/config.json /usr/local/etc/latinaserver/

if [ ! -f /etc/systemd/system/latinaserver.service ]; then
    sudo cp ./latinaserver.service /etc/systemd/system/
else 
    sudo systemctl stop latinaserver
fi
sudo cp $PROJECT/latinaserver /usr/local/bin/
sudo systemctl daemon-reload

sudo systemctl start latinaserver

sudo rm -rf $PROJECT/$ARCHIVE_NAME
sudo rm -rf $PROJECT/latinaserver