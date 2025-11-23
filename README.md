# E-Commerce Microservices Monorepo

A collection of gRPC (Go + TypeScript) microservices and an API Gateway (Kong) forming the backend foundation of an e-commerce platform. This root README gives a high-level view. Individual services have their own detailed README files.

> Status: Auth, Product, Cart, User, and Kong Gateway services have initial documentation. Order Service is partially scaffolded (no README yet). Payment Service planned.

---

## 🧩 Services Overview

| Service          | Language          | Purpose                                                       | Storage              | External Deps                    | README                      |
| ---------------- | ----------------- | ------------------------------------------------------------- | -------------------- | -------------------------------- | --------------------------- |
| Auth Service     | Go                | Authentication (JWT access + refresh), multi-device sessions  | MongoDB + Redis      | User Service                     | `auth_service/README.md`    |
| User Service     | TypeScript (Node) | User/customer CRUD & credentials                              | MongoDB              | (none)                           | `user_service/README.md`    |
| Product Service  | Go                | Product CRUD & DynamoDB persistence                           | DynamoDB Local (dev) | User Service (validation)        | `product_service/README.md` |
| Cart Service     | Go                | User cart management (Redis JSON blob)                        | Redis                | Product Service                  | `cart_service/README.md`    |
| Kong API Gateway | Kong + Lua        | HTTP entrypoint, auth/metadata plugins, JSON→gRPC translation | (stateless)          | All gRPC services                | `kongApiGateway/README.md`  |
| Order Service    | Go                | (In progress) Order creation, validation pipeline (async)     | Postgres (planned)   | Product (validation), User, Cart | (pending)                   |
| Payment Service  | Go (planned)      | Payment authorization/capture, refunds                        | TBD                  | Order Service                    | (future)                    |

---

## 🔐 Authentication Flow (High-Level)

1. Client hits Kong (HTTP) with credentials (cookies/JWT).
2. Lua plugins validate token & inject metadata (e.g., `x-user-email`).
3. Auth Service handles login/refresh; issues tokens (access + refresh); stores refresh tokens (Redis + MongoDB).
4. Downstream services rely on metadata for authorization/context.

---

## 🛒 Cart & Product Interaction

- Cart Service calls Product Service gRPC to fetch product data (name/price) before persisting item.
- Cart persists entire cart as one Redis JSON blob keyed by `cart:{email}` with TTL.
- Subtotal recalculated on every cart mutation (simple deterministic logic).

---

## 📦 Product Service Architecture

- DynamoDB Local for development (singleton AWS config).
- CRUD RPCs wrapped in a `StandardResponse` envelope.
- Potential future: add secondary indexes (category, price range) and cache layer.

---

## 🧬 User Service

- TypeScript + Node.js.
- Provides user/customer data, credentials, validation logic.
- Consumed by Auth (login) and Product (owner/vendor checks) — integration via gRPC/Connect generated code.

---

## 🚪 API Gateway (Kong)

- Declarative config (`kong.yml`).
- Custom Lua plugins:
  - `auth-token-validator` / `auth-cookie-clearer`
  - `grpc-cookie-transformer` / `auth-metadata-setter`
  - `user-context-injector`
- `grpc-gateway` plugin converts JSON HTTP requests into gRPC upstream calls using service `.proto` files.
- Observability stack: Prometheus + Grafana + Loki + Promtail.

---

## 📡 Communication Protocols

- Internal service-to-service: gRPC.
- External clients: HTTP → Kong → gRPC (JSON → Protobuf via plugin).
- Planned async events (Order validation): Kafka + Outbox for reliability.

---

## 🛠 Local Development Quick Start

Prerequisites: Go toolchain, Node.js, Docker.

```bash
# 1. Start infra (MongoDB, Redis, DynamoDB Local, Kong stack)
# Example simplified (adjust as needed):
docker run -d --name mongo -p 27017:27017 mongo:6

docker run -d --name redis -p 6379:6379 redis:7

docker run -d --name dynamodb -p 9000:8000 amazon/dynamodb-local

# Start Kong + observability (from kongApiGateway/)
cd kongApiGateway
docker-compose up -d --build

# 2. Generate protos (each Go service)
cd product_service/proto && buf generate && cd ../..
cd cart_service/proto && buf generate && cd ../..
cd auth_service/proto && buf generate && cd ../..
# user_service uses TS codegen already present

# 3. Run services
cd auth_service && go run ./cmd/api &
cd product_service && go run ./cmd/api &
cd cart_service && go run ./cmd/api &
# user_service (Node)
cd user_service && npm install && npm run start:dev &

# 4. Test via Kong (JSON -> product create example)
curl -i -X POST http://localhost:8000/products \
  -H 'Content-Type: application/json' \
  -H 'Cookie: access-token=YOURTOKEN' \
  -d '{"name":"T-Shirt","description":"100% cotton","category":"apparel","price":19.99}'
```

> Adjust ports and env variables to match each service `.env`.

---

## 🧪 Testing Strategy (Current + Suggested)

| Layer              | Existing              | Suggested                                      |
| ------------------ | --------------------- | ---------------------------------------------- |
| Unit (Go services) | Not yet               | Table-driven tests for handlers/services/repos |
| Integration        | Manual with grpcurl   | Docker Compose test profile (all deps)         |
| Contract (Proto)   | buf                   | Add buf lint + breaking change checks CI       |
| Load / Perf        | None                  | k6 / Vegeta against Kong for key flows         |
| Security           | Token validation only | Add static scan (gosec), dependency checks     |

---

## 📁 Directory Structure (Root)

```
./
├── auth_service/
├── cart_service/
├── product_service/
├── user_service/
├── order_service/        # (in progress)
├── kongApiGateway/
└── README.md              # This file
```

Each service directory is self-contained with its own `proto/`, `internal/` (or `src/` for TS), and README.

---

## 🔄 Planned Additions

| Upcoming        | Description                              | Notes                              |
| --------------- | ---------------------------------------- | ---------------------------------- |
| Order README    | Document async validation pipeline       | Will include Kafka + Outbox        |
| Payment Service | Payment authorization/capture, refunds   | Integrates with Order events       |
| Event Bus       | Kafka cluster + schema evolution process | Topic naming + DLQ strategy        |
| Tracing         | OpenTelemetry instrumentation            | Correlate requests across services |
| CI/CD           | Build, lint, test, deploy pipelines      | Buf breaking check, security scan  |

---

## 🧮 Cross-Service Dependencies

- Auth → User Service (credentials & account creation)
- Cart → Product Service (product details, pricing)
- Product → User Service (seller/user validation)
- Gateway → All (routing, auth, translation)
- Future: Order → Product (validation), Cart (cart snapshot), User (buyer info), Payment → Order

Graph (simplified):

```
[Client]
   | HTTP JSON
 [Kong Gateway]
   | gRPC metadata injected
   +--> Auth Service --> User Service
   +--> Product Service --> User Service
   +--> Cart Service --> Product Service
   +--> (Future) Order Service --> Product/User/Cart
```

---

## 🛡 Security (Current State)

- JWT verification in Auth Service + Kong plugins.
- Redis + Mongo isolation via container ports (dev).
- No role-based authorization yet inside product/cart (planned).

> Add role checks at handler/service level (e.g., product modifications require seller role) and propagate `x-user-role` metadata.

---

## 📊 Observability

- Metrics + Logs for gateway (Prometheus + Grafana + Loki).
- Services: basic logging; add Prometheus instrumentation (request latency, error counts) later.
- Future: distributed tracing with OpenTelemetry.

---

## ⚙️ Proto & Code Generation

Central principles:

- Each service keeps only its needed proto files + buf configs.
- Avoid duplicating common messages (consider shared proto module or import once; current duplication acceptable short-term).
- Run `buf generate` after proto changes; commit generated artifacts (for now) to simplify onboarding.

---

## 🧹 Conventions / Style

- Go services: layered (`handler -> service -> repo -> infra`).
- Error handling: map internal errors to gRPC statuses via interceptor.
- Response wrapping: consistent `StandardResponse` or `CartStandardResponse` envelope.
- Environment: `.env` files without spaces around `=`.
- IDs: string/UUID across services for compatibility.

---

## 🤝 Contribution Guidelines (Initial)

1. Update relevant service README when adding new RPC or behavior.
2. Run `buf lint` (add CI) before committing proto changes.
3. Keep changes focused—avoid cross-service edits unless strictly required.
4. Add unit tests for new service logic.
5. For new messages/events: document topic purpose and schema in future `/docs/events.md`.

---

## ❓ FAQ (Early)

| Question                          | Answer                                                                          |
| --------------------------------- | ------------------------------------------------------------------------------- |
| Why Redis blob for cart?          | Simple atomic pattern; low complexity for typical cart sizes                    |
| Why DynamoDB for product?         | Fits flexible product attributes & scales; local dev uses DynamoDB Local        |
| Why Kong?                         | Unified gateway for auth plugins + JSON→gRPC translation + observability        |
| Why separate Auth and User?       | Decouples authentication/session logic from core user domain CRUD               |
| Why not sync validation in Order? | Async design scales better; avoids coupling order throughput to product latency |

---

## 🔮 Roadmap (Next Priorities)

1. Finalize Order validation events (Kafka topics + outbox implementation).
2. Introduce Payment Service with idempotent charge API.
3. Shared proto module (reduce message duplication).
4. Metrics + tracing across all services.
5. CI pipeline (lint, test, security scan, build images).

---

## 📄 License & Legal

(Choose and add a LICENSE file; currently unspecified.)

---

## 🏁 Getting Help

- Read per-service README for deep dive.
- Check Kong gateway logs for integration issues (`docker logs -f kong`).
- Use `grpcurl` for quick RPC manual tests.

---

Feel free to request a diagram-enhanced or minimal quick-start version. This root README will evolve as Order/Payment services are completed.
