#!/bin/bash
set -e

mkdir ./scripts/files
dd if=/dev/zero of=./scripts/files/file.dat bs=1024 count=2048
docker-compose -f ./scripts/docker-compose.yaml up -d
