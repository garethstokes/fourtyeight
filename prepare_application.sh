#!/bin/sh

echo "creating library database"
cat << EOF > createdb.sql
CREATE DATABASE IF NOT EXISTS fourtyeight_development;
EOF
mysql -u root < createdb.sql
rm createdb.sql

echo "fetching libraries"
cd migrations
go get -v

echo "running migrations"
./reset.sh
