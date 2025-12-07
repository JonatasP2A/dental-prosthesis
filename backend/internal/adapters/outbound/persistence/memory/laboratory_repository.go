package memory

import (
	"context"
	"sync"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

// LaboratoryRepository is an in-memory implementation of the laboratory repository
type LaboratoryRepository struct {
	mu   sync.RWMutex
	data map[string]*laboratory.Laboratory
}

// NewLaboratoryRepository creates a new in-memory laboratory repository
func NewLaboratoryRepository() *LaboratoryRepository {
	return &LaboratoryRepository{
		data: make(map[string]*laboratory.Laboratory),
	}
}

// Create stores a new laboratory
func (r *LaboratoryRepository) Create(ctx context.Context, lab *laboratory.Laboratory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[lab.ID]; exists {
		return errors.ErrDuplicateEmail // ID already exists
	}

	// Clone to avoid external modifications
	r.data[lab.ID] = r.clone(lab)
	return nil
}

// GetByID retrieves a laboratory by ID (excludes soft-deleted)
func (r *LaboratoryRepository) GetByID(ctx context.Context, id string) (*laboratory.Laboratory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lab, exists := r.data[id]
	if !exists || lab.IsDeleted() {
		return nil, errors.ErrNotFound
	}

	return r.clone(lab), nil
}

// GetByEmail retrieves a laboratory by email (excludes soft-deleted)
func (r *LaboratoryRepository) GetByEmail(ctx context.Context, email string) (*laboratory.Laboratory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, lab := range r.data {
		if lab.Email == email && !lab.IsDeleted() {
			return r.clone(lab), nil
		}
	}

	return nil, errors.ErrNotFound
}

// Update updates an existing laboratory
func (r *LaboratoryRepository) Update(ctx context.Context, lab *laboratory.Laboratory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[lab.ID]
	if !exists || existing.IsDeleted() {
		return errors.ErrNotFound
	}

	r.data[lab.ID] = r.clone(lab)
	return nil
}

// Delete performs a soft delete on a laboratory
func (r *LaboratoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lab, exists := r.data[id]
	if !exists || lab.IsDeleted() {
		return errors.ErrNotFound
	}

	lab.Delete()
	return nil
}

// List retrieves all active (non-deleted) laboratories
func (r *LaboratoryRepository) List(ctx context.Context) ([]*laboratory.Laboratory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var labs []*laboratory.Laboratory
	for _, lab := range r.data {
		if !lab.IsDeleted() {
			labs = append(labs, r.clone(lab))
		}
	}

	return labs, nil
}

// clone creates a deep copy of a laboratory to avoid external modifications
func (r *LaboratoryRepository) clone(lab *laboratory.Laboratory) *laboratory.Laboratory {
	clone := *lab
	if lab.DeletedAt != nil {
		deletedAt := *lab.DeletedAt
		clone.DeletedAt = &deletedAt
	}
	return &clone
}

