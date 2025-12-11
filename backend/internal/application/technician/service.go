package technician

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/ports/outbound"
)

// Service provides technician use cases
type Service struct {
	techRepo outbound.TechnicianRepository
	labRepo  outbound.LaboratoryRepository
	idGen    IDGenerator
}

// IDGenerator generates unique IDs
type IDGenerator interface {
	Generate() string
}

// NewService creates a new technician service
func NewService(techRepo outbound.TechnicianRepository, labRepo outbound.LaboratoryRepository, idGen IDGenerator) *Service {
	return &Service{
		techRepo: techRepo,
		labRepo:  labRepo,
		idGen:    idGen,
	}
}

// CreateInput represents the input for creating a technician
type CreateInput struct {
	LaboratoryID    string
	Name            string
	Email           string
	Phone           string
	Role            technician.Role
	Specializations []string
}

// CreateTechnician creates a new technician
func (s *Service) CreateTechnician(ctx context.Context, input CreateInput) (*technician.Technician, error) {
	// Validate laboratory exists
	_, err := s.labRepo.GetByID(ctx, input.LaboratoryID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check if email already exists within the laboratory
	existing, err := s.techRepo.GetByEmail(ctx, input.LaboratoryID, input.Email)
	if err != nil && err != errors.ErrNotFound {
		return nil, errors.ErrInternal
	}
	if existing != nil {
		return nil, errors.ErrDuplicateEmail
	}

	// Create new technician
	id := s.idGen.Generate()
	tech, err := technician.NewTechnician(id, input.LaboratoryID, input.Name, input.Email, input.Phone, input.Role, input.Specializations)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.techRepo.Create(ctx, tech); err != nil {
		return nil, errors.ErrInternal
	}

	return tech, nil
}

// GetTechnician retrieves a technician by ID (laboratory-scoped)
func (s *Service) GetTechnician(ctx context.Context, id, laboratoryID string) (*technician.Technician, error) {
	tech, err := s.techRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if tech.LaboratoryID != laboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	return tech, nil
}

// UpdateInput represents the input for updating a technician
type UpdateInput struct {
	ID              string
	LaboratoryID    string
	Name            string
	Email           string
	Phone           string
	Role            technician.Role
	Specializations []string
}

// UpdateTechnician updates an existing technician
func (s *Service) UpdateTechnician(ctx context.Context, input UpdateInput) (*technician.Technician, error) {
	// Get existing technician
	tech, err := s.techRepo.GetByID(ctx, input.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if tech.LaboratoryID != input.LaboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	// Check if email changed and already exists
	if tech.Email != input.Email {
		existing, err := s.techRepo.GetByEmail(ctx, input.LaboratoryID, input.Email)
		if err != nil && err != errors.ErrNotFound {
			return nil, errors.ErrInternal
		}
		if existing != nil && existing.ID != tech.ID {
			return nil, errors.ErrDuplicateEmail
		}
	}

	// Update technician
	if err := tech.Update(input.Name, input.Email, input.Phone, input.Role, input.Specializations); err != nil {
		return nil, err
	}

	// Persist
	if err := s.techRepo.Update(ctx, tech); err != nil {
		return nil, errors.ErrInternal
	}

	return tech, nil
}

// ListTechnicians retrieves all active technicians for a laboratory, optionally filtered by role
func (s *Service) ListTechnicians(ctx context.Context, laboratoryID string, role *technician.Role) ([]*technician.Technician, error) {
	var techs []*technician.Technician
	var err error

	if role != nil {
		techs, err = s.techRepo.ListByRole(ctx, laboratoryID, *role)
	} else {
		techs, err = s.techRepo.List(ctx, laboratoryID)
	}

	if err != nil {
		return nil, errors.ErrInternal
	}

	return techs, nil
}

// DeleteTechnician performs a soft delete on a technician (laboratory-scoped)
func (s *Service) DeleteTechnician(ctx context.Context, id, laboratoryID string) error {
	// Check if technician exists and belongs to the laboratory
	tech, err := s.techRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.ErrInternal
	}

	// Check laboratory scope
	if tech.LaboratoryID != laboratoryID {
		return errors.ErrNotFound // Security: don't reveal existence
	}

	// Delete
	if err := s.techRepo.Delete(ctx, id); err != nil {
		return errors.ErrInternal
	}

	return nil
}
