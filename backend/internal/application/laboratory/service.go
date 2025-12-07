package laboratory

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/ports/outbound"
)

// Service provides laboratory use cases
type Service struct {
	repo      outbound.LaboratoryRepository
	idGen     IDGenerator
}

// IDGenerator generates unique IDs
type IDGenerator interface {
	Generate() string
}

// NewService creates a new laboratory service
func NewService(repo outbound.LaboratoryRepository, idGen IDGenerator) *Service {
	return &Service{
		repo:  repo,
		idGen: idGen,
	}
}

// CreateInput represents the input for creating a laboratory
type CreateInput struct {
	Name    string
	Email   string
	Phone   string
	Address laboratory.Address
}

// CreateLaboratory creates a new laboratory
func (s *Service) CreateLaboratory(ctx context.Context, input CreateInput) (*laboratory.Laboratory, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != errors.ErrNotFound {
		return nil, errors.ErrInternal
	}
	if existing != nil {
		return nil, errors.ErrDuplicateEmail
	}

	// Create new laboratory
	id := s.idGen.Generate()
	lab, err := laboratory.NewLaboratory(id, input.Name, input.Email, input.Phone, input.Address)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.repo.Create(ctx, lab); err != nil {
		return nil, errors.ErrInternal
	}

	return lab, nil
}

// GetLaboratory retrieves a laboratory by ID
func (s *Service) GetLaboratory(ctx context.Context, id string) (*laboratory.Laboratory, error) {
	lab, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	return lab, nil
}

// UpdateInput represents the input for updating a laboratory
type UpdateInput struct {
	ID      string
	Name    string
	Email   string
	Phone   string
	Address laboratory.Address
}

// UpdateLaboratory updates an existing laboratory
func (s *Service) UpdateLaboratory(ctx context.Context, input UpdateInput) (*laboratory.Laboratory, error) {
	// Get existing laboratory
	lab, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check if email changed and already exists
	if lab.Email != input.Email {
		existing, err := s.repo.GetByEmail(ctx, input.Email)
		if err != nil && err != errors.ErrNotFound {
			return nil, errors.ErrInternal
		}
		if existing != nil && existing.ID != lab.ID {
			return nil, errors.ErrDuplicateEmail
		}
	}

	// Update laboratory
	if err := lab.Update(input.Name, input.Email, input.Phone, input.Address); err != nil {
		return nil, err
	}

	// Persist
	if err := s.repo.Update(ctx, lab); err != nil {
		return nil, errors.ErrInternal
	}

	return lab, nil
}

// ListLaboratories retrieves all active laboratories
func (s *Service) ListLaboratories(ctx context.Context) ([]*laboratory.Laboratory, error) {
	labs, err := s.repo.List(ctx)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return labs, nil
}

// DeleteLaboratory performs a soft delete on a laboratory
func (s *Service) DeleteLaboratory(ctx context.Context, id string) error {
	// Check if laboratory exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.ErrInternal
	}

	// Delete
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.ErrInternal
	}

	return nil
}

