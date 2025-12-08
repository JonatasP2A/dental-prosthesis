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
│   │   ├── order/       # Order domain (planned)
│   │   ├── client/      # Client domain (planned)
│   │   ├── prosthesis/  # Prosthesis domain (planned)
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
│   │   └── laboratory/  # Laboratory use cases
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
3. Add `laboratory_id` to user's public metadata or JWT claims
4. Include the JWT token in the `Authorization` header:
   ```
   Authorization: Bearer <your-jwt-token>
   ```

### JWT Claims
The middleware expects the following claims:
- `sub` - User ID
- `laboratory_id` - Laboratory ID (custom claim)
- `exp` - Expiration timestamp
