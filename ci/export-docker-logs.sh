#!/usr/bin/env bash

set -eu -o pipefail

mkdir -p logs/

docker container ls -a --format '{{.ID}} {{.Names}}' | while read -r id name; do
    docker logs "$id" > "logs/docker-$name.log" 2>&1 || true
done