#!/bin/sh

root=$(pwd)
# rm -r $root/bin

cd $root/signaling_server/src
go build -o ../../bin/signaling_server/signaling_server main.go

cd $root/website_server/src
go build -o ../../bin/website_server/website_server main.go
cp -r templates ../../bin/website_server/

cd $root/pusher/src
apt install gcc-arm-linux-gnueabi
export GOOS=android
export GOARCH=arm64
go build -o ../../bin/pusher/pusher_android_arm64 main.go

export GOOS=windows
export GOARCH=amd64
go build -o ../../bin/pusher/pusher_windows_amd64.exe main.go
