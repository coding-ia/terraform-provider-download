#!/bin/bash
set -e

docker-compose -f ./scripts/docker-compose.yaml down
rm ./scripts/files -rf
