package outbound

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
)

// ProsthesisRepository defines the interface for prosthesis persistence operations
type ProsthesisRepository interface {
	// Create stores a new prosthesis
	Create(ctx context.Context, p *prosthesis.Prosthesis) error

	// GetByID retrieves a prosthesis by ID (excludes soft-deleted)
	GetByID(ctx context.Context, id string) (*prosthesis.Prosthesis, error)

	// Update updates an existing prosthesis
	Update(ctx context.Context, p *prosthesis.Prosthesis) error

	// Delete performs a soft delete on a prosthesis
	Delete(ctx context.Context, id string) error

	// List retrieves all active (non-deleted) prostheses for a laboratory
	List(ctx context.Context, laboratoryID string) ([]*prosthesis.Prosthesis, error)

	// FindByType retrieves prostheses filtered by type for a laboratory
	FindByType(ctx context.Context, laboratoryID string, prosthesisType prosthesis.ProsthesisType) ([]*prosthesis.Prosthesis, error)

	// FindByMaterial retrieves prostheses filtered by material for a laboratory
	FindByMaterial(ctx context.Context, laboratoryID string, material string) ([]*prosthesis.Prosthesis, error)
}
