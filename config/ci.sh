#!/bin/sh

echo "Dropping the CI like a bunch of motherfuckers!"

if [ -n "$GO_ENV" ] && [ "$GO_ENV" == "production" ]; then
  echo "set to production"
  cd /go/src/github.com/garethstokes/fourtyeight
fi

echo "reseting HEAD --hard"
echo "$(git reset HEAD --hard)"

LOCAL="$(git rev-parse master)"
REMOTE="$(git rev-parse origin/master)"

echo "$LOCAL"
echo "$REMOTE"

if [ "$LOCAL" != "$REMOTE" ]; then
  echo "updating repository"
  echo "$(git pull origin master)"

  echo "restarting server"
  echo "$(god restart webserver)"
fi

echo "local: $(git rev-parse master)"
echo "Dropped out, yo"
