package outbound

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
)

// OrderRepository defines the interface for order persistence operations
type OrderRepository interface {
	// Create stores a new order
	Create(ctx context.Context, o *order.Order) error

	// GetByID retrieves an order by ID (excludes soft-deleted)
	GetByID(ctx context.Context, id string) (*order.Order, error)

	// Update updates an existing order
	Update(ctx context.Context, o *order.Order) error

	// UpdateStatus updates only the order status
	UpdateStatus(ctx context.Context, id string, status order.Status) error

	// Delete performs a soft delete on an order
	Delete(ctx context.Context, id string) error

	// List retrieves all active (non-deleted) orders for a laboratory
	List(ctx context.Context, laboratoryID string) ([]*order.Order, error)

	// ListByClientID retrieves all active orders for a specific client
	ListByClientID(ctx context.Context, clientID string) ([]*order.Order, error)
}
