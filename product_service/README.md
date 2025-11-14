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
- Dependency injection for repo/client wiring and singleton DB config for stable infra

## Architecture & Patterns

- Layered structure:

  - `internal/api/handlers` — gRPC handlers (translate RPC → service layer)
  - `internal/services` — business logic and orchestration
  - `internal/repo` — persistence (DynamoDB) and migrations
  - `internal/client` — gRPC clients to other microservices (User Service)
  - `internal/infra` — infra helpers (DynamoDB configuration, connection helpers)
  - `proto/` — protobuf service definitions and generated code

- Dependency Injection: constructors like `NewRepo`, `NewService`, `NewHandler` are used during bootstrap (`internal/bootstrap/server.go`) to wire components. This enables easy testing by providing mocked implementations.

- Singleton pattern for DB configuration: `internal/infra/db/db.go` uses a package-level `sync.Once` and a `GetDBConfig()` function that initializes and returns an `aws.Config` exactly once. This ensures a single, thread-safe DynamoDB client configuration is shared across the service:

```go
var (
    once sync.Once
    cfg  aws.Config
)

func loadDBConfig() { /* loads config and sets cfg */ }

func GetDBConfig() aws.Config {
    once.Do(loadDBConfig)
    return cfg
}
```

This pattern prevents redundant SDK configuration and makes the infra deterministic for tests and runtime.

## Repo layout (important paths)

- `cmd/api/main.go` — service entrypoint
- `internal/bootstrap/server.go` — app wiring (clients, repo, service, gRPC registration)
- `internal/api/handlers/product/handler.go` — gRPC handlers
- `internal/services/productService/service.go` — business logic
- `internal/repo/productRepo/repo.go` — DynamoDB data access
- `internal/client/product/client.go` — any external client usage (e.g., user service)
- `internal/infra/db/db.go` — DynamoDB SDK configuration
- `migrations/` — table initialization helpers (called from bootstrap)
- `proto/product.proto` — protobuf definitions
- `proto/gen/` — generated Go protobuf code (the repo includes generated files for convenience)

## Environment

Create a `.env` in `product_service/` or provide the environment variables through your system. Example:

```env
VERSION=1
ADDR=":5003"
SERVICE_NAME=product_service
USER_SERVICE_ADDR=0.0.0.0:5001
```

Notes:

- Do NOT leave spaces around `=` (e.g., `ADDR = ":5003"`), `godotenv`/Make include expects `ADDR=":5003"`.
- `ADDR` is the gRPC server listen address.
- `USER_SERVICE_ADDR` must point to a running user service; the product service attempts to dial it during startup.

## Dependencies

- Go 1.20+
- buf (for protobuf generation)
- DynamoDB Local (for local dev) or real AWS credentials for production

## DynamoDB Local (local development)

Start DynamoDB Local (expected endpoint `http://localhost:9000`):

```bash
docker run -d -p 9000:8000 amazon/dynamodb-local
```

The service's DB loader uses static dummy credentials and explicit endpoint resolver in `internal/infra/db/db.go` so local development works without AWS credentials.

During bootstrap the migration helper `migrations.InitProductTable(client)` is invoked to create the `Products` table if needed.

## Build & Run

Install modules and tidy:

```bash
go mod tidy
```

Generate protobuf (if you changed proto):

```bash
cd proto
buf generate
cd ..
```

Build the binary:

```bash
go build -o bin/main ./cmd/api
```

Run the service:

```bash
go run ./cmd/api
# or
./bin/main
```

Logs: bootstrap prints messages when DynamoDB client is created and when gRPC server starts listening.

## Protobuf / gRPC API

Service proto: `proto/product.proto`. Key RPCs (each returns `StandardResponse` wrapper):

- `CreateProduct(CreateProductRequest) returns (StandardResponse)`
- `GetProduct(GetProductsRequest) returns (StandardResponse)`
- `GetProductById(GetProductByIdRequest) returns (StandardResponse)`
- `UpdateProduct(UpdateProductRequest) returns (StandardResponse)`
- `DeleteProduct(DeleteProductRequest) returns (StandardResponse)`

`StandardResponse` provides `success`, `message`, `status_code` and a `oneof result` for typed payloads.

When exposing HTTP endpoints through Kong (grpc-gateway plugin) or Envoy, configure the gateway with `product.proto` so it can convert JSON → Protobuf.

### Example Kong config snippet (gateway must have the proto available inside the container)

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

### Example JSON → Kong → gRPC call (add product)

Request to Kong (HTTP):

```http
POST http://localhost:8000/products
Content-Type: application/json

{
  "name": "T-Shirt",
  "description": "100% cotton",
  "category": "apparel",
  "price": 19.99
}
```

Kong (grpc-gateway) converts JSON to Protobuf and forwards to the product gRPC server.

## Direct gRPC client example (Go)

```go
conn, err := grpc.DialContext(ctx, ":5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
defer conn.Close()
client := productpb.NewProductServiceClient(conn)
resp, err := client.CreateProduct(ctx, &productpb.CreateProductRequest{...})
```

## Troubleshooting

- Startup fails with "missing core service environment variables": verify `.env` and that variables are present and spelled as expected.
- `dial user service: context canceled` — product service attempts a blocking dial to `USER_SERVICE_ADDR`. Start `user_service` first or change the address.
- DynamoDB errors — ensure DynamoDB Local is running on `http://localhost:9000`.
- `buf generate` errors — ensure `buf` is installed and `proto/buf.yaml` dependencies are reachable (googleapis/protoc-gen-validate deps). If you generate in a container, ensure network access.

## Tests

There are no tests included by default. Recommended additions:

- Unit tests for `internal/services/productService` (table-driven tests)
- Integration tests using a DynamoDB Local instance and test data

Run unit tests:

```bash
go test ./internal/services/...
```

## Deployment notes

- In production, remove the local endpoint resolver in `internal/infra/db/db.go` and rely on real AWS credentials (or IAM roles).
- Avoid blocking `grpc.Dial` calls at startup in production; use backoff/retry or health-checks so the service can start and recover if dependencies are temporarily unavailable.

## Contribution and extension

- To add a new RPC: update `proto/product.proto`, run `buf generate`, implement handler -> service -> repo changes, wire in `internal/bootstrap/server.go` and test locally.
- Follow the pattern: handler (minimal) → service (business logic) → repo (persistence). Use DI via constructors (`NewRepo`, `NewService`, `NewHandler`).

## Where to look in code

- `internal/services/productService/service.go` — business logic and response building
- `internal/repo/productRepo/repo.go` — how products are persisted and queries are done
- `internal/bootstrap/server.go` — wiring (user client, repo, service, handler)

---

If you want, I can add a Dockerfile + docker-compose snippet that brings up DynamoDB Local + the product service and a Kong gateway example for local end-to-end testing.
