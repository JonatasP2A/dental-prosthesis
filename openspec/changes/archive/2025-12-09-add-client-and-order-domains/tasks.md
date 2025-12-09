## 1. Client Domain Layer
- [x] 1.1 Enhance Client domain entity with validation logic
- [x] 1.2 Add domain value objects (Address validation, reuse from Laboratory)
- [x] 1.3 Add domain errors for Client operations
- [x] 1.4 Add validation for email format, phone format, required fields

## 2. Client Ports (Interfaces)
- [x] 2.1 Create ClientRepository port in `ports/outbound/`
- [x] 2.2 Define repository interface methods (Create, Get, Update, List, Delete)
- [x] 2.3 Add FindByEmail method for duplicate email checking

## 3. Client Application Layer (Use Cases)
- [x] 3.1 Implement CreateClient use case
- [x] 3.2 Implement GetClient use case
- [x] 3.3 Implement UpdateClient use case
- [x] 3.4 Implement ListClients use case (laboratory-scoped)
- [x] 3.5 Implement DeleteClient use case (soft delete)
- [x] 3.6 Add use case error handling
- [x] 3.7 Add laboratory existence validation

## 4. Client Adapters - Inbound (HTTP)
- [x] 4.1 Create HTTP handlers package structure
- [x] 4.2 Implement POST /api/v1/clients (create)
- [x] 4.3 Implement GET /api/v1/clients/:id (get)
- [x] 4.4 Implement PUT /api/v1/clients/:id (update)
- [x] 4.5 Implement GET /api/v1/clients (list, laboratory-scoped)
- [x] 4.6 Implement DELETE /api/v1/clients/:id (soft delete)
- [x] 4.7 Add request/response DTOs
- [x] 4.8 Add input validation middleware

## 5. Client Adapters - Outbound (Persistence)
- [x] 5.1 Create repository implementation structure
- [x] 5.2 Implement in-memory repository (for initial development)
- [x] 5.3 Add repository error handling
- [x] 5.4 Add laboratory-scoped queries

## 6. Order Domain Layer
- [x] 6.1 Enhance Order domain entity with validation logic
- [x] 6.2 Add OrderStatus enum with workflow states
- [x] 6.3 Add ProsthesisItem value object
- [x] 6.4 Add domain errors for Order operations
- [x] 6.5 Add status transition validation

## 7. Order Ports (Interfaces)
- [x] 7.1 Create OrderRepository port in `ports/outbound/`
- [x] 7.2 Define repository interface methods (Create, Get, Update, List, Delete)
- [x] 7.3 Add FindByClientID method for client-scoped queries
- [x] 7.4 Add UpdateStatus method for status transitions

## 8. Order Application Layer (Use Cases)
- [x] 8.1 Implement CreateOrder use case
- [x] 8.2 Implement GetOrder use case
- [x] 8.3 Implement UpdateOrder use case
- [x] 8.4 Implement ListOrders use case (laboratory-scoped)
- [x] 8.5 Implement UpdateOrderStatus use case with workflow validation
- [x] 8.6 Implement DeleteOrder use case (soft delete)
- [x] 8.7 Add use case error handling
- [x] 8.8 Add client existence validation
- [x] 8.9 Add laboratory existence validation (via client)

## 9. Order Adapters - Inbound (HTTP)
- [x] 9.1 Create HTTP handlers package structure
- [x] 9.2 Implement POST /api/v1/orders (create)
- [x] 9.3 Implement GET /api/v1/orders/:id (get)
- [x] 9.4 Implement PUT /api/v1/orders/:id (update)
- [x] 9.5 Implement PATCH /api/v1/orders/:id/status (update status)
- [x] 9.6 Implement GET /api/v1/orders (list, laboratory-scoped)
- [x] 9.7 Implement GET /api/v1/clients/:clientId/orders (list by client)
- [x] 9.8 Implement DELETE /api/v1/orders/:id (soft delete)
- [x] 9.9 Add request/response DTOs
- [x] 9.10 Add input validation middleware

## 10. Order Adapters - Outbound (Persistence)
- [x] 10.1 Create repository implementation structure
- [x] 10.2 Implement in-memory repository (for initial development)
- [x] 10.3 Add repository error handling
- [x] 10.4 Add laboratory-scoped queries
- [x] 10.5 Add client-scoped queries

## 11. Router Integration
- [x] 11.1 Add Client routes to router
- [x] 11.2 Add Order routes to router
- [x] 11.3 Apply Clerk authentication middleware
- [x] 11.4 Add laboratory-scoped route groups

## 12. Testing
- [x] 12.1 Write unit tests for Client domain entity (90%+ coverage)
- [x] 12.2 Write unit tests for Client use cases (90%+ coverage)
- [x] 12.3 Write unit tests for Client HTTP handlers
- [x] 12.4 Write unit tests for Order domain entity (90%+ coverage)
- [x] 12.5 Write unit tests for Order use cases (90%+ coverage)
- [x] 12.6 Write unit tests for Order HTTP handlers
- [x] 12.7 Write integration tests for Client repository
- [x] 12.8 Write integration tests for Order repository
- [x] 12.9 Write integration tests for Client HTTP endpoints
- [x] 12.10 Write integration tests for Order HTTP endpoints
- [x] 12.11 Test laboratory-scoped access control
- [x] 12.12 Test order status workflow transitions

## 13. Documentation
- [x] 13.1 Update backend README with Client domain examples
- [x] 13.2 Update backend README with Order domain examples
- [x] 13.3 Add API documentation comments
- [x] 13.4 Document order workflow states
- [x] 13.5 Add HTTP client examples for Client endpoints
- [x] 13.6 Add HTTP client examples for Order endpoints
