# Backend - Dental Prosthesis Laboratory SaaS

Backend service built with Go, following Hexagonal Architecture (Ports & Adapters) pattern.

## Architecture

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/          # Core business logic (entities, value objects)
│   │   ├── errors/      # Domain errors
│   │   ├── laboratory/  # Laboratory domain
│   │   ├── order/       # Order domain
│   │   ├── client/       # Client domain
│   │   ├── prosthesis/  # Prosthesis domain
│   │   └── technician/  # Technician domain (planned)
│   ├── ports/           # Interface definitions
│   │   ├── inbound/     # Use case interfaces (driving)
│   │   └── outbound/    # Repository/service interfaces (driven)
│   ├── adapters/        # Implementations
│   │   ├── inbound/     # HTTP handlers
│   │   │   └── http/
│   │   │       ├── dto/      # Request/Response DTOs
│   │   │       ├── handler/  # HTTP handlers
│   │   │       └── router/   # Gin router setup
│   │   └── outbound/    # Database, external APIs
│   │       └── persistence/
│   │           └── memory/   # In-memory repository
│   ├── application/     # Use cases / application services
│   │   ├── laboratory/  # Laboratory use cases
│   │   ├── client/      # Client use cases
│   │   ├── order/       # Order use cases
│   │   └── prosthesis/ # Prosthesis use cases
│   └── config/          # Viper configuration
├── pkg/                 # Shared packages
│   ├── auth/            # Clerk authentication
│   └── uuid/            # UUID generation
└── test/                # Integration tests
```

## Getting Started

### Prerequisites
- Go 1.21 or higher

### Configuration

1. Copy the example configuration:
```bash
cp config.example.yaml config.yaml
```

2. Update the configuration with your Clerk credentials:
```yaml
server:
  port: "8080"
  host: "0.0.0.0"

clerk:
  secret_key: "sk_test_your_secret_key"
```

Or use environment variables:
```bash
export CLERK_SECRET_KEY="sk_test_your_secret_key"
```

Note: The Clerk SDK automatically fetches JWKS (JSON Web Key Set) using the secret key, so `jwks_url` is no longer needed.

### Run Locally

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### API Endpoints

#### Health Check
```
GET /health
```

#### Laboratories
```
POST   /api/v1/laboratories     # Create laboratory
GET    /api/v1/laboratories     # List laboratories
GET    /api/v1/laboratories/:id # Get laboratory by ID
PUT    /api/v1/laboratories/:id # Update laboratory
DELETE /api/v1/laboratories/:id # Delete laboratory (soft delete)
```

#### Clients
```
POST   /api/v1/clients?laboratory_id=xxx     # Create client
GET    /api/v1/clients?laboratory_id=xxx     # List clients
GET    /api/v1/clients/:id?laboratory_id=xxx # Get client by ID
PUT    /api/v1/clients/:id?laboratory_id=xxx # Update client
DELETE /api/v1/clients/:id?laboratory_id=xxx # Delete client (soft delete)
```

#### Orders
```
POST   /api/v1/orders?laboratory_id=xxx           # Create order
GET    /api/v1/orders?laboratory_id=xxx           # List orders
GET    /api/v1/orders/:id?laboratory_id=xxx        # Get order by ID
PUT    /api/v1/orders/:id?laboratory_id=xxx        # Update order
PATCH  /api/v1/orders/:id/status?laboratory_id=xxx # Update order status
DELETE /api/v1/orders/:id?laboratory_id=xxx        # Delete order (soft delete)
GET    /api/v1/clients/:id/orders?laboratory_id=xxx # List orders by client
```

#### Prostheses
```
POST   /api/v1/prostheses?laboratory_id=xxx     # Create prosthesis
GET    /api/v1/prostheses?laboratory_id=xxx     # List prostheses
GET    /api/v1/prostheses?laboratory_id=xxx&type=crown # List prostheses filtered by type
GET    /api/v1/prostheses?laboratory_id=xxx&material=zirconia # List prostheses filtered by material
GET    /api/v1/prostheses/:id?laboratory_id=xxx # Get prosthesis by ID
PUT    /api/v1/prostheses/:id?laboratory_id=xxx # Update prosthesis
DELETE /api/v1/prostheses/:id?laboratory_id=xxx # Delete prosthesis (soft delete)
```

#### Example: Create Laboratory
```bash
curl -X POST http://localhost:8080/api/v1/laboratories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-clerk-jwt>" \
  -d '{
    "name": "Dental Lab São Paulo",
    "email": "contact@dentallab.com",
    "phone": "+5511999999999",
    "address": {
      "street": "Rua das Flores, 123",
      "city": "São Paulo",
      "state": "SP",
      "postal_code": "01234-567",
      "country": "Brazil"
    }
  }'
```

#### Example: Create Client (requires laboratory_id query parameter)
```bash
curl -X POST "http://localhost:8080/api/v1/clients?laboratory_id=lab-123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-clerk-jwt>" \
  -d '{
    "name": "Dr. João Silva",
    "email": "joao@clinicadental.com",
    "phone": "+5511999999999",
    "address": {
      "street": "Av. Paulista, 1000",
      "city": "São Paulo",
      "state": "SP",
      "postal_code": "01310-100",
      "country": "Brazil"
    }
  }'
```

#### Example: List Clients (requires laboratory_id query parameter)
```bash
curl -X GET "http://localhost:8080/api/v1/clients?laboratory_id=lab-123" \
  -H "Authorization: Bearer <your-clerk-jwt>"
```

#### Example: Create Prosthesis (requires laboratory_id query parameter)
```bash
curl -X POST "http://localhost:8080/api/v1/prostheses?laboratory_id=lab-123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-clerk-jwt>" \
  -d '{
    "type": "crown",
    "material": "zirconia",
    "shade": "A1",
    "specifications": "Full coverage crown",
    "notes": "High priority case"
  }'
```

#### Example: List Prostheses (requires laboratory_id query parameter)
```bash
# List all prostheses
curl -X GET "http://localhost:8080/api/v1/prostheses?laboratory_id=lab-123" \
  -H "Authorization: Bearer <your-clerk-jwt>"

# List prostheses filtered by type
curl -X GET "http://localhost:8080/api/v1/prostheses?laboratory_id=lab-123&type=crown" \
  -H "Authorization: Bearer <your-clerk-jwt>"

# List prostheses filtered by material
curl -X GET "http://localhost:8080/api/v1/prostheses?laboratory_id=lab-123&material=zirconia" \
  -H "Authorization: Bearer <your-clerk-jwt>"
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report (target: 90%)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/domain/laboratory/...
go test ./internal/application/laboratory/...
```

## Tech Stack

- **Framework**: Gin
- **Configuration**: Viper
- **Authentication**: Clerk (JWT validation)
- **Architecture**: Hexagonal Architecture (Ports & Adapters)

## Domain: Laboratory

The Laboratory domain is the foundational multi-tenant entity. All other entities belong to a laboratory.

### Entity Fields
- `ID` - Unique identifier
- `Name` - Laboratory name (required, max 200 chars)
- `Email` - Contact email (required, valid format)
- `Phone` - Phone number (required, E.164 format)
- `Address` - Full address (street, city, state, postal code, country)
- `CreatedAt` - Creation timestamp (UTC)
- `UpdatedAt` - Last update timestamp (UTC)
- `DeletedAt` - Soft delete timestamp (nullable)

### Use Cases
- **CreateLaboratory**: Creates a new laboratory with validation
- **GetLaboratory**: Retrieves a laboratory by ID
- **UpdateLaboratory**: Updates laboratory information
- **ListLaboratories**: Lists all active laboratories
- **DeleteLaboratory**: Soft deletes a laboratory

## Development Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Write table-driven tests
- Mock external dependencies via interfaces
- **Minimum 90% test coverage** for unit and integration tests

## Clerk Authentication

The API uses Clerk for authentication. To access protected endpoints:

1. Set up a Clerk application at [clerk.com](https://clerk.com)
2. Configure your frontend with `@clerk/nextjs`
3. Include the JWT token in the `Authorization` header:
   ```
   Authorization: Bearer <your-jwt-token>
   ```
4. **Important**: All API requests that require laboratory context must include `laboratory_id` as a query parameter:
   ```
   GET /api/v1/clients?laboratory_id=lab-123
   POST /api/v1/orders?laboratory_id=lab-123
   ```

### JWT Claims
The middleware expects the following claims:
- `sub` - User ID
- `exp` - Expiration timestamp

### Laboratory Context
**Breaking Change**: `laboratory_id` is no longer extracted from JWT tokens. Instead, it must be provided as a query parameter in all API requests that require laboratory context:

- **Client endpoints**: All client operations require `?laboratory_id=xxx`
- **Order endpoints**: All order operations require `?laboratory_id=xxx`
- **Prosthesis endpoints**: All prosthesis operations require `?laboratory_id=xxx`
- **Laboratory endpoints**: Laboratory endpoints do not require `laboratory_id` query parameter

**Example:**
```bash
# List clients for a laboratory
GET /api/v1/clients?laboratory_id=lab-123
Authorization: Bearer <your-jwt-token>

# Create an order
POST /api/v1/orders?laboratory_id=lab-123
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "client_id": "client-456",
  "prosthesis": [...]
}
```

If `laboratory_id` query parameter is missing, the API will return HTTP 400 Bad Request.

## Domain: Prosthesis

The Prosthesis domain represents individual dental prosthetic items (crowns, bridges, dentures, implants, etc.) that can be tracked independently from orders.

### Entity Fields
- `ID` - Unique identifier
- `LaboratoryID` - Laboratory identifier (required)
- `Type` - Prosthesis type (required, enum: crown, bridge, complete_denture, partial_denture, implant, veneer, inlay, onlay)
- `Material` - Material used (required, e.g., zirconia, porcelain, metal alloy, acrylic)
- `Shade` - Tooth color matching (optional, e.g., A1, A2, B1 using VITA scale)
- `Specifications` - Technical details (optional)
- `Notes` - Special instructions (optional)
- `CreatedAt` - Creation timestamp (UTC)
- `UpdatedAt` - Last update timestamp (UTC)
- `DeletedAt` - Soft delete timestamp (nullable)

### Prosthesis Types
- `crown` - Crown (coroa)
- `bridge` - Bridge (ponte)
- `complete_denture` - Complete denture (prótese total)
- `partial_denture` - Partial denture (prótese parcial)
- `implant` - Implant-supported prosthetic
- `veneer` - Veneer (faceta)
- `inlay` - Inlay
- `onlay` - Onlay

### Use Cases
- **CreateProsthesis**: Creates a new prosthesis with validation
- **GetProsthesis**: Retrieves a prosthesis by ID (laboratory-scoped)
- **UpdateProsthesis**: Updates prosthesis information
- **ListProstheses**: Lists all active prostheses for a laboratory, optionally filtered by type or material
- **DeleteProsthesis**: Soft deletes a prosthesis

### Filtering
The ListProstheses endpoint supports optional query parameters:
- `type` - Filter by prosthesis type (e.g., `?type=crown`)
- `material` - Filter by material (e.g., `?material=zirconia`)
