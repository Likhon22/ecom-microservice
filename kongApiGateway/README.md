# Kong API Gateway

Central HTTP entrypoint for the e-commerce microservices. Kong terminates HTTP(S) requests, applies custom Lua plugins (auth & metadata), translates selected JSON REST calls to backend gRPC services (auth, product, cart, user) using the `grpc-gateway` plugin, and exposes metrics/logs to Prometheus, Loki/Grafana.

---

## 🎯 Responsibilities

- Route external HTTP requests to internal gRPC services (gRPC-over-HTTP/2 upstreams)
- Perform auth header / cookie validation & transformation (custom plugins)
- Inject user context metadata into gRPC requests
- Translate JSON payloads to Protobuf via `grpc-gateway` (when configured)
- Provide centralized observability (Prometheus metrics, Loki logs, Grafana dashboards)

---

## 📁 Directory Layout

```
kongApiGateway/
├── docker-compose.yml        # Local stack: Kong + Prometheus + Grafana + Loki + Promtail
├── kong.yml                  # Declarative Kong configuration (routes/services/plugins)
├── kong.tpl.yml              # Template variant (if envsubst or CI substitution needed)
├── protos/                   # Protobuf files required for grpc-gateway translation
│   ├── auth.proto
│   ├── cart.proto
│   ├── product.proto
│   ├── user.proto
│   └── ... (google/ validate/ dependencies)
├── plugins/                  # Custom Lua plugins
│   ├── auth-cookie-clearer/
│   ├── auth-metadata-setter/
│   ├── auth-token-validator/
│   ├── grpc-cookie-transformer/
│   └── user-context-injector/
├── grafana-dashboards/       # JSON dashboard definitions
├── grafana-provisioning/     # Datasource auto-provision (Prometheus, Loki)
├── loki-config.yml           # Loki configuration
├── promtail-config.yml       # Promtail log shipper configuration
├── prometheus.yml            # Prometheus scrape config (includes Kong metrics)
├── .env                      # Environment variables for docker-compose / template
└── start-kong.sh             # Helper script (optional) to start Kong in custom ways
```

---

## 🔌 Custom Plugins (Lua)

Order of execution matters. Typical chain for an authenticated request:

1. `auth-token-validator` — validates JWT / refresh tokens or session cookies
2. `auth-cookie-clearer` — clears/rewrites cookies when invalid or expired
3. `grpc-cookie-transformer` — converts HTTP cookies/headers into gRPC metadata
4. `auth-metadata-setter` — sets auth/user metadata headers (e.g., x-user-email)
5. `user-context-injector` — injects additional user context into the outbound gRPC call

### Plugin Summaries

- **auth-token-validator**: Parses JWT/cookies, checks signature/expiry, sets failure status if invalid.
- **auth-cookie-clearer**: Removes stale cookies so clients do not keep re-sending invalid tokens.
- **grpc-cookie-transformer**: Normalizes incoming HTTP headers/cookies into gRPC metadata pairs (e.g., `x-user-email`).
- **auth-metadata-setter**: Adds standardized auth metadata fields required by downstream services.
- **user-context-injector**: Adds auxiliary context (roles, device id) to facilitate downstream authorization decisions.

> If any plugin detects invalid credentials it can short-circuit the request with appropriate HTTP status (401/403).

---

## 🧬 Protobuf & gRPC Gateway

The `grpc-gateway` plugin enables Kong to accept JSON over HTTP, convert it to Protobuf, and forward as gRPC to upstream services.

Requirements for a route using grpc-gateway:

- Corresponding `.proto` file present under `protos/` inside the Kong container
- `plugins:` list contains the `grpc-gateway` plugin with `config.proto: /path/to/file.proto`
- Upstream `service.url` uses `grpc://` scheme (HTTP/2 connection to backend port)

Example snippet (simplified) from `kong.yml`:

```yaml
services:
  - name: product-service
    url: grpc://host.docker.internal:5003
    routes:
      - name: product-create
        paths: [/products]
        methods: [POST]
    plugins:
      - name: grpc-gateway
        config:
          proto: /usr/local/kong/protos/product.proto
      - name: auth-token-validator
      - name: grpc-cookie-transformer
```

If `grpc-gateway` plugin is missing for a JSON request, you may see `415 Unsupported Media Type` because Kong cannot translate the body to Protobuf.

---

## ⚙️ Environment Variables (.env)

Typical entries:

```
KONG_DATABASE=off              # declarative config mode
KONG_DECLARATIVE_CONFIG=/usr/local/kong/kong.yml
KONG_PROXY_LISTEN=0.0.0.0:8000, 0.0.0.0:8443 ssl
KONG_ADMIN_LISTEN=0.0.0.0:8001
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
LOKI_PORT=3100
```

Adjust host ports as needed for local dev. `docker-compose.yml` reads these values.

---

## 🚀 Local Development

Prerequisites: Docker / Docker Compose installed.

Start the stack:

```bash
docker-compose up -d --build
```

Check running containers:

```bash
docker ps --format 'table {{.Names}}\t{{.Status}}'
```

Tail Kong logs (via Promtail/Loki or directly):

```bash
docker logs -f kong
```

Access:

- Proxy: `http://localhost:8000`
- Admin (if enabled): `http://localhost:8001`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (default admin/admin unless overridden)

Stop:

```bash
docker-compose down -v
```

---

## 🛣️ Adding a New gRPC Service Route

1. Place the service's `.proto` file in `protos/` (include dependencies like google/ or validate/ if needed).
2. Update `kong.yml`:
   - Add a `service:` entry with `url: grpc://host.docker.internal:<port>`
   - Add a `route:` (paths, methods) for HTTP exposure
   - Attach `grpc-gateway` plugin with the correct `proto` path
   - Chain auth/metadata plugins as required
3. Rebuild / restart Kong:
   ```bash
   docker-compose restart kong
   ```
4. Test request:
   ```bash
   curl -i -X POST http://localhost:8000/new-resource \
     -H 'Content-Type: application/json' \
     -d '{"field":"value"}'
   ```
5. Observe transformed gRPC request logs at the backend service.

---

## 🪵 Observability Stack

- **Prometheus**: Scrapes Kong metrics endpoint (enable the prometheus plugin if required).
- **Grafana**: Dashboards in `grafana-dashboards/` auto-provisioned (see `grafana-provisioning/`).
- **Loki + Promtail**: Promtail tails Kong container logs and ships to Loki; Grafana queries Loki for log panels.

Metrics to watch:

- Request rate, latency (p50/p95/p99)
- Upstream failures (5xx) / downstream failures (4xx)
- Plugin execution times (if exposed)

Logs to inspect for troubleshooting plugin or routing issues.

---

## 🔐 Authentication & Metadata Flow

Typical request steps:

1. Client sends cookies/JWT with HTTP request.
2. `auth-token-validator` verifies token; rejects or proceeds.
3. `grpc-cookie-transformer` normalizes tokens into gRPC metadata (e.g., `x-user-email`).
4. `auth-metadata-setter` adds consistent metadata fields used by microservices (user role, device id).
5. `user-context-injector` enriches with auxiliary context.
6. `grpc-gateway` plugin marshals JSON → Protobuf → forwards gRPC request upstream.
7. Upstream gRPC response → (optionally) converted back to JSON → client.

---

## 🧪 Testing Routes Quickly

JSON → product create (example):

```bash
curl -i -X POST http://localhost:8000/products \
  -H 'Content-Type: application/json' \
  -H 'Cookie: access-token=YOURTOKEN' \
  -d '{"name":"T-Shirt","description":"100% cotton","category":"apparel","price":19.99}'
```

If you receive `415 Unsupported Media Type` ensure `grpc-gateway` plugin is present and proto path correct.

Auth debugging tip:

```bash
curl -i http://localhost:8000/products -H 'Authorization: Bearer INVALID' -d '{}'
```

Expect `401` or `403` based on validator plugin logic.

---

## 🛠️ Troubleshooting

| Symptom                          | Likely Cause                                                   | Fix                                                        |
| -------------------------------- | -------------------------------------------------------------- | ---------------------------------------------------------- |
| 415 Unsupported Media Type       | Missing `grpc-gateway` plugin or wrong proto path              | Add plugin / correct path; restart Kong                    |
| 502 Bad Gateway                  | Upstream gRPC service down or port mismatch                    | Verify service running & port in `kong.yml`                |
| Auth always failing              | Token plugin misconfigured or secret mismatch                  | Check plugin config & token signing key                    |
| No metrics in Prometheus         | Prometheus plugin disabled or wrong scrape target              | Enable plugin & verify `prometheus.yml`                    |
| Missing user metadata downstream | `grpc-cookie-transformer` or `auth-metadata-setter` path error | Confirm plugin order and headers                           |
| Grafana dashboards empty         | Datasource provisioning failed                                 | Inspect `grafana-provisioning/datasources/datasources.yml` |

---

## 🧩 Performance & Scaling Notes

- Use upstream keep-alive (default) to reduce gRPC connection churn.
- Keep proto files small & versioned; regenerate clients when fields change.
- Add rate-limiting plugins for public endpoints.
- Horizontal scale: run multiple Kong instances behind a load balancer; share Prometheus/Loki.
- Use caching plugin (optional) for read-heavy endpoints.

---

## ➕ Adding a Custom Lua Plugin (Quick Guide)

1. Create directory under `plugins/<your-plugin>/` with `schema.lua` and `handler.lua`.
2. Reference plugin by name in `kong.yml` under `plugins:` for route or global.
3. Rebuild image (if Docker) so plugin code is available.
4. Validate with `kong check kong.yml` before starting.

Minimal `handler.lua` skeleton:

```lua
local MyPlugin = {
  PRIORITY = 800,
  VERSION = "1.0.0",
}
function MyPlugin:access(conf)
  -- modify request headers / abort if needed
end
return MyPlugin
```

---

## 🧪 Validating Configuration

Pre-flight check:

```bash
kong config parse kong.yml
```

(Or `kong check kong.yml` if available.)

---

## ♻️ Updating Protos

After changing a service proto:

1. Copy updated `.proto` into `protos/`
2. Restart Kong container
3. Test JSON → gRPC translation; if fields missing ensure plugin points to correct file

---

## ✅ Minimal Checklist Before Running

- [ ] Protos for all gRPC services present under `protos/`
- [ ] `kong.yml` references correct upstream ports (auth/cart/product/user)
- [ ] `grpc-gateway` plugin configured for JSON endpoints
- [ ] Auth plugins listed in proper order
- [ ] Observability stack containers up (Prometheus, Grafana, Loki, Promtail)
- [ ] `.env` variables match local port selections

---

## 📦 Future Enhancements

- Global rate limiting (per IP / user)
- Distributed tracing headers injection (OpenTelemetry plugin)
- Canary routing / A/B testing via traffic-splitting plugin
- JSON schema validation pre-gRPC translation
- mTLS between Kong and backend services

---

## ℹ️ References

- Kong Docs: https://docs.konghq.com/
- grpc-gateway plugin: https://docs.konghq.com/hub/kong-inc/grpc-gateway/
- Prometheus plugin: https://docs.konghq.com/hub/kong-inc/prometheus/
- Lua plugin development: https://docs.konghq.com/gateway/latest/plugin-development/

---

Feel free to request a lighter or expanded version (e.g., with diagrams). This README aims to be a comprehensive starting point for new contributors.
