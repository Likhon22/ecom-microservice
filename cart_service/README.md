# Cart Service

A gRPC microservice that manages user shopping carts. It stores cart state as a single JSON blob in Redis with a TTL, exposes protobuf RPCs for cart operations, calls the Product Service for product details/validation, and returns consistent wrapped responses.

---

## 🎯 Overview

The Cart Service lets clients:

- Create / add items to a cart
- Update item quantities
- Remove items
- Fetch the whole cart (with subtotal recalculated)
- Clear cart

It uses Redis for fast access and expiration, integrating product data (price, name) at insertion time so reads are cheap.

---

## 🏗 Architecture

Layered layout (similar to other services):

```
cart_service/
├── cmd/api/main.go                # Entrypoint (loads config & starts gRPC server)
├── internal/
│   ├── api/handlers/handlers.go   # gRPC handlers (metadata extraction, response wrapping)
│   ├── bootstrap/server.go        # Wiring (Redis client, product client, repo, service, handler)
│   ├── clients/product/client.go  # gRPC client to Product Service
│   ├── config/config.go           # Environment / config loader
│   ├── domain/cart.go             # Domain models (Cart, CartItem)
│   ├── infra/redis.go             # Redis connection setup
│   ├── interceptors/errorInterceptor.go # Maps domain errors to gRPC statuses
│   ├── repo/cart/repo.go          # Persistence (Redis blob)
│   ├── services/cart/service.go   # Business logic (add/update/remove/recalc)
│   └── utils/*.go                 # Pure helpers (subtotal, key creation, item index, mapping errors)
└── proto/                         # cart.proto + product.proto + buf configs
```

Key separation:

- Handler: gRPC → service call → wraps response (sets success/message/status_code).
- Service: orchestrates product lookups & repo persistence; recalculates subtotal.
- Repo: serializes whole cart to JSON and stores under a Redis key (email-based namespace) with TTL.
- Utils: pure functions (recalculate subtotal, find item, create cart item, error mapping).

---

## 🗄 Data Model

`Cart` is stored as one JSON blob keyed by `cart:{email}` (see `utils/key.go`).

```
Cart {
  email: string,
  items: [
    { product_id, name, price, quantity }, ...
  ],
  subtotal: number
}
```

Pros of single-blob pattern:

- Atomic updates (read-modify-write)
- Simplicity
  Cons:
- Large carts mean bigger payload rewrite on every change
- Concurrent updates require careful read-modify-write (possible lost update risk)

TTL (7 days) configured on write to auto-expire stale carts.

---

## 🔌 External Dependency

- **Product Service** (gRPC): used during Add or Update operations to verify product existence & price.

---

## 🔐 Metadata & Auth

Handlers expect `x-user-email` (or similar) in gRPC metadata (injected upstream by gateway plugins). That email becomes the cart key. If metadata missing, handler returns an error.

---

## 📦 Protobuf / RPCs (cart.proto)

Common RPC pattern (examples):

- `AddToCart(AddToCartRequest) returns (CartStandardResponse)`
- `GetCart(GetCartRequest) returns (CartStandardResponse)`
- `UpdateCartItem(UpdateItemRequest) returns (CartStandardResponse)`
- `RemoveCartItem(RemoveItemRequest) returns (CartStandardResponse)`
- `ClearCart(ClearCartRequest) returns (CartStandardResponse)`

`CartStandardResponse` contains:

- `success` (bool)
- `message` (string)
- `status_code` (int32)
- `cart_data` (Cart message) or relevant payload inside a `oneof` field

If you modify proto:

```bash
cd proto
buf generate
```

---

## ⚙️ Configuration (.env)

Example `.env` (no spaces around `=`):

```
VERSION=1
SERVICE_NAME=cart_service
ADDR=":5004"          # gRPC listen address
REDIS_ADDR=localhost:6379
REDIS_DB=0
PRODUCT_SERVICE_ADDR=localhost:5003
```

Load order: `config/config.go` reads env; bootstrap wires clients.

---

## 🚀 Run Locally

Prerequisites: Go, Redis, Product Service running.

```bash
# start Redis (if not already)
docker run -d --name redis -p 6379:6379 redis:7

# generate protos if changed
cd proto && buf generate && cd ..

# build
go build -o bin/main ./cmd/api

# run
./bin/main
```

Check healthy (example gRPC CLI using grpcurl):

```bash
grpcurl -plaintext localhost:5004 list
grpcurl -plaintext -d '{"email":"user@example.com"}' localhost:5004 cart.CartService/GetCart
```

(Adjust service/method names to match actual proto.)

---

## 🔁 Core Logic Details

### AddToCart Flow:

1. Handler extracts email from metadata.
2. Service fetches current cart via repo.GetCart.
3. Calls product client to get product details (name, price).
4. If item exists → increment quantity; else create new item.
5. Recalculate subtotal (sum of price \* quantity).
6. Persist entire cart (repo.AddToCart) writing JSON blob to Redis with TTL.
7. Return wrapped response.

### Update Quantity:

- Find item index.
- Set or increment based on request semantics (ensure you define which in proto).
- Recalculate subtotal.
- Persist.

### Remove Item / Clear Cart:

- Remove item slice entry (or delete key for clear).
- Recalculate subtotal (if not clearing).
- Persist updated cart or delete Redis key.

---

## 🧮 Subtotal Calculation

In `utils/recalculateSubTotal.go`: loops items and sums `price * quantity`. Always recomputed after any modification to keep source of truth simple.

---

## 🧠 Concurrency Considerations

Potential lost update scenario:

- Two simultaneous Add operations read old cart, both modify, last write wins.
  Mitigations (future improvements):
- Use Redis WATCH (optimistic lock) and retry on conflict.
- Switch to per-item keys (trade-off complexity).
- Introduce a small queue if high contention expected.
  For typical e-commerce browsing patterns, low probability of conflicting writes for a single user cart; acceptable simplification.

---

## 🧪 Testing Suggestions

Unit Tests:

- Service layer: Add, Update, Remove logic (mock product client, mock repo).
- Repo: JSON marshal/unmarshal roundtrip + TTL.
  Integration:
- Spin up real Redis and test full flow with grpcurl or a test harness.

Sample Go test snippet (pseudo):

```go
func TestAddToCart_NewItem(t *testing.T) {
  // mock product client returning name+price
  // repo with in-memory redis (use miniredis)
  // call service.AddToCart -> assert subtotal, items length = 1
}
```

---

## 🔧 Error Handling

`interceptors/errorInterceptor.go` maps internal errors to gRPC status codes (e.g., NotFound, InvalidArgument). Ensure service returns typed errors (wrap or sentinel) so interceptor can map correctly.

Common error sources:

- Missing email metadata → InvalidArgument
- Product not found → NotFound
- Redis connectivity issues → Internal

---

## 🛠 Extending the Service

Add field (e.g., currency) → steps:

1. Update `cart.proto` (add field to CartItem & Cart).
2. `buf generate`.
3. Adjust domain/cart.go and utils for subtotal (e.g., multi-currency logic).
4. Update service to populate currency.

Add discount codes:

- Store discount metadata in cart blob.
- Recalculate subtotal with discount rules.

Add event publishing (e.g., Kafka on checkout): integrate outbox pattern similar to Order Service.

---

## 🩹 Troubleshooting

| Symptom                                  | Cause                                                    | Fix                                                                         |
| ---------------------------------------- | -------------------------------------------------------- | --------------------------------------------------------------------------- |
| Empty cart returned when expected items  | Wrong email metadata or TTL expired                      | Verify metadata header & TTL setting                                        |
| Redis `json: Unmarshal(nil *Cart)` error | Attempt to unmarshal into nil pointer without allocation | Ensure repo creates struct before Unmarshal (already fixed pattern)         |
| 415 from gateway                         | Missing grpc-gateway plugin for cart route               | Add plugin with cart.proto path                                             |
| Subtotal incorrect                       | Price change in Product Service not synced               | Decide: refresh on each GetCart or store a timestamp & refresh stale prices |
| Add increments incorrectly               | Logic mismatch (increment vs set)                        | Clarify semantics in proto & update service code                            |

---

## 📏 Performance Notes

- Single key JSON blob → O(cart_size) update cost (serialize entire cart); fine for small carts (<100 items typical).
- Redis expiration reduces storage of abandoned carts.
- Consider enabling compression if cart blobs become large (application side before Set).

---

## ✅ Quick Checklist

- [ ] Redis reachable (REDIS_ADDR)
- [ ] Product Service reachable (PRODUCT_SERVICE_ADDR)
- [ ] Protos generated (`buf generate` ran successfully)
- [ ] gRPC server listening on configured ADDR
- [ ] Gateway route configured with grpc-gateway plugin + cart.proto

---

## 🔮 Future Enhancements

- Optimistic locking (Redis WATCH) for high-concurrency carts
- Partial item updates to avoid full blob rewrite
- Metrics (Prometheus): cart size, operation latency, Redis errors
- Preloading product data batch for multi-item add performance
- Incorporate inventory reserve step before checkout

---

## ℹ️ References

- Redis: https://redis.io/
- gRPC Go: https://grpc.io/docs/languages/go/
- Buf (Protobuf): https://buf.build/docs/introduction
- Segmentio kafka-go (if extended later): https://github.com/segmentio/kafka-go

---

Feel free to ask for a slimmer version or add diagrams (e.g., Add flow sequence). This README aims to onboard new contributors quickly.
