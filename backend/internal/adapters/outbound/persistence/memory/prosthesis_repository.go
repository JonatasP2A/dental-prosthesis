package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
)

// ProsthesisRepository is an in-memory implementation of the prosthesis repository
type ProsthesisRepository struct {
	mu   sync.RWMutex
	data map[string]*prosthesis.Prosthesis
}

// NewProsthesisRepository creates a new in-memory prosthesis repository
func NewProsthesisRepository() *ProsthesisRepository {
	return &ProsthesisRepository{
		data: make(map[string]*prosthesis.Prosthesis),
	}
}

// Create stores a new prosthesis
func (r *ProsthesisRepository) Create(ctx context.Context, p *prosthesis.Prosthesis) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[p.ID]; exists {
		return errors.ErrInternal // ID already exists
	}

	// Clone to avoid external modifications
	r.data[p.ID] = r.clone(p)
	return nil
}

// GetByID retrieves a prosthesis by ID (excludes soft-deleted)
func (r *ProsthesisRepository) GetByID(ctx context.Context, id string) (*prosthesis.Prosthesis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.data[id]
	if !exists || p.IsDeleted() {
		return nil, errors.ErrNotFound
	}

	return r.clone(p), nil
}

// Update updates an existing prosthesis
func (r *ProsthesisRepository) Update(ctx context.Context, p *prosthesis.Prosthesis) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[p.ID]
	if !exists || existing.IsDeleted() {
		return errors.ErrNotFound
	}

	r.data[p.ID] = r.clone(p)
	return nil
}

// Delete performs a soft delete on a prosthesis
func (r *ProsthesisRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.data[id]
	if !exists || p.IsDeleted() {
		return errors.ErrNotFound
	}

	p.Delete()
	return nil
}

// List retrieves all active (non-deleted) prostheses for a laboratory
func (r *ProsthesisRepository) List(ctx context.Context, laboratoryID string) ([]*prosthesis.Prosthesis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var prostheses []*prosthesis.Prosthesis
	for _, p := range r.data {
		if p.LaboratoryID == laboratoryID && !p.IsDeleted() {
			prostheses = append(prostheses, r.clone(p))
		}
	}

	return prostheses, nil
}

// FindByType retrieves prostheses filtered by type for a laboratory
func (r *ProsthesisRepository) FindByType(ctx context.Context, laboratoryID string, prosthesisType prosthesis.ProsthesisType) ([]*prosthesis.Prosthesis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var prostheses []*prosthesis.Prosthesis
	for _, p := range r.data {
		if p.LaboratoryID == laboratoryID && p.Type == prosthesisType && !p.IsDeleted() {
			prostheses = append(prostheses, r.clone(p))
		}
	}

	return prostheses, nil
}

// FindByMaterial retrieves prostheses filtered by material for a laboratory
func (r *ProsthesisRepository) FindByMaterial(ctx context.Context, laboratoryID string, material string) ([]*prosthesis.Prosthesis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var prostheses []*prosthesis.Prosthesis
	for _, p := range r.data {
		if p.LaboratoryID == laboratoryID && strings.EqualFold(p.Material, material) && !p.IsDeleted() {
			prostheses = append(prostheses, r.clone(p))
		}
	}

	return prostheses, nil
}

// clone creates a deep copy of a prosthesis to avoid external modifications
func (r *ProsthesisRepository) clone(p *prosthesis.Prosthesis) *prosthesis.Prosthesis {
	clone := *p
	if p.DeletedAt != nil {
		deletedAt := *p.DeletedAt
		clone.DeletedAt = &deletedAt
	}
	return &clone
}
