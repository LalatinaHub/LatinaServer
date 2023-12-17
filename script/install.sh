#!/usr/bin/env bash

DIR=$(dirname "$0")
PROJECT=$DIR/..

go build -tags with_grpc,with_clash_api,with_v2ray_api,with_wireguard,with_utls,with_reality_server,with_gvisor,with_quic -o $PROJECT/latinaserver $PROJECT/cmd/latinaserver/main.go

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
