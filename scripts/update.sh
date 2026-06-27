#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "=== Chele Update ==="
echo ""

# 1. Pull latest from git
echo ">>> Pulling latest changes..."
GIT_OUTPUT=$(git pull 2>&1)
echo "$GIT_OUTPUT"
echo ""

# 2. Rebuild and restart
echo ">>> Rebuilding images..."
docker compose build 2>&1
echo ""

echo ">>> Restarting containers..."
docker compose up -d 2>&1
echo ""

# 3. Prune unused images
echo ">>> Cleaning up unused images..."
docker image prune -f 2>&1
echo ""

echo "=== Done ==="
