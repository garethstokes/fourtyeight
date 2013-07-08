#!/bin/sh

cd web
echo "...fetching dependencies"
go get -v

echo "running webserver"
echo "================="
go run *.go
