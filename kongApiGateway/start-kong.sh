#!/bin/bash

# Load .env properly (preserves quotes and spaces)
set -a
source .env
set +a

# Debug - show what we got
echo "JWT_ACCESS_SECRET=[$JWT_ACCESS_SECRET]"

# Stop Kong
docker compose down

# Delete old kong.yml
rm -rf kong.yml

# Generate kong.yml from template
envsubst < kong.tpl.yml > kong.yml

# Show generated file
echo "Generated kong.yml:"
cat kong.yml

# Start Kong
docker compose up