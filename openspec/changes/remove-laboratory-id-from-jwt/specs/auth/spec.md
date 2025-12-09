## MODIFIED Requirements

### Requirement: JWT Authentication Middleware
The system SHALL authenticate requests using Clerk JWT tokens via the Authorization header. The middleware SHALL extract and validate the JWT token, verify its signature using Clerk's JWKS endpoint, and inject authenticated user information into the request context.

#### Scenario: Valid JWT token authentication
- **WHEN** a request includes a valid Bearer token in the Authorization header
- **THEN** the middleware validates the token signature
- **AND** extracts the user ID (`sub` claim) from the token
- **AND** injects the user ID into the request context
- **AND** allows the request to proceed

#### Scenario: Missing authorization header
- **WHEN** a request does not include an Authorization header
- **THEN** the middleware returns HTTP 401 Unauthorized
- **AND** includes an error message indicating missing authorization header

#### Scenario: Invalid token format
- **WHEN** a request includes an Authorization header with invalid format (not "Bearer <token>")
- **THEN** the middleware returns HTTP 401 Unauthorized
- **AND** includes an error message indicating invalid authorization header format

#### Scenario: Invalid or expired token
- **WHEN** a request includes an invalid, expired, or tampered JWT token
- **THEN** the middleware returns HTTP 401 Unauthorized
- **AND** includes an error message indicating token verification failed

## REMOVED Requirements

### Requirement: Laboratory ID from JWT
**Reason**: In the first stage of the product, laboratory context will not be extracted from JWT tokens. This simplifies authentication setup and allows for explicit laboratory selection via URL query parameters in the MVP.

**Migration**: Handlers and services that previously relied on `laboratory_id` from JWT context must now extract `laboratory_id` from URL query parameters.

## ADDED Requirements

### Requirement: Laboratory ID Query Parameter
The system SHALL require `laboratory_id` as a query parameter in all API requests that require laboratory context.

#### Scenario: Valid request with laboratory_id query parameter
- **WHEN** a request includes `laboratory_id` query parameter (e.g., `GET /api/v1/clients?laboratory_id=lab-123`)
- **THEN** the handler extracts the `laboratory_id` from the query parameter
- **AND** uses it for laboratory-scoped operations
- **AND** allows the request to proceed

#### Scenario: Missing laboratory_id query parameter
- **WHEN** a request does not include `laboratory_id` query parameter
- **THEN** the handler returns HTTP 400 Bad Request
- **AND** includes an error message indicating `laboratory_id` query parameter is required

#### Scenario: Empty laboratory_id query parameter
- **WHEN** a request includes `laboratory_id` query parameter with empty value (e.g., `?laboratory_id=`)
- **THEN** the handler returns HTTP 400 Bad Request
- **AND** includes an error message indicating `laboratory_id` is required
