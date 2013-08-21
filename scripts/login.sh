#!/bin/sh

echo "login"

HOST="http://shortfuse.io"

curl -XPOST $HOST/user/login \
  -H "Content-Type: application/json" \
  -d '{ "name":"shredder", "password":"bobafett" }'
