#!/bin/env bash

dir="$(mktemp -d)"
name=$(basename "$dir")
go build -o "$dir/$name" ./cmd/bingo/main.go
BINGO_DB_DIR=:memory: "$dir/$name" &
PID=$(pidof $name)

source .env
# wait for app to start
until [ "$(curl -s -w '%{http_code}' -o /dev/null "http://localhost$BINGO_ADDR/")" -eq 200 ]
do
    sleep 0.1
done

hurl --test ./hurl
#
kill -9 "$PID"
rm -r "$dir"
