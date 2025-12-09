# prosthesis Specification

## Purpose
TBD - created by archiving change add-prosthesis-domain. Update Purpose after archive.
## Requirements
### Requirement: Prosthesis Entity
The system SHALL represent a dental prosthetic item as a domain entity with the following attributes:
- Unique identifier (ID)
- Laboratory ID (required, references Laboratory)
- Type (required, enum: crown, bridge, complete_denture, partial_denture, implant, veneer, inlay, onlay)
- Material (required, e.g., zirconia, porcelain, metal alloy, acrylic)
- Shade (optional, tooth color matching, e.g., A1, A2, B1 using VITA scale)
- Specifications (optional, text field for technical details)
- Notes (optional, text field for special instructions)
- Created timestamp (UTC)
- Updated timestamp (UTC)
- Deleted timestamp (UTC, nullable for soft delete)

#### Scenario: Create valid prosthesis
- **WHEN** a prosthesis is created with valid type, material, laboratory_id, and optional shade/specifications
- **THEN** the prosthesis entity is created with generated ID and current UTC timestamps

#### Scenario: Reject invalid prosthesis type
- **WHEN** a prosthesis is created with invalid or empty type
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject empty material
- **WHEN** a prosthesis is created with empty material
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject invalid laboratory_id
- **WHEN** a prosthesis is created with empty or invalid laboratory_id
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Soft delete prosthesis
- **WHEN** a prosthesis is deleted
- **THEN** the DeletedAt timestamp is set to current UTC time
- **AND** the prosthesis is not returned in list queries

### Requirement: Prosthesis Type Enum
The system SHALL support the following prosthesis types:
- Crown (coroa)
- Bridge (ponte)
- Complete denture (prótese total)
- Partial denture (prótese parcial)
- Implant-supported prosthetic
- Veneer (faceta)
- Inlay
- Onlay

#### Scenario: Accept valid prosthesis type
- **WHEN** a prosthesis is created with valid type from the enum
- **THEN** the prosthesis is created successfully

#### Scenario: Reject invalid prosthesis type
- **WHEN** a prosthesis is created with type not in the enum
- **THEN** an ErrInvalidInput error is returned

### Requirement: Create Prosthesis Use Case
The system SHALL allow creating a new prosthesis with validation.

#### Scenario: Successfully create prosthesis
- **WHEN** CreateProsthesis use case is called with valid prosthesis data and existing laboratory_id
- **THEN** the prosthesis is persisted
- **AND** the created prosthesis entity is returned

#### Scenario: Validate required fields
- **WHEN** CreateProsthesis is called with missing required fields
- **THEN** an ErrInvalidInput error is returned with field validation details

#### Scenario: Validate laboratory exists
- **WHEN** CreateProsthesis is called with non-existent laboratory_id
- **THEN** an ErrNotFound error is returned

### Requirement: Get Prosthesis Use Case
The system SHALL allow retrieving a prosthesis by ID.

#### Scenario: Retrieve existing prosthesis
- **WHEN** GetProsthesis is called with valid prosthesis ID and matching laboratory_id
- **THEN** the prosthesis entity is returned

#### Scenario: Prosthesis not found
- **WHEN** GetProsthesis is called with non-existent prosthesis ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Exclude soft-deleted prostheses
- **WHEN** GetProsthesis is called with ID of soft-deleted prosthesis
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped access
- **WHEN** GetProsthesis is called with prosthesis ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Update Prosthesis Use Case
The system SHALL allow updating an existing prosthesis's information.

#### Scenario: Successfully update prosthesis
- **WHEN** UpdateProsthesis is called with valid ID, matching laboratory_id, and updated fields
- **THEN** the prosthesis is updated
- **AND** UpdatedAt timestamp is set to current UTC time
- **AND** the updated prosthesis entity is returned

#### Scenario: Update non-existent prosthesis
- **WHEN** UpdateProsthesis is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Validate updated type
- **WHEN** UpdateProsthesis is called with invalid prosthesis type
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Laboratory-scoped update
- **WHEN** UpdateProsthesis is called with prosthesis ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: List Prostheses Use Case
The system SHALL allow listing all active (non-deleted) prostheses for a laboratory, optionally filtered by type or material.

#### Scenario: List all active prostheses for laboratory
- **WHEN** ListProstheses is called with laboratory_id
- **THEN** a list of all active prostheses for that laboratory is returned
- **AND** soft-deleted prostheses are excluded
- **AND** prostheses from other laboratories are excluded

#### Scenario: List prostheses filtered by type
- **WHEN** ListProstheses is called with laboratory_id and type filter
- **THEN** only prostheses matching the type are returned
- **AND** prostheses from other laboratories are excluded

#### Scenario: List prostheses filtered by material
- **WHEN** ListProstheses is called with laboratory_id and material filter
- **THEN** only prostheses matching the material are returned
- **AND** prostheses from other laboratories are excluded

#### Scenario: Empty list
- **WHEN** ListProstheses is called and no prostheses exist for the laboratory
- **THEN** an empty list is returned

### Requirement: Delete Prosthesis Use Case
The system SHALL allow soft-deleting a prosthesis.

#### Scenario: Successfully delete prosthesis
- **WHEN** DeleteProsthesis is called with valid prosthesis ID and matching laboratory_id
- **THEN** the prosthesis's DeletedAt timestamp is set to current UTC time
- **AND** the prosthesis is no longer returned in list queries

#### Scenario: Delete non-existent prosthesis
- **WHEN** DeleteProsthesis is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped delete
- **WHEN** DeleteProsthesis is called with prosthesis ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Prosthesis Repository Port
The system SHALL provide a repository interface for prosthesis persistence operations.

#### Scenario: Repository interface contract
- **WHEN** a repository implementation is created
- **THEN** it MUST implement Create, Get, Update, List, Delete, FindByType, and FindByMaterial methods
- **AND** all methods MUST match the port interface signature
- **AND** all queries MUST be scoped by laboratory_id

### Requirement: Prosthesis HTTP API
The system SHALL provide REST endpoints for prosthesis management.

#### Scenario: Create prosthesis via POST
- **WHEN** POST /api/v1/prostheses?laboratory_id=xxx is called with valid JSON body and authenticated user
- **THEN** HTTP 201 Created is returned
- **AND** response body contains created prosthesis data
- **AND** laboratory_id is extracted from query parameter

#### Scenario: Get prosthesis via GET
- **WHEN** GET /api/v1/prostheses/:id?laboratory_id=xxx is called with valid ID and authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains prosthesis data
- **AND** only prostheses from user's laboratory are accessible

#### Scenario: Update prosthesis via PUT
- **WHEN** PUT /api/v1/prostheses/:id?laboratory_id=xxx is called with valid ID and JSON body
- **THEN** HTTP 200 OK is returned
- **AND** response body contains updated prosthesis data

#### Scenario: List prostheses via GET
- **WHEN** GET /api/v1/prostheses?laboratory_id=xxx is called with authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains array of prostheses from user's laboratory

#### Scenario: List prostheses filtered by type
- **WHEN** GET /api/v1/prostheses?laboratory_id=xxx&type=crown is called
- **THEN** HTTP 200 OK is returned
- **AND** response body contains only prostheses of type "crown"

#### Scenario: Delete prosthesis via DELETE
- **WHEN** DELETE /api/v1/prostheses/:id?laboratory_id=xxx is called with valid ID
- **THEN** HTTP 204 No Content is returned
- **AND** the prosthesis is soft-deleted

#### Scenario: Invalid JSON request body
- **WHEN** POST /api/v1/prostheses?laboratory_id=xxx is called with invalid JSON
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Prosthesis not found
- **WHEN** GET /api/v1/prostheses/:id?laboratory_id=xxx is called with non-existent ID
- **THEN** HTTP 404 Not Found is returned

#### Scenario: Missing laboratory_id query parameter
- **WHEN** POST /api/v1/prostheses is called without laboratory_id query parameter
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Unauthorized access
- **WHEN** GET /api/v1/prostheses/:id?laboratory_id=xxx is called without authentication
- **THEN** HTTP 401 Unauthorized is returned

#### Scenario: Cross-laboratory access attempt
- **WHEN** GET /api/v1/prostheses/:id?laboratory_id=xxx is called with prosthesis ID from different laboratory
- **THEN** HTTP 404 Not Found is returned (security: don't reveal existence)

