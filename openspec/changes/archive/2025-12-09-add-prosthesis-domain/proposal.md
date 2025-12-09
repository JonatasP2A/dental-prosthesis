# Change: Add Prosthesis Domain Implementation

## Why
The Prosthesis domain represents individual dental prosthetic items (crowns, bridges, dentures, implants, etc.) that can be tracked independently from orders. While orders contain prosthesis items as value objects, a full Prosthesis domain enables laboratories to manage prosthesis templates, specifications, materials catalog, and track individual prosthesis items through their lifecycle. This domain supports advanced features like prosthesis inventory management, material cost tracking, and prosthesis-specific workflows.

## What Changes
- **ADDED** Prosthesis domain implementation following hexagonal architecture
  - Domain entity with business logic and validation
  - Repository port (interface) for persistence
  - Use cases: CreateProsthesis, GetProsthesis, UpdateProsthesis, ListProstheses, DeleteProsthesis
  - HTTP adapter with REST endpoints (`/api/v1/prostheses`)
  - Laboratory-scoped access control (prostheses belong to a laboratory)
  - Prosthesis type validation (crown, bridge, denture, implant, veneer, inlay/onlay)
  - Material validation and specifications
  - Unit tests with 90%+ coverage requirement

## Impact
- **Affected specs**: New capability `prosthesis` added
- **Affected code**: 
  - `backend/internal/domain/prosthesis/` - Prosthesis domain entity
  - `backend/internal/ports/outbound/` - Repository interface
  - `backend/internal/application/prosthesis/` - Prosthesis use cases
  - `backend/internal/adapters/inbound/http/` - HTTP handlers
  - `backend/internal/adapters/outbound/persistence/` - Repository implementation
  - `backend/test/` - Integration tests
- **Dependencies**: 
  - Laboratory domain (must exist)
- **Relationships**:
  - Prosthesis belongs to Laboratory (via laboratory_id)
  - Orders reference Prosthesis items (via prosthesis_id in future enhancements)
