# Product Service

This is the Product microservice for the ecom_microservice project. It is a gRPC service (no HTTP frontend inside the service) that stores product data in DynamoDB (development uses DynamoDB Local) and exposes a protobuf-defined gRPC API. The service follows the common project pattern used in this repo: `internal/api` (handlers), `internal/services` (business logic), `internal/repo` (data access), `internal/client` (external service clients), `internal/infra` (infrastructure), and `proto` (protobuf definitions).

This README documents how to build, run, and debug the Product Service locally.

## Quick facts

- Language: Go
- gRPC/proto toolchain: buf + protoc (project already contains `proto/buf.yaml`)
- Persistent store (dev): DynamoDB Local (SDK v2)
- Auth: The service uses a User Service client to verify user/customer info; the address is read from env (`USER_SERVICE_ADDR`).

## Repo layout (important paths)

- `internal/bootstrap/server.go` — wiring (user client, repo, service, handler)

# Product Service

This is the Product microservice for the ecom_microservice project. It's a gRPC-based service that manages product data and persists it in DynamoDB (DynamoDB Local is used for development). The project follows the common microservice layout used across this repository and uses dependency injection for wiring components and a singleton pattern for DB configuration.

This README covers features, architecture, dependency injection details (including the singleton DB config pattern), build/run instructions, proto/gRPC API, Kong integration notes, troubleshooting, and next steps.

## Features

- Complete product lifecycle: create, list, fetch-by-id, update, delete
- Protobuf-defined gRPC API with a `StandardResponse` wrapper for consistent responses
- DynamoDB storage (development: DynamoDB Local)
- Uses a User Service client to validate user/customer-related operations
- Clear Handler → Service → Repo separation for testability and maintainability

# Product Service

A gRPC microservice that manages product data for the ecom_microservice project. It persists product records in DynamoDB (DynamoDB Local for development) and exposes a Protobuf-defined gRPC API. The service follows the repo-wide layering convention (handler → service → repo) and uses constructor-based dependency injection and a singleton DB config for infra.

## 🎯 Overview

The Product Service provides the product lifecycle: create, read (list & by-id), update, and delete. It is intended to be consumed by other microservices (cart, order) and by an HTTP gateway (Kong) that translates JSON requests into gRPC calls.

Key responsibilities:

- Store and query product data in DynamoDB
- Expose a protobuf/gRPC API for product operations
- Validate and shape responses using a `StandardResponse` wrapper
- Call User Service where product operations require user validation

## 🏗️ Technology & Stack

- Language: Go (1.20+)
- gRPC + Protocol Buffers (buf toolchain)
- Persistent store (dev): DynamoDB Local (AWS SDK v2)
- Gateway: Kong (grpc-gateway plugin) for JSON→gRPC translation

## 📁 Project Structure

```
product_service/
├── cmd/api/                    # Application entrypoint
│   └── main.go                 # Starts the server
├── internal/
│   ├── api/handlers/           # gRPC request handlers
│   │   └── product/handler.go  # Handler implementations
│   ├── bootstrap/              # Application initialization & DI
│   │   └── server.go
│   ├── client/                 # gRPC clients to other services
│   │   └── product/            # (if needed)
│   ├── config/                 # Configuration loader
│   ├── domain/                 # Domain models
│   ├── infra/                  # Infrastructure helpers (DynamoDB config)
│   │   └── db/db.go
│   ├── repo/                   # Persistence layer (DynamoDB)
│   │   └── productRepo/repo.go
│   └── services/               # Business logic
│       └── productService/service.go
├── migrations/                 # Table init helpers
├── proto/                      # Protobuf definitions & buf config
└── proto/gen/                  # Generated code (checked-in for convenience)
```

## ✨ Features

- Create, list, fetch-by-id, update, and delete products
- gRPC API with `StandardResponse` wrapper for consistent responses
- DynamoDB-backed persistence with migration helper for table creation
- Clear DI pattern: `NewRepo`, `NewService`, `NewHandler` used during bootstrap

## Architecture & Patterns

- Layered design: `handlers` translate RPC → `services` (business rules) → `repo` (persistence)
- Dependency injection in `internal/bootstrap/server.go` for testability
- Singleton DB configuration in `internal/infra/db/db.go` using `sync.Once` to ensure a single aws.Config / DynamoDB client instance

## Environment

Create a `.env` file in `product_service/` or export these variables in your environment. Example:

```env
VERSION=1
SERVICE_NAME=product_service
ADDR=":5003"
USER_SERVICE_ADDR=0.0.0.0:5001
```

Notes:

- Avoid spaces around `=` in `.env` (e.g., `ADDR = ":5003"` will break parsers)
- `ADDR` is the gRPC server listen address
- `USER_SERVICE_ADDR` should point to a running `user_service` (product service dials it at startup)

## DynamoDB Local (local development)

Run DynamoDB Local for development (expected endpoint `http://localhost:9000`):

```bash
docker run -d -p 9000:8000 amazon/dynamodb-local
```

The code in `internal/infra/db/db.go` uses an explicit endpoint resolver and dummy credentials so local development works without AWS credentials. During bootstrap the migration helper `migrations.InitProductTable` will create the `Products` table if missing.

## Build & Run

Install modules and tidy dependencies:

```bash
go mod tidy
```

Generate protobuf artifacts (only if you changed `.proto` files):

```bash
cd proto
buf generate
cd ..
```

Build:

```bash
go build -o bin/main ./cmd/api
```

Run:

```bash
go run ./cmd/api
# or
./bin/main
```

Logs from bootstrap include DynamoDB client creation and gRPC listen address.

## Protobuf / gRPC API

The service contract is defined in `proto/product.proto`. RPCs are wrapped with `StandardResponse` for a consistent response envelope. Example RPCs:

- `CreateProduct(CreateProductRequest) returns (StandardResponse)`
- `GetProducts(GetProductsRequest) returns (StandardResponse)`
- `GetProductById(GetProductByIdRequest) returns (StandardResponse)`
- `UpdateProduct(UpdateProductRequest) returns (StandardResponse)`
- `DeleteProduct(DeleteProductRequest) returns (StandardResponse)`

When exposing HTTP endpoints via Kong (or other gateways) enable the grpc-gateway plugin and provide `product.proto` inside the gateway container so JSON requests can be converted to Protobuf.

## Kong Example (gateway side)

Gateway must have access to `product.proto`. Example plugin snippet:

```yaml
plugins:
  - name: grpc-gateway
    config:
      proto: /kong/protos/product.proto

service:
  name: product-service
  url: grpc://host.docker.internal:5003
  routes:
    - name: product-create
      paths: ["/products"]
      methods: [POST]
```

## Example: JSON → Kong → gRPC (Add Product)

POST to Kong's HTTP endpoint with JSON body (Content-Type: application/json). Kong converts and forwards to gRPC. Example payload:

```json
{
  "name": "T-Shirt",
  "description": "100% cotton",
  "category": "apparel",
  "price": 19.99
}
```

## Troubleshooting

- Startup fails with `dial user service: context canceled`: ensure `user_service` is running and `USER_SERVICE_ADDR` is correct. The product service attempts a blocking dial to the user service during bootstrap.
- DynamoDB errors: ensure DynamoDB Local is running on `http://localhost:9000`.
- `buf generate` errors: install `buf` and ensure `proto/buf.yaml` dependencies are reachable.

## Tests

There are no unit/integration tests included by default. Recommended:

- Unit tests for `internal/services/productService`
- Integration tests using DynamoDB Local and a test table

Run unit tests:

```bash
go test ./internal/services/...
```

## Deployment notes

- In production remove local endpoint resolver and use real AWS credentials / IAM roles in `internal/infra/db/db.go`.
- Avoid blocking `grpc.Dial` calls in bootstrap for production; prefer retries/backoff or non-blocking health checks so the service can start even if dependencies are temporarily unavailable.

## Contribution & Extension

- To add a new RPC: update `proto/product.proto`, run `buf generate`, implement the handler → service → repo changes, wire in `internal/bootstrap/server.go`, and test locally.
- Follow the repository pattern: handler (thin) → service (business rules) → repo (persistence). Use constructor DI for each component.

## Where to look in code

- `internal/bootstrap/server.go` — wiring and DI
- `internal/api/handlers/product/handler.go` — gRPC handlers
- `internal/services/productService/service.go` — business logic
- `internal/repo/productRepo/repo.go` — DynamoDB persistence

---

If you'd like, I can also add a Docker Compose snippet that brings up DynamoDB Local plus the product service and a Kong gateway for local end-to-end testing.
