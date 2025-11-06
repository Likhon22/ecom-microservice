# Auth Service

A high-performance authentication and authorization microservice built with Go and gRPC, providing JWT-based authentication with dual-token strategy (access + refresh tokens) and multi-device session management.

## 🎯 Overview

The Auth Service is a core component of the e-commerce microservices architecture, responsible for:

- User account creation (delegates to User Service)
- Authentication (login/logout)
- JWT token generation and validation
- Refresh token management with Redis caching
- Multi-device session tracking
- Secure cookie-based token delivery

## 🏗️ Architecture

### Technology Stack

- **Language**: Go 1.25.1
- **Framework**: gRPC with Protocol Buffers
- **Databases**:
  - MongoDB (persistent refresh token storage)
  - Redis (fast token caching & lookup)
- **Authentication**: JWT (HMAC-SHA256)
- **Password Hashing**: bcrypt
- **Validation**: protoc-gen-validate
- **Logging**: zerolog
- **Hot Reload**: Air

### Project Structure

```
auth_service/
├── cmd/api/                    # Application entry point
│   └── main.go                # Initializes config, app, and starts server
├── internal/
│   ├── api/handler/           # gRPC request handlers
│   │   ├── auth/              # Authentication endpoints
│   │   └── user/              # User account endpoints
│   ├── bootstrap/             # Application initialization
│   │   └── server.go          # gRPC server setup & DI
│   ├── clients/               # External service clients
│   │   └── usersvc/           # User service gRPC client
│   ├── config/                # Configuration management
│   │   ├── config.go          # Main config loader
│   │   ├── authConfig.go      # JWT secrets & expiration
│   │   └── dbConfig.go        # Database connections
│   ├── domain/                # Domain models
│   │   └── token.go           # Token entities
│   ├── infra/db/              # Infrastructure layer
│   │   ├── db.go              # MongoDB connection
│   │   └── redis.go           # Redis connection
│   ├── interceptors/          # gRPC interceptors
│   │   └── errorInterceptor.go # Error handling & HTTP code mapping
│   ├── repo/auth/             # Data access layer
│   │   └── repo.go            # Token CRUD with Redis + MongoDB
│   ├── services/auth/         # Business logic layer
│   │   └── service.go         # Auth operations
│   ├── types/                 # Shared types
│   │   └── types.go
│   └── utils/                 # Helper functions
│       ├── jwt.go             # Token signing & parsing
│       ├── mapError.go        # Error mapping
│       ├── metadata.go        # Metadata handling
│       ├── tokenFromCtx.go    # Extract tokens from context
│       └── wrapSuccess.go     # Success response wrapper
├── proto/                     # Protocol Buffer definitions
│   ├── auth.proto             # Auth service contract
│   ├── user.proto             # User service contract (for client)
│   ├── buf.gen.yaml           # Buf code generation config
│   └── gen/                   # Generated proto code
├── .env                       # Environment variables
├── .air.toml                  # Air hot-reload configuration
└── makefile                   # Build commands

```

## 🔐 Security Features

### Dual Token Strategy

1. **Access Token**

   - Short-lived (default: 5 minutes)
   - Used for API authentication
   - Stored in HttpOnly, Secure cookies
   - Contains: email, role, device_id

2. **Refresh Token**
   - Long-lived (default: 24 hours)
   - Used to obtain new access tokens
   - Stored in both Redis (cache) and MongoDB (persistence)
   - Revocable per device

### Token Storage Architecture

```
┌─────────────────────────────────────────────────┐
│  Login Request                                  │
│  ↓                                              │
│  Generate Access + Refresh Tokens               │
│  ↓                                              │
│  Store Refresh Token:                           │
│    • Redis (key: email:device_id) ←─ Fast      │
│    • MongoDB (upsert by email+device)           │
│  ↓                                              │
│  Return tokens via Set-Cookie headers           │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  Token Validation Flow                          │
│  ↓                                              │
│  Check Redis first (O(1) lookup)                │
│  ↓                                              │
│  If miss → Fetch from MongoDB → Cache in Redis  │
│  ↓                                              │
│  Verify: signature, expiry, revoked status      │
└─────────────────────────────────────────────────┘
```

### Multi-Device Support

- Each user can have multiple active sessions (different devices)
- Tokens tracked by: `email + device_id`
- Logout revokes token for specific device only
- Redis key pattern: `refreshToken{email}:{device_id}`

## 📡 gRPC API

### Service Definition

```protobuf
service AuthService {
  rpc Login(LoginRequest) returns (StandardResponse);
  rpc ValidateRefreshToken(ValidateRefreshTokenRequest) returns (StandardResponse);
  rpc Logout(LogoutRequest) returns (StandardResponse);
  rpc CreateUserAccount(CreateUserRequest) returns (StandardResponse);
}
```

### Endpoints

#### 1. Create User Account

**Route**: `POST /users`

Creates a new user account by delegating to User Service.

**Request**:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123",
  "phone": "+1234567890",
  "address": "123 Main St",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

**Response**:

```json
{
  "success": true,
  "message": "user created successfully",
  "status_code": 201,
  "user_data": {
    "name": "John Doe",
    "email": "john@example.com",
    "role": "customer",
    "status": "active"
  }
}
```

**Validations**:

- `name`: min_len = 1
- `email`: valid email format
- `password`: min_len = 6

---

#### 2. Login

**Route**: `POST /auth/login`

Authenticates user and returns JWT tokens via cookies.

**Request**:

```json
{
  "email": "john@example.com",
  "password": "securepass123",
  "device_id": "web-chrome-mac"
}
```

**Response**:

```json
{
  "success": true,
  "message": "login successful",
  "status_code": 200,
  "login_data": {
    "message": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Set-Cookie Headers**:

```
Set-Cookie: access-token={jwt}; HttpOnly; Secure; Max-Age=300
Set-Cookie: refresh-token={jwt}; HttpOnly; Secure; Max-Age=86400
```

**Flow**:

1. Fetch user credentials from User Service
2. Verify password with bcrypt
3. Generate access token (5min TTL)
4. Generate refresh token (24h TTL)
5. Store refresh token in Redis + MongoDB
6. Return tokens in cookies

---

#### 3. Validate Refresh Token

**Route**: `POST /auth/refresh/validate`

Issues new access token using valid refresh token.

**Request**:

```json
{
  "refresh_token": "from-metadata"
}
```

**Response**:

```json
{
  "success": true,
  "message": "new access token generated",
  "status_code": 200,
  "refreshTokenData": {
    "message": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Validation Checks**:

- Parse JWT and verify signature
- Check expiration time
- Verify token exists in database
- Check not revoked
- Match stored token with provided token

**Set-Cookie Header**:

```
Set-Cookie: access-token={new_jwt}; HttpOnly; Secure; Max-Age=300
```

---

#### 4. Logout

**Route**: `POST /auth/logout`

Revokes refresh token for current device.

**Request**:

```json
{
  "refresh_token": "from-metadata"
}
```

**Response**:

```json
{
  "success": true,
  "message": "logout successful",
  "status_code": 200,
  "logout_data": {
    "message": "logout successful"
  }
}
```

**Flow**:

1. Parse refresh token to extract email + device_id
2. Delete from Redis cache
3. Delete from MongoDB
4. Return success

---

## 🔌 External Dependencies

### User Service Client

The Auth Service communicates with User Service via gRPC for:

- `CreateCustomer` - Account creation
- `GetCustomerByEmail` - Profile lookup
- `GetCustomerCredentials` - Password verification
- `GetCustomers` - List users
- `DeleteCustomer` - Account deletion

**Client Setup**:

```go
conn, err := grpc.DialContext(ctx,
    userServiceAddr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithBlock(),
    grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
)
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the root directory:

```bash
# Service Configuration
VERSION=1.0.0
SERVICE_NAME=auth-service
ADDR=:5002
USER_SERVICE_ADDR=localhost:5001

# JWT Configuration
JWT_ACCESS_TOKEN_SECRET=your-super-secret-access-key-min-32-chars
JWT_REFRESH_TOKEN_SECRET=your-super-secret-refresh-key-min-32-chars
ACCESS_TOKEN_EXP_DURATION=5m
RESET_TOKEN_EXP_DURATION=24h

# Database Configuration
DB_URL=mongodb://localhost:27017
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_MAX_IDLE_TIME=15m

# Redis Configuration (hardcoded in bootstrap)
REDIS_ADDR=localhost:6379
REDIS_DB=0
```

### Required External Services

1. **MongoDB** (port 27017)

   - Database: `auth_service`
   - Collection: `refresh_tokens`
   - Schema:
     ```javascript
     {
       email: String,
       device_id: String,
       token: String,
       created_at: Date,
       expires_at: Date,
       revoked: Boolean
     }
     ```

2. **Redis** (port 6379)

   - Used for token caching
   - TTL matches refresh token expiration

3. **User Service** (gRPC on port 5001)
   - Must be running for auth operations
   - Provides user CRUD and credentials

## 🚀 Getting Started

### Prerequisites

- Go 1.25.1 or higher
- MongoDB 4.4+
- Redis 6.0+
- Protocol Buffers compiler (`protoc`)
- Buf CLI (for proto generation)
- Air (for hot reload)

### Installation

1. **Clone the repository**

   ```bash
   cd auth_service
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Install development tools**

   ```bash
   # Install Air for hot reload
   go install github.com/air-verse/air@latest

   # Install Buf for proto generation
   brew install bufbuild/buf/buf  # macOS
   # or visit https://buf.build/docs/installation
   ```

4. **Generate proto files**

   ```bash
   cd proto
   buf generate
   cd ..
   ```

5. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

6. **Start dependencies**

   ```bash
   # Start MongoDB
   docker run -d -p 27017:27017 --name mongo mongo:latest

   # Start Redis
   docker run -d -p 6379:6379 --name redis redis:latest
   ```

7. **Run the service**

   ```bash
   # Development (with hot reload)
   make run

   # Or build and run
   go build -o bin/main cmd/api/main.go
   ./bin/main
   ```

### Testing

**Test with grpcurl**:

```bash
# Login
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "password123",
  "device_id": "web-1"
}' localhost:5002 auth_service.AuthService/Login

# Validate Refresh Token (requires metadata)
grpcurl -plaintext \
  -H "refresh-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{}' \
  localhost:5002 auth_service.AuthService/ValidateRefreshToken
```

**Test via Kong Gateway** (if configured):

```bash
# Login
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "device_id": "web-1"
  }'

# Refresh token
curl -X POST http://localhost:8000/auth/refresh/validate \
  -H "Cookie: refresh-token=..." \
  -H "Content-Type: application/json"
```

## 🏛️ Design Patterns

### Dependency Injection

- Constructor-based injection
- Interface-based abstractions
- Single Responsibility Principle

### Repository Pattern

- `AuthRepo` interface for data access
- Dual storage strategy (Redis + MongoDB)
- Read-through cache pattern

### Service Layer

- Business logic separated from handlers
- Single service interface with multiple methods
- Error handling via custom mapping

### Middleware (Interceptors)

- Error interceptor converts gRPC errors to HTTP codes
- Graceful error response wrapping
- Consistent API responses

## 📊 Data Flow

```
┌──────────────┐
│  gRPC Client │
└──────┬───────┘
       │
       ↓
┌──────────────────────────────────────┐
│  gRPC Server (interceptors)          │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│  Handlers (auth/user)                │
│  • Request validation                │
│  • Response wrapping                 │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│  Service Layer                       │
│  • Business logic                    │
│  • User service client calls         │
│  • Token generation                  │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│  Repository Layer                    │
│  • Token storage/retrieval           │
│  • Redis caching                     │
│  • MongoDB persistence               │
└──────────────────────────────────────┘
```

## 🔄 Token Refresh Flow

```
Client                 Auth Service              Redis         MongoDB
  │                         │                      │              │
  │─ POST /auth/login ─────→│                      │              │
  │                         │─ Verify password ────→ User Service │
  │                         │                      │              │
  │                         │─ Generate tokens     │              │
  │                         │─ SET refresh token ─→│              │
  │                         │─ UPSERT refresh ─────────────────→│
  │←─ 200 + cookies ────────│                      │              │
  │                         │                      │              │
  │ (5 min later...)        │                      │              │
  │                         │                      │              │
  │─ POST /auth/refresh ───→│                      │              │
  │                         │─ GET token ─────────→│              │
  │                         │←─ HIT ───────────────│              │
  │                         │─ Verify signature    │              │
  │                         │─ Check expiry        │              │
  │                         │─ Generate new access │              │
  │←─ 200 + new cookie ─────│                      │              │
```

## 🛡️ Error Handling

### Error Mapping

gRPC codes are mapped to HTTP status codes:

| gRPC Code        | HTTP Status | Use Case            |
| ---------------- | ----------- | ------------------- |
| InvalidArgument  | 400         | Bad request data    |
| Unauthenticated  | 401         | Invalid credentials |
| PermissionDenied | 403         | Forbidden action    |
| NotFound         | 404         | User not found      |
| AlreadyExists    | 409         | Duplicate email     |
| Internal         | 500         | Server errors       |

### Standard Response Format

```json
{
  "success": false,
  "message": "password does not match",
  "status_code": 401,
  "result": null
}
```

## 🧩 Integration with Kong Gateway

The service is designed to work behind Kong API Gateway with custom plugins:

1. **grpc-gateway**: Converts HTTP → gRPC
2. **grpc-cookie-transformer**: Extracts cookies → metadata
3. **auth-metadata-setter**: Injects refresh token from cookies
4. **auth-cookie-clearer**: Clears cookies on logout

See `kongApiGateway/` directory for plugin configurations.

## 📈 Performance Considerations

### Redis Caching Strategy

- **Read Pattern**: Redis first, MongoDB fallback
- **Write Pattern**: Write to both (Redis for speed, MongoDB for durability)
- **TTL Management**: Auto-expiry in Redis matches token expiration

### Connection Pooling

- MongoDB: Configurable max open/idle connections
- gRPC: Persistent connection to User Service
- Redis: Single client instance (thread-safe)

### Graceful Shutdown

- Context-based cancellation
- gRPC server graceful stop
- Database connections closed on shutdown

## 🔮 Future Enhancements

- [ ] Rate limiting on login attempts
- [ ] Password reset flow
- [ ] Email verification
- [ ] OAuth2/OIDC integration
- [ ] Token blacklisting for compromised tokens
- [ ] Audit logging for security events
- [ ] Metrics (Prometheus)
- [ ] Distributed tracing (OpenTelemetry)

## 📝 License

This project is part of a larger e-commerce microservices architecture.

## 👥 Contributing

This is a learning/portfolio project. For production use, consider additional security hardening:

- Secret rotation mechanism
- Token introspection endpoint
- MFA support
- IP-based restrictions
- Anomaly detection

---

**Built with ❤️ using Go, gRPC, MongoDB, and Redis**
