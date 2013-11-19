#!/bin/sh

if [[ -n "$GO_ENV" && $GO_ENV == "production" ]]; then
  echo "set to production"
  cd /go/src/github.com/garethstokes/fourtyeight
fi

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
