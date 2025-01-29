#!/bin/sh

root=$(pwd)
#rm -r $root/bin

cd $root/signaling_server/src
go build -o ../../bin/signaling_server/signaling_server main.go

cd $root/website_server/src
go build -o ../../bin/website_server/website_server main.go
cp -r templates ../../bin/website_server/
