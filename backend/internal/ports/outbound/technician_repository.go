package outbound

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
)

// TechnicianRepository defines the interface for technician persistence operations
type TechnicianRepository interface {
	// Create stores a new technician
	Create(ctx context.Context, tech *technician.Technician) error

	// GetByID retrieves a technician by ID (excludes soft-deleted)
	GetByID(ctx context.Context, id string) (*technician.Technician, error)

	// GetByEmail retrieves a technician by email within a laboratory (excludes soft-deleted)
	GetByEmail(ctx context.Context, laboratoryID, email string) (*technician.Technician, error)

	// Update updates an existing technician
	Update(ctx context.Context, tech *technician.Technician) error

	// Delete performs a soft delete on a technician
	Delete(ctx context.Context, id string) error

	// List retrieves all active (non-deleted) technicians for a laboratory
	List(ctx context.Context, laboratoryID string) ([]*technician.Technician, error)

	// ListByRole retrieves all active technicians for a laboratory filtered by role
	ListByRole(ctx context.Context, laboratoryID string, role technician.Role) ([]*technician.Technician, error)
}
