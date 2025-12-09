package outbound

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
)

// ClientRepository defines the interface for client persistence operations
type ClientRepository interface {
	// Create stores a new client
	Create(ctx context.Context, c *client.Client) error

	// GetByID retrieves a client by ID (excludes soft-deleted)
	GetByID(ctx context.Context, id string) (*client.Client, error)

	// GetByEmail retrieves a client by email within a laboratory (excludes soft-deleted)
	GetByEmail(ctx context.Context, laboratoryID, email string) (*client.Client, error)

	// Update updates an existing client
	Update(ctx context.Context, c *client.Client) error

	// Delete performs a soft delete on a client
	Delete(ctx context.Context, id string) error

	// List retrieves all active (non-deleted) clients for a laboratory
	List(ctx context.Context, laboratoryID string) ([]*client.Client, error)
}
