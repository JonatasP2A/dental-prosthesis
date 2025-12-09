# Change: Add Client and Order Domains Implementation

## Why
The Client and Order domains are core business entities that enable the dental prosthesis SaaS platform's primary workflow. Clients (dental clinics/dentists) place Orders for prosthetic work, which laboratories then fulfill. Implementing these domains establishes the relationship chain: Laboratory → Client → Order, enabling the platform to track work orders from receipt through delivery. These domains build upon the Laboratory domain's multi-tenancy foundation and follow the same hexagonal architecture patterns.

## What Changes
- **ADDED** Client domain implementation following hexagonal architecture
  - Domain entity with business logic and validation
  - Repository port (interface) for persistence
  - Use cases: CreateClient, GetClient, UpdateClient, ListClients, DeleteClient
  - HTTP adapter with REST endpoints (`/api/v1/clients`)
  - Laboratory-scoped access control (clients belong to a laboratory)
  - Unit tests with 90%+ coverage requirement

- **ADDED** Order domain implementation following hexagonal architecture
  - Domain entity with workflow state management
  - Repository port (interface) for persistence
  - Use cases: CreateOrder, GetOrder, UpdateOrder, ListOrders, UpdateOrderStatus
  - HTTP adapter with REST endpoints (`/api/v1/orders`)
  - Client relationship validation (orders must belong to valid client)
  - Laboratory-scoped access control (orders belong to a laboratory via client)
  - Order status workflow enforcement (Received → In Production → Quality Check → Ready → Delivered/Revision)
  - Unit tests with 90%+ coverage requirement

## Impact
- **Affected specs**: New capabilities `client` and `order` added
- **Affected code**: 
  - `backend/internal/domain/client/` - Client domain entity
  - `backend/internal/domain/order/` - Order domain entity
  - `backend/internal/ports/outbound/` - Repository interfaces
  - `backend/internal/application/client/` - Client use cases
  - `backend/internal/application/order/` - Order use cases
  - `backend/internal/adapters/inbound/http/` - HTTP handlers for both domains
  - `backend/internal/adapters/outbound/persistence/` - Repository implementations
  - `backend/test/` - Integration tests
- **Dependencies**: 
  - Laboratory domain (must exist)
  - Client domain (required for Order domain)
- **Relationships**:
  - Client belongs to Laboratory (via laboratory_id)
  - Order belongs to Client (via client_id) and Laboratory (via laboratory_id from client)
