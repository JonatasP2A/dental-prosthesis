# Client Domain Specification

## ADDED Requirements

### Requirement: Client Entity
The system SHALL represent a dental clinic or dentist (client) as a domain entity with the following attributes:
- Unique identifier (ID)
- Laboratory ID (required, references Laboratory)
- Name (required, 1-200 characters)
- Email (required, valid email format)
- Phone (required, valid phone format)
- Address (street, city, state, postal code, country)
- Created timestamp (UTC)
- Updated timestamp (UTC)
- Deleted timestamp (UTC, nullable for soft delete)

#### Scenario: Create valid client
- **WHEN** a client is created with valid name, email, phone, address, and laboratory_id
- **THEN** the client entity is created with generated ID and current UTC timestamps

#### Scenario: Reject invalid email
- **WHEN** a client is created with invalid email format
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject empty name
- **WHEN** a client is created with empty name
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject invalid laboratory_id
- **WHEN** a client is created with empty or invalid laboratory_id
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Soft delete client
- **WHEN** a client is deleted
- **THEN** the DeletedAt timestamp is set to current UTC time
- **AND** the client is not returned in list queries

### Requirement: Create Client Use Case
The system SHALL allow creating a new client with validation.

#### Scenario: Successfully create client
- **WHEN** CreateClient use case is called with valid client data and existing laboratory_id
- **THEN** the client is persisted
- **AND** the created client entity is returned

#### Scenario: Reject duplicate email within laboratory
- **WHEN** CreateClient is called with an email that already exists for the same laboratory
- **THEN** an ErrDuplicateEmail error is returned
- **AND** no client is created

#### Scenario: Allow duplicate email across laboratories
- **WHEN** CreateClient is called with an email that exists for a different laboratory
- **THEN** the client is created successfully
- **AND** email uniqueness is enforced per laboratory

#### Scenario: Validate required fields
- **WHEN** CreateClient is called with missing required fields
- **THEN** an ErrInvalidInput error is returned with field validation details

#### Scenario: Validate laboratory exists
- **WHEN** CreateClient is called with non-existent laboratory_id
- **THEN** an ErrNotFound error is returned

### Requirement: Get Client Use Case
The system SHALL allow retrieving a client by ID.

#### Scenario: Retrieve existing client
- **WHEN** GetClient is called with valid client ID and matching laboratory_id
- **THEN** the client entity is returned

#### Scenario: Client not found
- **WHEN** GetClient is called with non-existent client ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Exclude soft-deleted clients
- **WHEN** GetClient is called with ID of soft-deleted client
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped access
- **WHEN** GetClient is called with client ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Update Client Use Case
The system SHALL allow updating an existing client's information.

#### Scenario: Successfully update client
- **WHEN** UpdateClient is called with valid ID, matching laboratory_id, and updated fields
- **THEN** the client is updated
- **AND** UpdatedAt timestamp is set to current UTC time
- **AND** the updated client entity is returned

#### Scenario: Update non-existent client
- **WHEN** UpdateClient is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Validate updated email format
- **WHEN** UpdateClient is called with invalid email format
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Laboratory-scoped update
- **WHEN** UpdateClient is called with client ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: List Clients Use Case
The system SHALL allow listing all active (non-deleted) clients for a laboratory.

#### Scenario: List all active clients for laboratory
- **WHEN** ListClients is called with laboratory_id
- **THEN** a list of all active clients for that laboratory is returned
- **AND** soft-deleted clients are excluded
- **AND** clients from other laboratories are excluded

#### Scenario: Empty list
- **WHEN** ListClients is called and no clients exist for the laboratory
- **THEN** an empty list is returned

### Requirement: Delete Client Use Case
The system SHALL allow soft-deleting a client.

#### Scenario: Successfully delete client
- **WHEN** DeleteClient is called with valid client ID and matching laboratory_id
- **THEN** the client's DeletedAt timestamp is set to current UTC time
- **AND** the client is no longer returned in list queries

#### Scenario: Delete non-existent client
- **WHEN** DeleteClient is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped delete
- **WHEN** DeleteClient is called with client ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Client Repository Port
The system SHALL provide a repository interface for client persistence operations.

#### Scenario: Repository interface contract
- **WHEN** a repository implementation is created
- **THEN** it MUST implement Create, Get, Update, List, Delete, and FindByEmail methods
- **AND** all methods MUST match the port interface signature
- **AND** all queries MUST be scoped by laboratory_id

### Requirement: Client HTTP API
The system SHALL provide REST endpoints for client management.

#### Scenario: Create client via POST
- **WHEN** POST /api/v1/clients is called with valid JSON body and authenticated user
- **THEN** HTTP 201 Created is returned
- **AND** response body contains created client data
- **AND** laboratory_id is automatically set from JWT claim

#### Scenario: Get client via GET
- **WHEN** GET /api/v1/clients/:id is called with valid ID and authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains client data
- **AND** only clients from user's laboratory are accessible

#### Scenario: Update client via PUT
- **WHEN** PUT /api/v1/clients/:id is called with valid ID and JSON body
- **THEN** HTTP 200 OK is returned
- **AND** response body contains updated client data

#### Scenario: List clients via GET
- **WHEN** GET /api/v1/clients is called with authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains array of clients from user's laboratory

#### Scenario: Delete client via DELETE
- **WHEN** DELETE /api/v1/clients/:id is called with valid ID
- **THEN** HTTP 204 No Content is returned
- **AND** the client is soft-deleted

#### Scenario: Invalid JSON request body
- **WHEN** POST /api/v1/clients is called with invalid JSON
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Client not found
- **WHEN** GET /api/v1/clients/:id is called with non-existent ID
- **THEN** HTTP 404 Not Found is returned

#### Scenario: Unauthorized access
- **WHEN** GET /api/v1/clients/:id is called without authentication
- **THEN** HTTP 401 Unauthorized is returned

#### Scenario: Cross-laboratory access attempt
- **WHEN** GET /api/v1/clients/:id is called with client ID from different laboratory
- **THEN** HTTP 404 Not Found is returned (security: don't reveal existence)
