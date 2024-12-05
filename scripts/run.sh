#!/bin/bash

container_running() {
  docker ps --filter "name=$1" --filter "status=running" | grep -q "$1"
}

if ! container_running "stream_db" || ! container_running "stream_redis"; then
  docker-compose up -d
fi

make build > /dev/null
./bin/gostream "$@"
