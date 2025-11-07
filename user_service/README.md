# User Service

A TypeScript-based user management microservice built with Node.js, Express, gRPC, and MongoDB, providing dual interface (REST + gRPC) for customer account management in an e-commerce platform.

## 🎯 Overview

The User Service is a foundational microservice responsible for:

- Customer account creation and management
- User credential storage and retrieval
- Profile management (name, email, phone, address, avatar)
- Soft deletion support
- Dual protocol support (REST HTTP and gRPC)
- Integration with Auth Service for authentication workflows

## 🏗️ Architecture

### Technology Stack

- **Runtime**: Node.js with TypeScript
- **Framework**: Express.js (REST API)
- **RPC**: gRPC with @grpc/grpc-js
- **Database**: MongoDB with Mongoose ODM
- **Validation**: Zod + protoc-gen-validate
- **Password Hashing**: bcrypt
- **Development**: tsx (hot reload), ESLint, Prettier
- **Protocol Buffers**: @bufbuild packages

### Project Structure

```
user_service/
├── src/
│   ├── cmd/api/                   # Application entry point
│   │   └── main.ts               # Bootstrap logic, DI setup
│   ├── internal/
│   │   ├── server.ts             # HTTP & gRPC server initialization
│   │   ├── api/
│   │   │   ├── grpc/             # gRPC layer
│   │   │   │   ├── handler/      # gRPC request handlers
│   │   │   │   ├── types/        # gRPC type definitions
│   │   │   │   └── utils/        # Error mapping & response factory
│   │   │   ├── handlers/         # REST API handlers
│   │   │   │   └── userCustomer/ # Customer CRUD endpoints
│   │   │   ├── middleware/       # Express middleware
│   │   │   └── routes/           # REST route definitions
│   │   ├── config/               # Configuration management
│   │   │   └── index.ts          # Environment variable loader
│   │   ├── domain/               # Domain models & DTOs
│   │   │   ├── user.domain.ts    # User entity type
│   │   │   ├── customer.domain.ts # Customer entity type
│   │   │   └── dtos/             # Data transfer objects
│   │   ├── error/                # Error handling
│   │   │   ├── appError.ts       # Custom API error
│   │   │   ├── zodError.ts       # Zod validation error handler
│   │   │   ├── validationError.ts # Mongoose validation handler
│   │   │   ├── duplicateError.ts # MongoDB duplicate key handler
│   │   │   └── castError.ts      # Mongoose cast error handler
│   │   ├── infra/db/             # Infrastructure layer
│   │   │   └── connection.ts     # MongoDB connection singleton
│   │   ├── models/               # Mongoose schemas
│   │   │   ├── user.model.ts     # User collection schema
│   │   │   └── customer.model.ts # Customer collection schema
│   │   ├── repo/                 # Data access layer
│   │   │   └── userCustomer.repo.ts # CRUD operations
│   │   ├── service/              # Business logic layer
│   │   │   └── userCustomer.service.ts # User-Customer logic
│   │   ├── types/                # Shared TypeScript types
│   │   └── utils/                # Helper functions
│   │       ├── catchAsync.ts     # Async error wrapper
│   │       ├── hashPassword.ts   # bcrypt password hashing
│   │       ├── sendResponse.ts   # Standardized response
│   │       └── dateToTimeStamp.ts # Date to protobuf timestamp
│   └── proto/                    # Protocol Buffer definitions
│       ├── user.proto            # User service contract
│       ├── buf.gen.yaml          # Buf code generation config
│       └── gen/                  # Generated proto code
├── .env                          # Environment variables
├── tsconfig.json                 # TypeScript configuration
├── eslint.config.cjs             # ESLint configuration
├── .prettierrc.json              # Prettier configuration
└── package.json                  # Dependencies & scripts

```

## 🔐 Data Model

### Dual Entity Architecture

The service uses a **User-Customer** split model:

1. **User** (Authentication Entity)
   - Stores credentials and auth-related data
   - Fields: `email`, `password`, `role`, `status`, `isDeleted`, `passwordChangedAt`
   - Roles: `customer`, `admin`, `superAdmin`
   - Status: `in-progress`, `blocked`

2. **Customer** (Profile Entity)
   - Stores customer profile information
   - Fields: `name`, `email`, `phone`, `address`, `avatarUrl`, `user` (reference)
   - One-to-one relationship with User

### Database Schema

**Users Collection:**

```typescript
{
  _id: ObjectId,
  email: String (unique, required),
  password: String (hashed, required),
  role: "customer" | "admin" | "superAdmin",
  status: "in-progress" | "blocked",
  isDeleted: Boolean (default: false),
  passwordChangedAt?: Date,
  createdAt: Date,
  updatedAt: Date
}
```

**Customers Collection:**

```typescript
{
  _id: ObjectId,
  name: String (required),
  email: String (unique, required),
  user: ObjectId (ref: "User", unique),
  phone?: String,
  address?: String,
  avatarUrl?: String,
  createdAt: Date,
  updatedAt: Date
}
```

### Transaction Support

Account creation uses **MongoDB transactions** to ensure atomicity:

1. Create User (auth credentials)
2. Create Customer (profile linked to User)
3. Commit or rollback both operations together

## 📡 API Interfaces

### gRPC Service Definition

```protobuf
service UserService {
  rpc CreateCustomer(CreateCustomerRequest) returns (CreateCustomerResponse);
  rpc GetCustomerByEmail(GetCustomerByEmailRequest) returns (CreateCustomerResponse);
  rpc GetCustomers(GetCustomersRequest) returns (GetCustomersResponse);
  rpc DeleteCustomer(DeleteCustomerRequest) returns (DeleteCustomerResponse);
  rpc GetCustomerCredentials(GetCustomerByEmailRequest) returns (CustomerCredentialsResponse);
}
```

### Endpoints

#### 1. Create Customer

**gRPC**: `CreateCustomer`  
**REST**: `POST /api/v1/customers`

Creates a new customer account with user credentials.

**Request**:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123",
  "phone": "+1234567890",
  "address": "123 Main St, City",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

**Response**:

```json
{
  "success": true,
  "message": "Customer created successfully",
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "role": "customer",
    "status": "in-progress",
    "phone": "+1234567890",
    "address": "123 Main St, City",
    "avatar_url": "https://example.com/avatar.jpg",
    "is_deleted": false
  }
}
```

**Validations**:

- `name`: min_len = 1
- `email`: valid email format
- `password`: min_len = 6
- `phone`, `address`, `avatar_url`: optional

**Process**:

1. Validate input data (Zod + protoc-gen-validate)
2. Start MongoDB transaction
3. Hash password with bcrypt
4. Create User document
5. Create Customer document (linked to User)
6. Commit transaction
7. Return sanitized response (no password)

**Error Handling**:

- `409 Conflict`: Email already exists
- `400 Bad Request`: Validation errors
- `500 Internal`: Transaction failed

---

#### 2. Get Customer By Email

**gRPC**: `GetCustomerByEmail`  
**REST**: `GET /api/v1/customers/:email`

Retrieves customer profile by email address.

**Request**:

```protobuf
message GetCustomerByEmailRequest {
  string email = 1; // Must be valid email format
}
```

**Response**:

```json
{
  "success": true,
  "message": "Customer fetched successfully",
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "role": "customer",
    "status": "in-progress",
    "phone": "+1234567890",
    "address": "123 Main St, City",
    "avatar_url": "https://example.com/avatar.jpg",
    "is_deleted": false
  }
}
```

**Error Handling**:

- `404 Not Found`: Customer not found or deleted
- `400 Bad Request`: Invalid email format

---

#### 3. Get All Customers

**gRPC**: `GetCustomers`  
**REST**: `GET /api/v1/customers`

Retrieves all customer profiles (excludes deleted).

**Request**: Empty

**Response**:

```json
{
  "success": true,
  "message": "Customers fetched successfully",
  "data": {
    "customers": [
      {
        "name": "John Doe",
        "email": "john@example.com",
        "role": "customer",
        "status": "in-progress",
        "is_deleted": false
      },
      {
        "name": "Jane Smith",
        "email": "jane@example.com",
        "role": "customer",
        "status": "in-progress",
        "is_deleted": false
      }
    ]
  }
}
```

**Features**:

- Populates user data (role, status)
- Filters out sensitive information (password)
- Includes optional fields only if present

---

#### 4. Delete Customer (Soft Delete)

**gRPC**: `DeleteCustomer`  
**REST**: `DELETE /api/v1/customers/:email`

Soft deletes a customer by setting `isDeleted` flag.

**Request**:

```json
{
  "email": "john@example.com"
}
```

**Response**:

```json
{
  "success": true,
  "message": "user deleted successfully",
  "data": {
    "msg": "user deleted successfully"
  }
}
```

**Notes**:

- Does NOT physically delete data
- Sets `user.isDeleted = true`
- Customer profile remains in database
- Deleted users cannot log in

---

#### 5. Get Customer Credentials (Internal)

**gRPC**: `GetCustomerCredentials`

Retrieves hashed password and auth metadata for login verification.

**⚠️ Internal Use Only** - Called by Auth Service

**Request**:

```protobuf
message GetCustomerByEmailRequest {
  string email = 1;
}
```

**Response**:

```protobuf
message CustomerCredentialsResponse {
  string email = 1;
  string password = 2;  // bcrypt hashed
  string status = 3;
  string role = 4;
  bool is_deleted = 5;
  google.protobuf.Timestamp password_changed_at = 6;
}
```

**Use Case**:

- Auth Service calls this during login
- Verifies password with bcrypt.compare()
- Checks if account is deleted/blocked

---

## 🔌 Integration Points

### Called By (Consumers)

1. **Auth Service** (gRPC Client)
   - `CreateCustomer` - User registration
   - `GetCustomerByEmail` - Profile lookup
   - `GetCustomerCredentials` - Login verification
   - `GetCustomers` - Admin user list
   - `DeleteCustomer` - Account deletion

2. **Kong API Gateway** (HTTP → gRPC)
   - Routes REST requests to gRPC handlers
   - Applies authentication/authorization plugins
   - Rate limiting and logging

### External Dependencies

- **MongoDB** (port 27017)
  - Primary data store
  - Requires connection string in `DB_URL`

## ⚙️ Configuration

### Environment Variables

Create a `.env` file:

```bash
# HTTP Server
PORT=3001
NODE_ENV=development

# gRPC Server
GRPC_PORT=5001

# Database
DB_URL=mongodb://localhost:27017/user_service

# Password Hashing
SALTROUND=10
```

### Port Allocation

- **3001**: REST API (Express)
- **5001**: gRPC server
- **27017**: MongoDB (external)

## 🚀 Getting Started

### Prerequisites

- Node.js 18+ (LTS recommended)
- MongoDB 6.0+
- Protocol Buffers compiler (`protoc`)
- Buf CLI (for proto generation)

### Installation

1. **Install dependencies**

   ```bash
   cd user_service
   npm install
   ```

2. **Install development tools**

   ```bash
   # Buf for proto generation
   npm install -g @bufbuild/buf

   # Or use npx for one-time use
   npx @bufbuild/buf --version
   ```

3. **Generate proto files**

   ```bash
   cd src/proto
   buf generate
   cd ../..
   ```

4. **Set up environment**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Start MongoDB**

   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest

   # Or start local MongoDB service
   sudo systemctl start mongodb
   ```

6. **Run the service**

   ```bash
   # Development (with hot reload)
   npm run start:dev

   # Production build
   npm run build
   npm run start:prod

   # Linting
   npm run lint
   npm run lint:fix

   # Code formatting
   npm run format
   ```

### Testing

**Test with cURL (REST API)**:

```bash
# Create customer
curl -X POST http://localhost:3001/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "phone": "+1234567890"
  }'

# Get customer by email
curl http://localhost:3001/api/v1/customers/test@example.com

# Get all customers
curl http://localhost:3001/api/v1/customers

# Delete customer
curl -X DELETE http://localhost:3001/api/v1/customers/test@example.com
```

**Test with grpcurl (gRPC)**:

```bash
# List services
grpcurl -plaintext localhost:5001 list

# Create customer
grpcurl -plaintext -d '{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password123"
}' localhost:5001 user_service.UserService/CreateCustomer

# Get customer by email
grpcurl -plaintext -d '{
  "email": "test@example.com"
}' localhost:5001 user_service.UserService/GetCustomerByEmail

# Get customer credentials (internal)
grpcurl -plaintext -d '{
  "email": "test@example.com"
}' localhost:5001 user_service.UserService/GetCustomerCredentials
```

## 🏛️ Design Patterns

### Layered Architecture

```
┌─────────────────────────────────────────┐
│  Presentation Layer                     │
│  • REST API Handlers                    │
│  • gRPC Handlers                        │
│  • Request Validation                   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  Service Layer                          │
│  • Business Logic                       │
│  • Transaction Management               │
│  • Data Transformation                  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  Repository Layer                       │
│  • Database Operations                  │
│  • Mongoose Queries                     │
│  • Session Management                   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  Data Layer                             │
│  • MongoDB Collections                  │
│  • Mongoose Models                      │
└─────────────────────────────────────────┘
```

### Dependency Injection

Constructor-based DI for loose coupling:

```typescript
// main.ts
const repo = new UserCustomerRepo(UserModel, CustomerModel);
const service = new UserCustomerService(repo);
const handler = new UserCustomerHandler(service);
const grpcHandler = new UserCustomerGrpcHandler(service);
```

### Repository Pattern

- Abstracts database operations
- Interface-based contracts
- Easy to mock for testing
- Supports multiple data sources

### Error Handling Strategy

**Centralized Error Mapping**:

- Zod validation errors → `INVALID_ARGUMENT`
- Mongoose validation errors → `INVALID_ARGUMENT`
- Mongoose cast errors → `INVALID_ARGUMENT`
- Duplicate key errors → `ALREADY_EXISTS`
- Custom ApiError → HTTP → gRPC status codes
- Unknown errors → `UNKNOWN`

**gRPC Status Code Mapping**:

| HTTP | gRPC Status       | Use Case                    |
| ---- | ----------------- | --------------------------- |
| 400  | INVALID_ARGUMENT  | Bad request data            |
| 401  | UNAUTHENTICATED   | Missing/invalid credentials |
| 403  | PERMISSION_DENIED | Forbidden action            |
| 404  | NOT_FOUND         | Resource not found          |
| 409  | ALREADY_EXISTS    | Duplicate email             |
| 500  | INTERNAL          | Server errors               |

## 📊 Data Flow Diagrams

### Create Customer Flow

```
┌──────────┐         ┌──────────┐         ┌─────────┐         ┌──────────┐
│  Client  │────────→│ Handler  │────────→│ Service │────────→│   Repo   │
└──────────┘  POST   └──────────┘  create └─────────┘  create └──────────┘
    /customers           │                     │                    │
                         │                     │         ┌──────────▼────────┐
                         │                     │         │ Start Transaction │
                         │                     │         └──────────┬────────┘
                         │                     │                    │
                         │                     │         ┌──────────▼────────┐
                         │                     │         │  Hash Password    │
                         │                     │         └──────────┬────────┘
                         │                     │                    │
                         │                     │         ┌──────────▼────────┐
                         │                     │         │  Create User      │
                         │                     │         └──────────┬────────┘
                         │                     │                    │
                         │                     │         ┌──────────▼────────┐
                         │                     │         │ Create Customer   │
                         │                     │         └──────────┬────────┘
                         │                     │                    │
                         │                     │         ┌──────────▼────────┐
                         │                     │         │Commit Transaction │
                         │                     │         └──────────┬────────┘
                         │                     │◄───────────────────┘
                         │◄────────────────────┘
                         │
    ┌────────────────────▼───────────────────┐
    │     Return Response (no password)      │
    └────────────────────────────────────────┘
```

### Get Customer Flow

```
Client → Handler → Service → Repo → MongoDB
   │                              ↓
   │                    Find customer by email
   │                              ↓
   │                    Populate user fields
   │                              ↓
   │                    Filter sensitive data
   │                              ↑
   └──────────────────────────────┘
         Return sanitized profile
```

## 🛡️ Security Considerations

### Password Security

- **bcrypt hashing** with configurable salt rounds
- Passwords NEVER returned in responses
- Only hashed password in `GetCustomerCredentials`

### Data Sanitization

- Sensitive fields filtered in responses
- No password exposure in REST/gRPC responses
- Deleted users cannot authenticate

### Input Validation

- **Dual validation**: Zod (TypeScript) + protoc-gen-validate (protobuf)
- Email format validation
- Password minimum length enforcement
- SQL injection protection via Mongoose

### Soft Delete Pattern

- Preserves data for audit trails
- Prevents accidental data loss
- Allows account recovery

## 📈 Performance Optimizations

### Database Connection

- **Singleton pattern** prevents connection leaks
- Connection pooling via Mongoose
- Automatic reconnection handling

### Transaction Management

- Atomic operations for data consistency
- Rollback on partial failures
- Session cleanup in finally blocks

### Query Optimization

- Selective field population (`select`)
- Lean queries for better performance
- Index on `email` field (unique constraint)

### Error Handling

- Early validation to reduce DB load
- Structured error responses
- Graceful degradation

## 🧩 Integration with Kong Gateway

The service works behind Kong API Gateway with these plugins:

1. **grpc-gateway**: HTTP → gRPC translation
2. **auth-token-validator**: JWT validation (on create endpoint)
3. **rate-limiting**: Prevent abuse
4. **user-context-injector**: Injects user metadata

**Kong Route Configuration**:

```yaml
- name: create-user-route
  paths: ['/users']
  methods: [POST]
  plugins:
    - name: auth-token-validator
    - name: rate-limiting

- name: get-all-user-routes
  paths: ['/customers']
  methods: [GET]
  plugins:
    - name: auth-token-validator
```

## 🔮 Future Enhancements

- [ ] Email verification flow
- [ ] Phone number verification (OTP)
- [ ] Profile picture upload (S3 integration)
- [ ] User preferences and settings
- [ ] Activity logging and audit trail
- [ ] Admin role management
- [ ] Customer search and filtering
- [ ] Pagination for customer list
- [ ] Redis caching for frequently accessed profiles
- [ ] Metrics and monitoring (Prometheus)
- [ ] Distributed tracing (OpenTelemetry)
- [ ] GraphQL API support

## 🧪 Testing Strategy

### Unit Tests (TODO)

- Service layer logic
- Repository methods
- Utility functions

### Integration Tests (TODO)

- gRPC endpoint testing
- REST API testing
- Database transaction scenarios

### E2E Tests (TODO)

- Full customer lifecycle
- Error scenarios
- Authentication integration

## 📝 API Response Format

### Success Response

```typescript
{
  success: true,
  message: "Operation successful",
  data: { ... }
}
```

### Error Response

```typescript
{
  success: false,
  message: "Error description",
  errorMessages: [
    {
      path: "field_name",
      message: "Validation error"
    }
  ],
  stack: "..." // Only in development
}
```

## 🐛 Error Types

1. **ZodError**: Input validation failures
2. **ValidationError**: Mongoose schema validation
3. **CastError**: Invalid ObjectId format
4. **DuplicateError**: Unique constraint violation (email)
5. **ApiError**: Custom business logic errors

## 📦 Build & Deployment

### Development

```bash
npm run start:dev  # tsx watch for hot reload
```

### Production

```bash
npm run build      # Compile TypeScript → JavaScript
npm run start:prod # Run compiled code
```

### Docker (Future)

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY dist ./dist
EXPOSE 3001 5001
CMD ["node", "dist/cmd/api/main.js"]
```

## 📖 License

This project is part of a larger e-commerce microservices architecture.

## 👥 Contributing

This is a learning/portfolio project demonstrating:

- Clean architecture principles
- TypeScript best practices
- gRPC implementation in Node.js
- MongoDB transaction handling
- Dual interface design (REST + gRPC)

---

**Built with ❤️ using TypeScript, Express, gRPC, and MongoDB**
