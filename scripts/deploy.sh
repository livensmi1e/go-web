#!/bin/bash

set -e

IMAGE="ghcr.io/livensmi1e/go-web:latest"
echo "==> Pulling latest image: $IMAGE"
docker pull $IMAGE

echo "==> Starting services with docker-compose"
cd ~/app
docker compose --env-file .env.prod -f docker-compose.prod.yml up -d --remove-orphans