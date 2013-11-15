#!/bin/sh

echo "changing dir to project root"
cd ../

LOCAL="$(git rev-parse master)"
REMOTE="$(git rev-parse origin/master)"

echo $LOCAL
echo $REMOTE

if [ "$LOCAL" != "$REMOTE" ]; then
  echo "updating repository"
  git pull origin master

  echo "restarting server"
  god restart webserver
fi
