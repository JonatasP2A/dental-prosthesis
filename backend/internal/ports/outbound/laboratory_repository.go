package outbound

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

// LaboratoryRepository defines the interface for laboratory persistence operations
type LaboratoryRepository interface {
	// Create stores a new laboratory
	Create(ctx context.Context, lab *laboratory.Laboratory) error

	// GetByID retrieves a laboratory by ID (excludes soft-deleted)
	GetByID(ctx context.Context, id string) (*laboratory.Laboratory, error)

	// GetByEmail retrieves a laboratory by email (excludes soft-deleted)
	GetByEmail(ctx context.Context, email string) (*laboratory.Laboratory, error)

	// Update updates an existing laboratory
	Update(ctx context.Context, lab *laboratory.Laboratory) error

	// Delete performs a soft delete on a laboratory
	Delete(ctx context.Context, id string) error

	// List retrieves all active (non-deleted) laboratories
	List(ctx context.Context) ([]*laboratory.Laboratory, error)
}

