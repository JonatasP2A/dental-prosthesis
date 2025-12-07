# Design: Laboratory Domain Implementation

## Context
The Laboratory domain is the first domain to be implemented in the dental prosthesis SaaS platform. It serves as the multi-tenant root aggregate, meaning all other entities (clients, orders, technicians) belong to a laboratory. This implementation establishes the architectural patterns and conventions that will be followed for all subsequent domains.

**Constraints:**
- Multi-tenant architecture (laboratory isolation)
- Clerk authentication for JWT validation
- Hexagonal architecture (ports & adapters)
- 90%+ test coverage requirement
- Go 1.21+, Gin framework, Viper for configuration

## Goals / Non-Goals

### Goals
- Establish hexagonal architecture pattern
- Implement multi-tenancy foundation
- Create reusable authentication middleware pattern
- Set testing patterns and conventions
- Provide CRUD operations for laboratories

### Non-Goals
- Database persistence (start with in-memory, add DB adapter later)
- Complex business rules (keep domain logic simple initially)
- Laboratory user management (handled by Clerk)
- Billing/subscription management (future capability)

## Decisions

### Decision: Start with In-Memory Repository
**What**: Use an in-memory map-based repository implementation initially.
**Why**: 
- Allows development and testing without database setup
- Faster iteration on domain and use case logic
- Database adapter can be added later without changing domain/ports
- Aligns with hexagonal architecture (adapter can be swapped)

**Alternatives considered**:
- PostgreSQL from start → Rejected: Adds complexity, slows initial development
- SQLite → Rejected: Still requires DB setup, in-memory is simpler for MVP

### Decision: Clerk JWT Claims Structure
**What**: Extract laboratory ID from Clerk JWT custom claims (`laboratory_id`).
**Why**:
- Clerk handles user authentication and management
- Custom claims allow embedding laboratory context
- Middleware extracts lab ID and injects into request context
- Enables laboratory-scoped queries automatically

**Alternatives considered**:
- Separate laboratory-user mapping table → Rejected: Adds complexity, Clerk custom claims simpler
- Header-based laboratory selection → Rejected: Security risk, users could access wrong lab

### Decision: Use Case Error Handling Pattern
**What**: Return domain-specific errors from use cases, convert to HTTP status codes in adapters.
**Why**:
- Keeps domain layer independent of HTTP concerns
- Allows reuse of use cases with different adapters (CLI, gRPC, etc.)
- Clear separation of concerns

**Error types**:
- `ErrLaboratoryNotFound` → 404 Not Found
- `ErrInvalidInput` → 400 Bad Request
- `ErrUnauthorized` → 401 Unauthorized
- `ErrInternal` → 500 Internal Server Error

### Decision: API Versioning
**What**: Use `/api/v1/` prefix for all endpoints.
**Why**:
- Allows breaking changes in future versions
- Industry standard practice
- Explicit versioning prevents accidental breaking changes

## Risks / Trade-offs

### Risk: In-Memory Repository Data Loss
**Mitigation**: 
- Clearly document this is development-only
- Add database adapter as next priority after domain validation
- Use integration tests to validate repository interface contract

### Risk: Clerk Custom Claims Not Set Up
**Mitigation**:
- Document required Clerk configuration
- Provide example JWT structure
- Add validation middleware that fails gracefully with clear error

### Risk: Multi-Tenancy Bypass
**Mitigation**:
- Always extract laboratory_id from authenticated context
- Never allow laboratory_id in request body (only from JWT)
- Add integration tests that verify laboratory isolation

## Migration Plan

### Phase 1: In-Memory Implementation (Current)
- Implement domain, ports, use cases, in-memory repository
- HTTP handlers with Clerk middleware
- Full test coverage

### Phase 2: Database Adapter (Future)
- Add PostgreSQL adapter implementing same repository interface
- Migrate integration tests to use database
- Keep in-memory adapter for unit tests

### Rollback
- No production data to migrate (initial implementation)
- Can revert to in-memory if database issues occur

## Open Questions
- Should laboratory creation be restricted to super-admin role? → **Decision**: Yes, handled by Clerk role-based access
- Should we support laboratory soft-delete? → **Decision**: Yes, add `DeletedAt` timestamp field
- How to handle laboratory name uniqueness? → **Decision**: Enforce uniqueness per tenant (Clerk organization), not globally

