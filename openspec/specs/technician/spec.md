# technician Specification

## Purpose
TBD - created by archiving change add-technician-domain. Update Purpose after archive.
## Requirements
### Requirement: Technician Entity
The system SHALL represent a laboratory technician as a domain entity with the following attributes:
- Unique identifier (ID)
- Laboratory ID (required, references Laboratory)
- Name (required, 1-200 characters)
- Email (required, valid email format)
- Phone (required, valid phone format)
- Role (required, enum: senior_technician, technician, apprentice)
- Specializations (optional, array of strings, e.g., ["crowns", "dentures", "implants"])
- Created timestamp (UTC)
- Updated timestamp (UTC)
- Deleted timestamp (UTC, nullable for soft delete)

#### Scenario: Create valid technician
- **WHEN** a technician is created with valid name, email, phone, role, and laboratory_id
- **THEN** the technician entity is created with generated ID and current UTC timestamps

#### Scenario: Reject invalid email
- **WHEN** a technician is created with invalid email format
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject empty name
- **WHEN** a technician is created with empty name
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject invalid role
- **WHEN** a technician is created with invalid or empty role
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject invalid laboratory_id
- **WHEN** a technician is created with empty or invalid laboratory_id
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Soft delete technician
- **WHEN** a technician is deleted
- **THEN** the DeletedAt timestamp is set to current UTC time
- **AND** the technician is not returned in list queries

### Requirement: Technician Role Enum
The system SHALL support the following technician roles:
- Senior Technician (senior_technician) - Experienced technician with advanced skills
- Technician (technician) - Standard technician
- Apprentice (apprentice) - Junior technician in training

#### Scenario: Accept valid technician role
- **WHEN** a technician is created with valid role from the enum
- **THEN** the technician is created successfully

#### Scenario: Reject invalid technician role
- **WHEN** a technician is created with role not in the enum
- **THEN** an ErrInvalidInput error is returned

### Requirement: Create Technician Use Case
The system SHALL allow creating a new technician with validation.

#### Scenario: Successfully create technician
- **WHEN** CreateTechnician use case is called with valid technician data and existing laboratory_id
- **THEN** the technician is persisted
- **AND** the created technician entity is returned

#### Scenario: Reject duplicate email within laboratory
- **WHEN** CreateTechnician is called with an email that already exists for the same laboratory
- **THEN** an ErrDuplicateEmail error is returned
- **AND** no technician is created

#### Scenario: Allow duplicate email across laboratories
- **WHEN** CreateTechnician is called with an email that exists for a different laboratory
- **THEN** the technician is created successfully
- **AND** email uniqueness is enforced per laboratory

#### Scenario: Validate required fields
- **WHEN** CreateTechnician is called with missing required fields
- **THEN** an ErrInvalidInput error is returned with field validation details

#### Scenario: Validate laboratory exists
- **WHEN** CreateTechnician is called with non-existent laboratory_id
- **THEN** an ErrNotFound error is returned

### Requirement: Get Technician Use Case
The system SHALL allow retrieving a technician by ID.

#### Scenario: Retrieve existing technician
- **WHEN** GetTechnician is called with valid technician ID and matching laboratory_id
- **THEN** the technician entity is returned

#### Scenario: Technician not found
- **WHEN** GetTechnician is called with non-existent technician ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Exclude soft-deleted technicians
- **WHEN** GetTechnician is called with ID of soft-deleted technician
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped access
- **WHEN** GetTechnician is called with technician ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Update Technician Use Case
The system SHALL allow updating an existing technician's information.

#### Scenario: Successfully update technician
- **WHEN** UpdateTechnician is called with valid ID, matching laboratory_id, and updated fields
- **THEN** the technician is updated
- **AND** UpdatedAt timestamp is set to current UTC time
- **AND** the updated technician entity is returned

#### Scenario: Update non-existent technician
- **WHEN** UpdateTechnician is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Validate updated email format
- **WHEN** UpdateTechnician is called with invalid email format
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Validate updated role
- **WHEN** UpdateTechnician is called with invalid role
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Laboratory-scoped update
- **WHEN** UpdateTechnician is called with technician ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: List Technicians Use Case
The system SHALL allow listing all active (non-deleted) technicians for a laboratory, optionally filtered by role.

#### Scenario: List all active technicians for laboratory
- **WHEN** ListTechnicians is called with laboratory_id
- **THEN** a list of all active technicians for that laboratory is returned
- **AND** soft-deleted technicians are excluded
- **AND** technicians from other laboratories are excluded

#### Scenario: List technicians filtered by role
- **WHEN** ListTechnicians is called with laboratory_id and role filter
- **THEN** only technicians matching the role are returned
- **AND** technicians from other laboratories are excluded

#### Scenario: Empty list
- **WHEN** ListTechnicians is called and no technicians exist for the laboratory
- **THEN** an empty list is returned

### Requirement: Delete Technician Use Case
The system SHALL allow soft-deleting a technician.

#### Scenario: Successfully delete technician
- **WHEN** DeleteTechnician is called with valid technician ID and matching laboratory_id
- **THEN** the technician's DeletedAt timestamp is set to current UTC time
- **AND** the technician is no longer returned in list queries

#### Scenario: Delete non-existent technician
- **WHEN** DeleteTechnician is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped delete
- **WHEN** DeleteTechnician is called with technician ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Technician Repository Port
The system SHALL provide a repository interface for technician persistence operations.

#### Scenario: Repository interface contract
- **WHEN** a repository implementation is created
- **THEN** it MUST implement Create, Get, Update, List, Delete, FindByEmail, and FindByRole methods
- **AND** all methods MUST match the port interface signature
- **AND** all queries MUST be scoped by laboratory_id

### Requirement: Technician HTTP API
The system SHALL provide REST endpoints for technician management.

#### Scenario: Create technician via POST
- **WHEN** POST /api/v1/technicians?laboratory_id=xxx is called with valid JSON body and authenticated user
- **THEN** HTTP 201 Created is returned
- **AND** response body contains created technician data
- **AND** laboratory_id is extracted from query parameter

#### Scenario: Get technician via GET
- **WHEN** GET /api/v1/technicians/:id?laboratory_id=xxx is called with valid ID and authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains technician data
- **AND** only technicians from user's laboratory are accessible

#### Scenario: Update technician via PUT
- **WHEN** PUT /api/v1/technicians/:id?laboratory_id=xxx is called with valid ID and JSON body
- **THEN** HTTP 200 OK is returned
- **AND** response body contains updated technician data

#### Scenario: List technicians via GET
- **WHEN** GET /api/v1/technicians?laboratory_id=xxx is called with authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains array of technicians from user's laboratory

#### Scenario: List technicians filtered by role
- **WHEN** GET /api/v1/technicians?laboratory_id=xxx&role=senior_technician is called
- **THEN** HTTP 200 OK is returned
- **AND** response body contains only technicians with role "senior_technician"

#### Scenario: Delete technician via DELETE
- **WHEN** DELETE /api/v1/technicians/:id?laboratory_id=xxx is called with valid ID
- **THEN** HTTP 204 No Content is returned
- **AND** the technician is soft-deleted

#### Scenario: Invalid JSON request body
- **WHEN** POST /api/v1/technicians?laboratory_id=xxx is called with invalid JSON
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Technician not found
- **WHEN** GET /api/v1/technicians/:id?laboratory_id=xxx is called with non-existent ID
- **THEN** HTTP 404 Not Found is returned

#### Scenario: Missing laboratory_id query parameter
- **WHEN** POST /api/v1/technicians is called without laboratory_id query parameter
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Unauthorized access
- **WHEN** GET /api/v1/technicians/:id?laboratory_id=xxx is called without authentication
- **THEN** HTTP 401 Unauthorized is returned

#### Scenario: Cross-laboratory access attempt
- **WHEN** GET /api/v1/technicians/:id?laboratory_id=xxx is called with technician ID from different laboratory
- **THEN** HTTP 404 Not Found is returned (security: don't reveal existence)

