# Kong API Gateway Setup

This directory contains the Kong API Gateway configuration for the eCommerce microservices.

## Overview

Kong acts as an API Gateway that routes HTTP requests to gRPC backend services.

## Configuration Files

- **kong.tpl.yml** - Kong configuration template with placeholder for protocol buffer descriptor
- **descriptor.b64** - Base64-encoded Protocol Buffer descriptor file (generated from proto files)
- **generate-config.sh** - Script that generates the final kong.yml from the template
- **docker-compose.yml** - Docker Compose configuration to run Kong

## How It Works

1. When Kong container starts, `generate-config.sh` runs
2. The script reads `descriptor.b64` and replaces `${PROTO_B64}` in `kong.tpl.yml`
3. Generated config is written to `/tmp/kong.yml` inside the container
4. Kong starts with the generated configuration

## Starting Kong

```bash
cd kongApiGateway
docker compose up -d
```

## Checking Status

```bash
# Check if Kong is running
docker compose ps

# View Kong logs
docker compose logs kong

# Test Kong endpoint
curl http://localhost:8000/
```

## Configured Routes

- **POST /auth/create-account** - Routes to auth service at grpc://host.docker.internal:5002

## Prerequisites

Before Kong can successfully proxy requests, ensure:

1. **Auth Service** is running on port 5002
2. **User Service** is running on port 5001 (if configured)
3. Services are accessible via `host.docker.internal` from Docker container

## Troubleshooting

### Connection Refused on localhost:8000

- Check if Kong container is running: `docker compose ps`
- Check Kong logs: `docker compose logs kong`

### 500 Internal Server Error

This means Kong is running but the backend service (e.g., auth service on port 5002) is not available.

**Solution**: Start the backend services first:

```bash
# Start auth service
cd ../auth_service
make run

# Start user service (in another terminal)
cd ../user_service
npm run start:dev
```

### Permission Denied Errors

The generate-config.sh script writes to /tmp inside the container, which should always be writable. If you see permission errors, check the volume mounts in docker-compose.yml.

## Ports

- **8000** - Kong HTTP Proxy (main entry point)
- **8001** - Kong Admin API
- **8443** - Kong HTTPS Proxy

## Regenerating descriptor.b64

If you modify the proto files, regenerate the descriptor:

```bash
# This depends on your proto build setup
# Typically would be something like:
buf build -o descriptor.bin
base64 descriptor.bin > descriptor.b64
```
