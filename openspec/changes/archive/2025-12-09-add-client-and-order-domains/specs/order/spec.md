# Order Domain Specification

## ADDED Requirements

### Requirement: Order Entity
The system SHALL represent a prosthesis work order as a domain entity with the following attributes:
- Unique identifier (ID)
- Client ID (required, references Client)
- Laboratory ID (required, derived from Client, references Laboratory)
- Status (required, enum: received, in_production, quality_check, ready, delivered, revision)
- Prosthesis items (array of ProsthesisItem)
- Created timestamp (UTC)
- Updated timestamp (UTC)
- Deleted timestamp (UTC, nullable for soft delete)

#### Scenario: Create valid order
- **WHEN** an order is created with valid client_id, status, and at least one prosthesis item
- **THEN** the order entity is created with generated ID and current UTC timestamps
- **AND** laboratory_id is derived from the client

#### Scenario: Reject order without prosthesis items
- **WHEN** an order is created with empty prosthesis items array
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Reject invalid client_id
- **WHEN** an order is created with empty or invalid client_id
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Soft delete order
- **WHEN** an order is deleted
- **THEN** the DeletedAt timestamp is set to current UTC time
- **AND** the order is not returned in list queries

### Requirement: ProsthesisItem Value Object
The system SHALL represent a prosthesis item within an order with the following attributes:
- Type (required, e.g., crown, bridge, denture, implant, veneer)
- Material (required, e.g., zirconia, porcelain, metal alloy)
- Shade (optional, tooth color matching, e.g., A1, A2, B1 using VITA scale)
- Quantity (required, integer > 0)
- Notes (optional, text field for special instructions)

#### Scenario: Reject invalid prosthesis item quantity
- **WHEN** an order is created with prosthesis item having quantity <= 0
- **THEN** an ErrInvalidInput error is returned

#### Scenario: Accept valid prosthesis item
- **WHEN** an order is created with prosthesis item having valid type, material, and quantity > 0
- **THEN** the prosthesis item is included in the order

### Requirement: Order Status Workflow
The system SHALL enforce valid status transitions following the workflow: Received → In Production → Quality Check → Ready → Delivered/Revision.

#### Scenario: Valid status transition
- **WHEN** UpdateOrderStatus is called with a valid next status in the workflow
- **THEN** the order status is updated
- **AND** UpdatedAt timestamp is set to current UTC time

#### Scenario: Invalid status transition
- **WHEN** UpdateOrderStatus is called with an invalid status transition (e.g., Received → Delivered)
- **THEN** an ErrInvalidStatusTransition error is returned
- **AND** the order status remains unchanged

#### Scenario: Revision can return to In Production
- **WHEN** UpdateOrderStatus is called to change Revision → In Production
- **THEN** the status transition is allowed
- **AND** the order status is updated

#### Scenario: Initial status is Received
- **WHEN** CreateOrder is called
- **THEN** the order is created with status "received" by default

### Requirement: Create Order Use Case
The system SHALL allow creating a new order with validation.

#### Scenario: Successfully create order
- **WHEN** CreateOrder use case is called with valid order data and existing client_id
- **THEN** the order is persisted with status "received"
- **AND** the created order entity is returned
- **AND** laboratory_id is derived from the client

#### Scenario: Validate required fields
- **WHEN** CreateOrder is called with missing required fields
- **THEN** an ErrInvalidInput error is returned with field validation details

#### Scenario: Validate client exists
- **WHEN** CreateOrder is called with non-existent client_id
- **THEN** an ErrNotFound error is returned

#### Scenario: Validate client belongs to laboratory
- **WHEN** CreateOrder is called with client_id from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Get Order Use Case
The system SHALL allow retrieving an order by ID.

#### Scenario: Retrieve existing order
- **WHEN** GetOrder is called with valid order ID and matching laboratory_id
- **THEN** the order entity is returned

#### Scenario: Order not found
- **WHEN** GetOrder is called with non-existent order ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Exclude soft-deleted orders
- **WHEN** GetOrder is called with ID of soft-deleted order
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped access
- **WHEN** GetOrder is called with order ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Update Order Use Case
The system SHALL allow updating an existing order's information (excluding status).

#### Scenario: Successfully update order
- **WHEN** UpdateOrder is called with valid ID, matching laboratory_id, and updated fields
- **THEN** the order is updated
- **AND** UpdatedAt timestamp is set to current UTC time
- **AND** the updated order entity is returned

#### Scenario: Update non-existent order
- **WHEN** UpdateOrder is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped update
- **WHEN** UpdateOrder is called with order ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Update Order Status Use Case
The system SHALL allow updating an order's status with workflow validation.

#### Scenario: Successfully update status
- **WHEN** UpdateOrderStatus is called with valid order ID and valid next status
- **THEN** the order status is updated
- **AND** UpdatedAt timestamp is set to current UTC time
- **AND** the updated order entity is returned

#### Scenario: Invalid status transition
- **WHEN** UpdateOrderStatus is called with invalid status transition
- **THEN** an ErrInvalidStatusTransition error is returned
- **AND** the order status remains unchanged

#### Scenario: Update status of non-existent order
- **WHEN** UpdateOrderStatus is called with non-existent order ID
- **THEN** an ErrNotFound error is returned

### Requirement: List Orders Use Case
The system SHALL allow listing orders for a laboratory, optionally filtered by client.

#### Scenario: List all active orders for laboratory
- **WHEN** ListOrders is called with laboratory_id
- **THEN** a list of all active orders for that laboratory is returned
- **AND** soft-deleted orders are excluded
- **AND** orders from other laboratories are excluded

#### Scenario: List orders by client
- **WHEN** ListOrders is called with laboratory_id and client_id
- **THEN** a list of all active orders for that client is returned
- **AND** only orders from the specified laboratory are included

#### Scenario: Empty list
- **WHEN** ListOrders is called and no orders exist for the laboratory
- **THEN** an empty list is returned

### Requirement: Delete Order Use Case
The system SHALL allow soft-deleting an order.

#### Scenario: Successfully delete order
- **WHEN** DeleteOrder is called with valid order ID and matching laboratory_id
- **THEN** the order's DeletedAt timestamp is set to current UTC time
- **AND** the order is no longer returned in list queries

#### Scenario: Delete non-existent order
- **WHEN** DeleteOrder is called with non-existent ID
- **THEN** an ErrNotFound error is returned

#### Scenario: Laboratory-scoped delete
- **WHEN** DeleteOrder is called with order ID from different laboratory
- **THEN** an ErrNotFound error is returned (access denied)

### Requirement: Order Repository Port
The system SHALL provide a repository interface for order persistence operations.

#### Scenario: Repository interface contract
- **WHEN** a repository implementation is created
- **THEN** it MUST implement Create, Get, Update, List, Delete, FindByClientID, and UpdateStatus methods
- **AND** all methods MUST match the port interface signature
- **AND** all queries MUST be scoped by laboratory_id

### Requirement: Order HTTP API
The system SHALL provide REST endpoints for order management.

#### Scenario: Create order via POST
- **WHEN** POST /api/v1/orders is called with valid JSON body and authenticated user
- **THEN** HTTP 201 Created is returned
- **AND** response body contains created order data
- **AND** laboratory_id is automatically derived from client

#### Scenario: Get order via GET
- **WHEN** GET /api/v1/orders/:id is called with valid ID and authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains order data
- **AND** only orders from user's laboratory are accessible

#### Scenario: Update order via PUT
- **WHEN** PUT /api/v1/orders/:id is called with valid ID and JSON body
- **THEN** HTTP 200 OK is returned
- **AND** response body contains updated order data

#### Scenario: Update order status via PATCH
- **WHEN** PATCH /api/v1/orders/:id/status is called with valid ID and status in body
- **THEN** HTTP 200 OK is returned
- **AND** response body contains updated order data with new status

#### Scenario: List orders via GET
- **WHEN** GET /api/v1/orders is called with authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains array of orders from user's laboratory

#### Scenario: List orders by client via GET
- **WHEN** GET /api/v1/clients/:clientId/orders is called with authenticated user
- **THEN** HTTP 200 OK is returned
- **AND** response body contains array of orders for that client
- **AND** only orders from user's laboratory are included

#### Scenario: Delete order via DELETE
- **WHEN** DELETE /api/v1/orders/:id is called with valid ID
- **THEN** HTTP 204 No Content is returned
- **AND** the order is soft-deleted

#### Scenario: Invalid JSON request body
- **WHEN** POST /api/v1/orders is called with invalid JSON
- **THEN** HTTP 400 Bad Request is returned

#### Scenario: Order not found
- **WHEN** GET /api/v1/orders/:id is called with non-existent ID
- **THEN** HTTP 404 Not Found is returned

#### Scenario: Invalid status transition
- **WHEN** PATCH /api/v1/orders/:id/status is called with invalid status transition
- **THEN** HTTP 400 Bad Request is returned
- **AND** error message indicates invalid transition

#### Scenario: Unauthorized access
- **WHEN** GET /api/v1/orders/:id is called without authentication
- **THEN** HTTP 401 Unauthorized is returned

#### Scenario: Cross-laboratory access attempt
- **WHEN** GET /api/v1/orders/:id is called with order ID from different laboratory
- **THEN** HTTP 404 Not Found is returned (security: don't reveal existence)
