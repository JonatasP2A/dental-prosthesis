package prosthesis

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/ports/outbound"
)

// Service provides prosthesis use cases
type Service struct {
	prosthesisRepo outbound.ProsthesisRepository
	labRepo        outbound.LaboratoryRepository
	idGen          IDGenerator
}

// IDGenerator generates unique IDs
type IDGenerator interface {
	Generate() string
}

// NewService creates a new prosthesis service
func NewService(prosthesisRepo outbound.ProsthesisRepository, labRepo outbound.LaboratoryRepository, idGen IDGenerator) *Service {
	return &Service{
		prosthesisRepo: prosthesisRepo,
		labRepo:        labRepo,
		idGen:          idGen,
	}
}

// CreateInput represents the input for creating a prosthesis
type CreateInput struct {
	LaboratoryID   string
	Type           prosthesis.ProsthesisType
	Material       string
	Shade          string
	Specifications string
	Notes          string
}

// CreateProsthesis creates a new prosthesis
func (s *Service) CreateProsthesis(ctx context.Context, input CreateInput) (*prosthesis.Prosthesis, error) {
	// Validate laboratory exists
	_, err := s.labRepo.GetByID(ctx, input.LaboratoryID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Create new prosthesis
	id := s.idGen.Generate()
	p, err := prosthesis.NewProsthesis(id, input.LaboratoryID, input.Type, input.Material, input.Shade, input.Specifications, input.Notes)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.prosthesisRepo.Create(ctx, p); err != nil {
		return nil, errors.ErrInternal
	}

	return p, nil
}

// GetProsthesis retrieves a prosthesis by ID (laboratory-scoped)
func (s *Service) GetProsthesis(ctx context.Context, id, laboratoryID string) (*prosthesis.Prosthesis, error) {
	p, err := s.prosthesisRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if p.LaboratoryID != laboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	return p, nil
}

// UpdateInput represents the input for updating a prosthesis
type UpdateInput struct {
	ID             string
	LaboratoryID   string
	Type           prosthesis.ProsthesisType
	Material       string
	Shade          string
	Specifications string
	Notes          string
}

// UpdateProsthesis updates an existing prosthesis
func (s *Service) UpdateProsthesis(ctx context.Context, input UpdateInput) (*prosthesis.Prosthesis, error) {
	// Get existing prosthesis
	p, err := s.prosthesisRepo.GetByID(ctx, input.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if p.LaboratoryID != input.LaboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	// Update prosthesis
	if err := p.Update(input.Type, input.Material, input.Shade, input.Specifications, input.Notes); err != nil {
		return nil, err
	}

	// Persist
	if err := s.prosthesisRepo.Update(ctx, p); err != nil {
		return nil, errors.ErrInternal
	}

	return p, nil
}

// ListProstheses retrieves all active prostheses for a laboratory, optionally filtered by type or material
func (s *Service) ListProstheses(ctx context.Context, laboratoryID string, prosthesisType *prosthesis.ProsthesisType, material *string) ([]*prosthesis.Prosthesis, error) {
	var prostheses []*prosthesis.Prosthesis
	var err error

	if prosthesisType != nil {
		prostheses, err = s.prosthesisRepo.FindByType(ctx, laboratoryID, *prosthesisType)
	} else if material != nil {
		prostheses, err = s.prosthesisRepo.FindByMaterial(ctx, laboratoryID, *material)
	} else {
		prostheses, err = s.prosthesisRepo.List(ctx, laboratoryID)
	}

	if err != nil {
		return nil, errors.ErrInternal
	}

	return prostheses, nil
}

// DeleteProsthesis performs a soft delete on a prosthesis (laboratory-scoped)
func (s *Service) DeleteProsthesis(ctx context.Context, id, laboratoryID string) error {
	// Check if prosthesis exists and belongs to the laboratory
	p, err := s.prosthesisRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.ErrInternal
	}

	// Check laboratory scope
	if p.LaboratoryID != laboratoryID {
		return errors.ErrNotFound // Security: don't reveal existence
	}

	// Delete
	if err := s.prosthesisRepo.Delete(ctx, id); err != nil {
		return errors.ErrInternal
	}

	return nil
}
