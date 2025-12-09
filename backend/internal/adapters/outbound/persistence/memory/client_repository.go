package memory

import (
	"context"
	"sync"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

// ClientRepository is an in-memory implementation of the client repository
type ClientRepository struct {
	mu   sync.RWMutex
	data map[string]*client.Client
}

// NewClientRepository creates a new in-memory client repository
func NewClientRepository() *ClientRepository {
	return &ClientRepository{
		data: make(map[string]*client.Client),
	}
}

// Create stores a new client
func (r *ClientRepository) Create(ctx context.Context, c *client.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[c.ID]; exists {
		return errors.ErrDuplicateEmail // ID already exists
	}

	// Clone to avoid external modifications
	r.data[c.ID] = r.clone(c)
	return nil
}

// GetByID retrieves a client by ID (excludes soft-deleted)
func (r *ClientRepository) GetByID(ctx context.Context, id string) (*client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.data[id]
	if !exists || c.IsDeleted() {
		return nil, errors.ErrNotFound
	}

	return r.clone(c), nil
}

// GetByEmail retrieves a client by email within a laboratory (excludes soft-deleted)
func (r *ClientRepository) GetByEmail(ctx context.Context, laboratoryID, email string) (*client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.data {
		if c.Email == email && c.LaboratoryID == laboratoryID && !c.IsDeleted() {
			return r.clone(c), nil
		}
	}

	return nil, errors.ErrNotFound
}

// Update updates an existing client
func (r *ClientRepository) Update(ctx context.Context, c *client.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[c.ID]
	if !exists || existing.IsDeleted() {
		return errors.ErrNotFound
	}

	r.data[c.ID] = r.clone(c)
	return nil
}

// Delete performs a soft delete on a client
func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, exists := r.data[id]
	if !exists || c.IsDeleted() {
		return errors.ErrNotFound
	}

	c.Delete()
	return nil
}

// List retrieves all active (non-deleted) clients for a laboratory
func (r *ClientRepository) List(ctx context.Context, laboratoryID string) ([]*client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var clients []*client.Client
	for _, c := range r.data {
		if c.LaboratoryID == laboratoryID && !c.IsDeleted() {
			clients = append(clients, r.clone(c))
		}
	}

	return clients, nil
}

// clone creates a deep copy of a client to avoid external modifications
func (r *ClientRepository) clone(c *client.Client) *client.Client {
	clone := *c
	if c.DeletedAt != nil {
		deletedAt := *c.DeletedAt
		clone.DeletedAt = &deletedAt
	}
	return &clone
}
