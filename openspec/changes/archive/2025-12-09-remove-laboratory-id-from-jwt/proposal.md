# Change: Remove laboratory_id from JWT Authentication

## Why
In the first stage of the product, we want to simplify authentication by removing the requirement for `laboratory_id` in JWT tokens. This simplifies the initial setup and reduces complexity for early adopters. Multi-tenancy will be handled differently in the MVP stage, potentially through a single laboratory per user or through explicit laboratory selection in the application flow rather than requiring it in every authenticated request.

## What Changes
- **BREAKING** Remove requirement for `laboratory_id` custom claim from JWT tokens
- **MODIFIED** Authentication middleware no longer extracts or validates `laboratory_id` from JWT
- **MODIFIED** Remove `laboratory_id` from request context injection
- **MODIFIED** `laboratory_id` must now be provided as a query parameter in API requests (e.g., `?laboratory_id=xxx`)
- **MODIFIED** Update all handlers to extract `laboratory_id` from URL query parameters instead of context
- **MODIFIED** Update route handlers to validate `laboratory_id` query parameter is present and valid
- **REMOVED** `GetLaboratoryID` helper function from auth package (or make it return empty string)
- **MODIFIED** Update all tests to include `laboratory_id` query parameter in test requests

## Impact
- **Affected specs**: `auth` capability (authentication/authorization)
- **Affected code**: 
  - `backend/pkg/auth/clerk.go` - Remove laboratory_id extraction and validation
  - `backend/internal/adapters/inbound/http/handler/order.go` - Update to handle missing laboratory_id
  - `backend/internal/adapters/inbound/http/handler/client.go` - Update to handle missing laboratory_id
  - `backend/internal/adapters/inbound/http/handler/laboratory.go` - May need updates
  - All handlers that use `auth.GetLaboratoryID()` - Need alternative approach
- **Breaking changes**: 
  - Existing JWT tokens with `laboratory_id` will no longer be validated for this claim
  - **BREAKING** All API endpoints now require `laboratory_id` as a query parameter (e.g., `GET /api/v1/clients?laboratory_id=xxx`)
  - Frontend will no longer need to include `laboratory_id` in JWT custom claims
  - Frontend must include `laboratory_id` query parameter in all API requests
- **Migration considerations**: 
  - Update frontend to remove `laboratory_id` from JWT setup and add it to all API request URLs
  - All existing API calls must be updated to include `laboratory_id` query parameter
  - Test suites must be updated to include `laboratory_id` in test requests
