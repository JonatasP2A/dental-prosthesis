# Design: Remove laboratory_id from JWT Authentication

## Context
In the first stage of the product (MVP), we want to simplify authentication by removing the requirement for `laboratory_id` in JWT tokens. Currently, the system extracts `laboratory_id` from JWT custom claims and requires it for all authenticated requests. This change removes that requirement to simplify the initial product setup.

**Current State:**
- JWT middleware extracts `laboratory_id` from custom claims
- Middleware returns 403 if `laboratory_id` is missing
- All handlers use `auth.GetLaboratoryID()` to get laboratory context
- Multi-tenancy enforced through laboratory_id in context

**Constraints:**
- Must maintain backward compatibility where possible during transition
- Multi-tenancy must still be enforced (cannot be removed entirely)
- Changes affect authentication middleware, all HTTP handlers, and domain services
- Frontend will need updates to remove laboratory_id from JWT setup

## Goals / Non-Goals

### Goals
- Simplify JWT authentication setup for MVP
- Remove dependency on Clerk custom claims for laboratory context
- Maintain security and authentication integrity
- Enable alternative multi-tenancy approaches (user-laboratory relationship, explicit selection, etc.)

### Non-Goals
- Removing multi-tenancy entirely (laboratory isolation still required)
- Changing domain entities (Client and Order still have laboratory_id fields)
- Removing laboratory_id from database schema
- Changing how laboratory context is used in business logic

## Decisions

### Decision: Remove laboratory_id from JWT Validation
**What**: Remove extraction and validation of `laboratory_id` from JWT custom claims in authentication middleware.

**Why**:
- Simplifies Clerk setup (no need to configure custom claims)
- Reduces complexity for early adopters
- Allows flexibility in how laboratory context is determined (user-laboratory relationship, explicit selection, etc.)
- Aligns with MVP philosophy of starting simple

**Alternatives considered**:
- Keep laboratory_id optional in JWT → Rejected: Adds complexity, better to remove entirely
- Derive laboratory_id from user metadata → Considered: May be future approach, but not in scope for this change
- Single laboratory per user → Considered: May be future approach, but requires user-laboratory relationship table

### Decision: Make GetLaboratoryID Return Empty String
**What**: Update `GetLaboratoryID()` function to return empty string instead of removing it entirely.

**Why**:
- Maintains API compatibility for code that calls it
- Allows gradual migration of handlers
- Prevents breaking changes in dependent code
- Empty string return clearly indicates laboratory_id not available from JWT

**Alternatives considered**:
- Remove function entirely → Rejected: Would break all handlers immediately
- Return error → Rejected: Changes function signature, breaking change
- Deprecate function → Considered: Good approach, but returning empty string is simpler for MVP

### Decision: laboratory_id as Query Parameter
**What**: All handlers must extract `laboratory_id` from URL query parameters (e.g., `?laboratory_id=xxx`).

**Why**:
- Simple and explicit approach for MVP
- No route restructuring required (can add query param to existing routes)
- Frontend can easily include laboratory_id in all requests
- Clear API contract - laboratory_id is visible in URL
- Allows for future migration to path parameters or other approaches

**Implementation approach**:
- Handlers extract `laboratory_id` from `c.Query("laboratory_id")` in Gin context
- Validate that `laboratory_id` is present and non-empty
- Return HTTP 400 Bad Request if `laboratory_id` is missing or invalid
- Use extracted `laboratory_id` for all laboratory-scoped operations

**Alternatives considered**:
- Path parameters (e.g., `/api/v1/laboratories/:lab_id/clients`) → Rejected: Requires route restructuring, more complex migration
- Require laboratory_id in request body → Rejected: Security risk, users could access wrong lab, inconsistent with REST principles
- Lookup user-laboratory relationship → Considered: Requires user-laboratory mapping table (future work, not MVP)
- Keep in JWT but optional → Rejected: Adds complexity, better to remove entirely for MVP

## Risks / Trade-offs

### Risk: Multi-Tenancy Bypass
**Impact**: Users might access data from wrong laboratory if laboratory context is not properly enforced.

**Mitigation**:
- Document that handlers must still enforce laboratory scoping
- Add validation in handlers to ensure laboratory_id matches user's laboratory
- Consider adding user-laboratory relationship table in future
- Add integration tests to verify laboratory isolation

### Risk: Breaking Existing Frontend Integration
**Impact**: Frontend code that sets `laboratory_id` in JWT custom claims will need updates.

**Mitigation**:
- Update frontend documentation
- Provide migration guide
- Frontend can remove laboratory_id from Clerk custom claims setup

### Risk: Handler Implementation Inconsistency
**Impact**: Different handlers might implement laboratory context resolution differently.

**Mitigation**:
- Document expected approach for MVP (may vary by endpoint)
- Consider creating helper function for laboratory context resolution
- Add code review checklist for laboratory scoping

### Trade-off: Simplicity vs Security
**Benefit**: Simpler authentication setup
**Cost**: Less automatic laboratory context enforcement (must be handled per-handler)

## Migration Plan

### Phase 1: Backend Changes (This Change)
1. Remove laboratory_id extraction from JWT middleware
2. Remove laboratory_id validation check
3. Update `GetLaboratoryID()` to return empty string
4. Update handlers to handle missing laboratory_id
5. Update tests

### Phase 2: Frontend Updates (Future)
1. Remove laboratory_id from Clerk custom claims configuration
2. Update frontend code that references laboratory_id from JWT
3. Implement alternative laboratory context mechanism (if needed)

### Phase 3: Alternative Multi-Tenancy (Future)
1. Implement user-laboratory relationship table
2. Add lookup service for user's laboratory
3. Update handlers to use relationship lookup
4. Remove need for laboratory_id in requests entirely

### Rollback
- Can revert middleware changes if issues arise
- Frontend can re-add laboratory_id to JWT if needed
- No database migrations required (schema unchanged)

## Open Questions
- How will laboratory context be determined in MVP? → **Decision**: Query parameter `?laboratory_id=xxx` in all API requests
- Should we add user-laboratory relationship table now? → **Decision**: No, defer to future change
- Should laboratory_id be required in request body for certain operations? → **Decision**: No, query parameter only for consistency
- How to handle users who belong to multiple laboratories? → **Decision**: Out of scope for MVP, users specify laboratory_id explicitly per request
