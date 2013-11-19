#!/bin/sh

cd web

echo "...resetting log file"
rm -f log/web.log

echo "...fetching dependencies"
go get -v

echo "running webserver"
echo "================="
go run *.go
