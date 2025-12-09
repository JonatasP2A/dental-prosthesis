package memory

import (
	"context"
	"sync"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
)

// OrderRepository is an in-memory implementation of the order repository
type OrderRepository struct {
	mu   sync.RWMutex
	data map[string]*order.Order
}

// NewOrderRepository creates a new in-memory order repository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		data: make(map[string]*order.Order),
	}
}

// Create stores a new order
func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[o.ID]; exists {
		return errors.ErrInternal // ID already exists
	}

	// Clone to avoid external modifications
	r.data[o.ID] = r.clone(o)
	return nil
}

// GetByID retrieves an order by ID (excludes soft-deleted)
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, exists := r.data[id]
	if !exists || o.IsDeleted() {
		return nil, errors.ErrNotFound
	}

	return r.clone(o), nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[o.ID]
	if !exists || existing.IsDeleted() {
		return errors.ErrNotFound
	}

	r.data[o.ID] = r.clone(o)
	return nil
}

// UpdateStatus updates only the order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status order.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, exists := r.data[id]
	if !exists || o.IsDeleted() {
		return errors.ErrNotFound
	}

	o.Status = status
	return nil
}

// Delete performs a soft delete on an order
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, exists := r.data[id]
	if !exists || o.IsDeleted() {
		return errors.ErrNotFound
	}

	o.Delete()
	return nil
}

// List retrieves all active (non-deleted) orders for a laboratory
func (r *OrderRepository) List(ctx context.Context, laboratoryID string) ([]*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var orders []*order.Order
	for _, o := range r.data {
		if o.LaboratoryID == laboratoryID && !o.IsDeleted() {
			orders = append(orders, r.clone(o))
		}
	}

	return orders, nil
}

// ListByClientID retrieves all active orders for a specific client
func (r *OrderRepository) ListByClientID(ctx context.Context, clientID string) ([]*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var orders []*order.Order
	for _, o := range r.data {
		if o.ClientID == clientID && !o.IsDeleted() {
			orders = append(orders, r.clone(o))
		}
	}

	return orders, nil
}

// clone creates a deep copy of an order to avoid external modifications
func (r *OrderRepository) clone(o *order.Order) *order.Order {
	clone := *o
	if o.DeletedAt != nil {
		deletedAt := *o.DeletedAt
		clone.DeletedAt = &deletedAt
	}
	// Clone prosthesis items
	clone.Prosthesis = make([]order.ProsthesisItem, len(o.Prosthesis))
	copy(clone.Prosthesis, o.Prosthesis)
	return &clone
}
