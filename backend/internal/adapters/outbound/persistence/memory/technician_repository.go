package memory

import (
	"context"
	"sync"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
)

// TechnicianRepository is an in-memory implementation of the technician repository
type TechnicianRepository struct {
	mu   sync.RWMutex
	data map[string]*technician.Technician
}

// NewTechnicianRepository creates a new in-memory technician repository
func NewTechnicianRepository() *TechnicianRepository {
	return &TechnicianRepository{
		data: make(map[string]*technician.Technician),
	}
}

// Create stores a new technician
func (r *TechnicianRepository) Create(ctx context.Context, tech *technician.Technician) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[tech.ID]; exists {
		return errors.ErrInternal // ID already exists
	}

	// Clone to avoid external modifications
	r.data[tech.ID] = r.clone(tech)
	return nil
}

// GetByID retrieves a technician by ID (excludes soft-deleted)
func (r *TechnicianRepository) GetByID(ctx context.Context, id string) (*technician.Technician, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tech, exists := r.data[id]
	if !exists || tech.IsDeleted() {
		return nil, errors.ErrNotFound
	}

	return r.clone(tech), nil
}

// GetByEmail retrieves a technician by email within a laboratory (excludes soft-deleted)
func (r *TechnicianRepository) GetByEmail(ctx context.Context, laboratoryID, email string) (*technician.Technician, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tech := range r.data {
		if tech.Email == email && tech.LaboratoryID == laboratoryID && !tech.IsDeleted() {
			return r.clone(tech), nil
		}
	}

	return nil, errors.ErrNotFound
}

// Update updates an existing technician
func (r *TechnicianRepository) Update(ctx context.Context, tech *technician.Technician) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[tech.ID]
	if !exists || existing.IsDeleted() {
		return errors.ErrNotFound
	}

	r.data[tech.ID] = r.clone(tech)
	return nil
}

// Delete performs a soft delete on a technician
func (r *TechnicianRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tech, exists := r.data[id]
	if !exists || tech.IsDeleted() {
		return errors.ErrNotFound
	}

	tech.Delete()
	return nil
}

// List retrieves all active (non-deleted) technicians for a laboratory
func (r *TechnicianRepository) List(ctx context.Context, laboratoryID string) ([]*technician.Technician, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var techs []*technician.Technician
	for _, tech := range r.data {
		if tech.LaboratoryID == laboratoryID && !tech.IsDeleted() {
			techs = append(techs, r.clone(tech))
		}
	}

	return techs, nil
}

// ListByRole retrieves all active technicians for a laboratory filtered by role
func (r *TechnicianRepository) ListByRole(ctx context.Context, laboratoryID string, role technician.Role) ([]*technician.Technician, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var techs []*technician.Technician
	for _, tech := range r.data {
		if tech.LaboratoryID == laboratoryID && tech.Role == role && !tech.IsDeleted() {
			techs = append(techs, r.clone(tech))
		}
	}

	return techs, nil
}

// clone creates a deep copy of a technician to avoid external modifications
func (r *TechnicianRepository) clone(tech *technician.Technician) *technician.Technician {
	clone := *tech
	if tech.DeletedAt != nil {
		deletedAt := *tech.DeletedAt
		clone.DeletedAt = &deletedAt
	}
	if tech.Specializations != nil {
		clone.Specializations = make([]string, len(tech.Specializations))
		copy(clone.Specializations, tech.Specializations)
	}
	return &clone
}
