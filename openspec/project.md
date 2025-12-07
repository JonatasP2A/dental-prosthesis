# Project Context

## Purpose
A SaaS platform for dental prosthesis laboratories to manage their operations, including order management, client relationships, production workflows, and billing. The platform streamlines communication between dental clinics and prosthesis labs, enabling efficient tracking of prosthetic work from order to delivery.

## Tech Stack

### Backend
- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Configuration**: Viper
- **Authentication**: Clerk (JWT validation)
- **Architecture**: Hexagonal Architecture (Ports & Adapters)

### Frontend
- **Framework**: Next.js (App Router)
- **UI Components**: shadcn/ui
- **Language**: TypeScript
- **Authentication**: Clerk (@clerk/nextjs)
- **Styling**: Tailwind CSS (via shadcn)

### Shared
- **Authentication Provider**: Clerk
- **API Communication**: REST (JSON)

## Project Conventions

### Code Style

#### Backend (Go)
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Package names: lowercase, single word when possible
- Interface names: verb-er pattern (e.g., `Reader`, `OrderCreator`)
- Private types/functions: lowercase first letter
- Public types/functions: uppercase first letter (PascalCase)
- Constants: `UPPER_SNAKE_CASE` for exported, `camelCase` for unexported
- Error handling: return errors explicitly, no panic in business logic
- Use structured logging

#### Frontend (TypeScript/Next.js)
- Use TypeScript strict mode
- Components: PascalCase (e.g., `OrderForm.tsx`)
- Hooks: camelCase with `use` prefix (e.g., `useOrders.ts`)
- Utilities: camelCase (e.g., `formatCurrency.ts`)
- Use `interface` for object shapes, `type` for unions/primitives
- Prefer named exports over default exports
- Use Prettier for formatting
- ESLint for linting

### Architecture Patterns

#### Backend - Hexagonal Architecture
```
backend/
├── cmd/                    # Application entry points
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/             # Core business logic (entities, value objects)
│   │   ├── order/
│   │   ├── client/
│   │   └── prosthesis/
│   ├── ports/              # Interface definitions
│   │   ├── inbound/        # Use case interfaces (driving)
│   │   └── outbound/       # Repository/service interfaces (driven)
│   ├── adapters/           # Implementations
│   │   ├── inbound/        # HTTP handlers, gRPC, CLI
│   │   │   └── http/
│   │   └── outbound/       # Database, external APIs
│   │       ├── persistence/
│   │       └── external/
│   ├── application/        # Use cases / application services
│   └── config/             # Viper configuration
├── pkg/                    # Shared packages
└── test/                   # Integration tests
```

**Key Principles:**
- Domain layer has NO external dependencies
- Ports define contracts, adapters implement them
- Use cases orchestrate domain logic
- Dependency injection via constructor functions

#### Frontend - Feature-Based Structure
```
frontend/
├── app/                    # Next.js App Router
│   ├── (auth)/             # Auth-required routes
│   ├── (public)/           # Public routes
│   └── api/                # API routes (if needed)
├── components/
│   ├── ui/                 # shadcn components
│   └── [feature]/          # Feature-specific components
├── lib/                    # Utilities, API clients
├── hooks/                  # Custom React hooks
├── types/                  # TypeScript types/interfaces
└── services/               # API service layer
```

### Testing Strategy

#### Backend
- **Minimum Coverage**: 90% for unit and integration tests combined
- **Unit Tests**: Test domain logic and use cases in isolation
  - Use table-driven tests
  - Mock external dependencies via interfaces
  - File naming: `*_test.go` in same package
- **Integration Tests**: Test adapters with real dependencies
  - Use test containers for database tests
  - Located in `test/` directory
- **Tools**: Go's built-in testing, testify for assertions

#### Frontend
- **Unit Tests**: Components and hooks (Jest + React Testing Library)
- **E2E Tests**: Critical user flows (Playwright or Cypress)
- **Coverage Target**: 80%+ for critical paths

### Git Workflow
- **Branching Strategy**: GitHub Flow
  - `main` - production-ready code
  - `feature/*` - new features
  - `fix/*` - bug fixes
  - `refactor/*` - code improvements
- **Commit Convention**: Conventional Commits
  - `feat:` - new feature
  - `fix:` - bug fix
  - `refactor:` - code refactoring
  - `test:` - adding tests
  - `docs:` - documentation
  - `chore:` - maintenance tasks
- **PR Requirements**:
  - All tests passing
  - Code review approval
  - No linting errors

## Domain Context

### Key Entities
- **Laboratory**: The dental prosthesis lab (tenant)
- **Client**: Dental clinics/dentists who order prosthetic work
- **Order**: A prosthesis work order from a client
- **Prosthesis**: The dental prosthetic item (crown, bridge, denture, implant, etc.)
- **Material**: Materials used (zirconia, porcelain, metal alloys, etc.)
- **Technician**: Lab staff who produce the prosthetics

### Prosthesis Types
- Crowns (coroas)
- Bridges (pontes)
- Complete dentures (próteses totais)
- Partial dentures (próteses parciais)
- Implant-supported prosthetics
- Veneers (facetas)
- Inlays/Onlays

### Workflow States
1. **Received** - Order received from client
2. **In Production** - Work in progress
3. **Quality Check** - Inspection before delivery
4. **Ready** - Completed, awaiting pickup/delivery
5. **Delivered** - Sent to client
6. **Revision** - Returned for adjustments

### Industry Terms
- **Shade/Color**: Tooth color matching (e.g., A1, A2, B1 using VITA scale)
- **Antagonist**: Opposing teeth consideration
- **Occlusion**: How teeth bite together
- **Try-in**: Fitting appointment before final delivery

## Important Constraints

### Technical
- Backend API must be stateless for horizontal scaling
- All timestamps in UTC
- Multi-tenant architecture (laboratory isolation)
- API versioning required for breaking changes

### Business
- LGPD compliance (Brazilian data protection law) for patient data
- Audit trail for all order modifications
- Financial data must be immutable (soft deletes only)

### Security
- All endpoints authenticated via Clerk JWT
- Role-based access control (Admin, Manager, Technician)
- Sensitive data encrypted at rest
- HTTPS only in production

## External Dependencies

### Authentication
- **Clerk**: User management, authentication, JWT tokens
  - Frontend: `@clerk/nextjs` SDK
  - Backend: JWT validation middleware

### Infrastructure (Future)
- **Database**: PostgreSQL (primary data store)
- **Cache**: Redis (session, caching)
- **Storage**: S3-compatible (images, documents)
- **Email**: Transactional email service (order notifications)

### Integrations (Planned)
- Payment gateway (billing/invoicing)
- WhatsApp Business API (client notifications)
