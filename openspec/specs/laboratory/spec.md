# laboratory Specification

## Purpose
TBD - created by archiving change add-laboratory-domain. Update Purpose after archive.
## Requirements
### Requirement: Laboratory Entity
The system SHALL represent a dental prosthesis laboratory as a domain entity with the following attributes:
- Unique identifier (ID)
- Name (required, 1-200 characters)
- Email (required, valid email format)
- Phone (required, valid phone format)
- Address (street, city, state, postal code, country)
- Created timestamp (UTC)
- Updated timestamp (UTC)
- Deleted timestamp (UTC, nullable for soft delete)

#### Scenario: Create valid laboratory
- **WHEN** a laboratory is created with valid name, email, phone, and address
- **THEN** the laboratory entity is created with generated ID and current UTC timestamps

#### Scenario: Reject invalid email
- **WHEN** a laboratory is created with invalid email format
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject empty name
- **WHEN** a laboratory is created with empty name
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Soft delete laboratory
- **WHEN** a laboratory is deleted
- **THEN** the DeletedAt timestamp is set to current UTC time
- **AND** the laboratory is not returned in list queries

### Requirement: Create Laboratory Use Case
The system SHALL allow creating a new laboratory with validation.

#### Scenario: Successfully create laboratory
- **WHEN** CreateLaboratory use case is called with valid laboratory data
- **THEN** the laboratory is persisted
- **AND** the created laboratory entity is returned

#### Scenario: Reject duplicate email
- **WHEN** CreateLaboratory is called with an email that already exists
- **THEN** an ErrDuplicateEmail error is returned
- **AND** no laboratory is created

#### Scenario: Validate required fields
- **WHEN** CreateLaboratory is called with missing required fields
- **THEN** an ErrInvalidInput error is returned with field validation details

### Requirement: Get Laboratory Use Case
The system SHALL allow retrieving a laboratory by ID.

#### Scenario: Retrieve existing laboratory
- **WHEN** GetLaboratory is called with valid laboratory ID
- **THEN** the laboratory entity is returned

#### Scenario: Laboratory not found
- **WHEN** GetLaboratory is called with non-existent laboratory ID
- **THEN** an ErrLaboratoryNotFound error is returned

#### Scenario: Exclude soft-deleted laboratories
- **WHEN** GetLaboratory is called with ID of soft-deleted laboratory
- **THEN** an ErrLaboratoryNotFound error is returned

### Requirement: Update Laboratory Use Case
The system SHALL allow updating an existing laboratory's information.

#### Scenario: Successfully update laboratory
- **WHEN** UpdateLaboratory is called with valid ID and updated fields
- **THEN** the laboratory is updated
- **AND** UpdatedAt timestamp is set to current UTC time
- **AND** the updated laboratory entity is returned

#### Scenario: Update non-existent laboratory
- **WHEN** UpdateLaboratory is called with non-existent ID
- **THEN** an ErrLaboratoryNotFound error is returned

#### Scenario: Validate updated email format
- **WHEN** UpdateLaboratory is called with invalid email format
- **THEN** an ErrInvalidInput error is returned

### Requirement: List Laboratories Use Case
The system SHALL allow listing all active (non-deleted) laboratories.

#### Scenario: List all active laboratories
- **WHEN** ListLaboratories is called
- **THEN** a list of all active laboratories is returned
- **AND** soft-deleted laboratories are excluded

#### Scenario: Empty list
- **WHEN** ListLaboratories is called and no laboratories exist
- **THEN** an empty list is returned

### Requirement: Laboratory Repository Port
The system SHALL provide a repository interface for laboratory persistence operations.

#### Scenario: Repository interface contract
- **WHEN** a repository implementation is created
- **THEN** it MUST implement Create, Get, Update, List, and Delete methods
- **AND** all methods MUST match the port interface signature

### Requirement: Laboratory HTTP API
The system SHALL provide REST endpoints for laboratory management.

#### Scenario: Create laboratory via POST
- **WHEN** POST /api/v1/laboratories is called with valid JSON body
- **THEN** HTTP 201 Created is returned
- **AND** response body contains created laboratory data

#### Scenario: Get laboratory via GET
- **WHEN** GET /api/v1/laboratories/:id is called with valid ID
- **THEN** HTTP 200 OK is returned
- **AND** response body contains laboratory data

#### Scenario: Update laboratory via PUT
- **WHEN** PUT /api/v1/laboratories/:id is called with valid ID and JSON body
- **THEN** HTTP 200 OK is returned
- **AND** response body contains updated laboratory data

#### Scenario: List laboratories via GET
- **WHEN** GET /api/v1/laboratories is called
- **THEN** HTTP 200 OK is returned
- **AND** response body contains array of laboratories

#### Scenario: Invalid JSON request body
- **WHEN** POST /api/v1/laboratories is called with invalid JSON
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Laboratory not found
- **WHEN** GET /api/v1/laboratories/:id is called with non-existent ID
- **THEN** HTTP 404 Not Found is returned

### Requirement: Clerk Authentication Middleware
The system SHALL authenticate requests using Clerk JWT tokens and extract laboratory context.

#### Scenario: Valid JWT token
- **WHEN** a request includes valid Clerk JWT token with laboratory_id claim
- **THEN** the request is authenticated
- **AND** laboratory_id is extracted and available in request context

#### Scenario: Missing JWT token
- **WHEN** a request does not include Authorization header
- **THEN** HTTP 401 Unauthorized is returned

#### Scenario: Invalid JWT token
- **WHEN** a request includes invalid or expired JWT token
- **THEN** HTTP 401 Unauthorized is returned

#### Scenario: Missing laboratory_id claim
- **WHEN** a request includes valid JWT but without laboratory_id claim
- **THEN** HTTP 403 Forbidden is returned

### Requirement: Multi-Tenancy Isolation
The system SHALL ensure laboratory data is isolated by laboratory_id from authenticated context.

#### Scenario: Laboratory-scoped queries
- **WHEN** GetLaboratory or ListLaboratories is called
- **THEN** only laboratories matching the authenticated user's laboratory_id are returned
- **AND** users cannot access laboratories they don't belong to

#### Scenario: Cross-laboratory access attempt
- **WHEN** a user attempts to access a laboratory with different laboratory_id than their claim
- **THEN** HTTP 403 Forbidden is returned

