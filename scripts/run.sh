#!/bin/bash

# This scripts allows fast local runs and re-runs of gostream

container_running() {
  docker ps --filter "name=$1" --filter "status=running" | grep -q "$1"
}

if ! container_running "stream_db" || ! container_running "stream_redis"; then
  docker-compose up postgres redis -d
fi

make build > /dev/null
./bin/gostream "$@"
