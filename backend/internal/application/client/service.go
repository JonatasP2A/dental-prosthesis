package client

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/ports/outbound"
)

// Service provides client use cases
type Service struct {
	clientRepo outbound.ClientRepository
	labRepo    outbound.LaboratoryRepository
	idGen      IDGenerator
}

// IDGenerator generates unique IDs
type IDGenerator interface {
	Generate() string
}

// NewService creates a new client service
func NewService(clientRepo outbound.ClientRepository, labRepo outbound.LaboratoryRepository, idGen IDGenerator) *Service {
	return &Service{
		clientRepo: clientRepo,
		labRepo:    labRepo,
		idGen:      idGen,
	}
}

// CreateInput represents the input for creating a client
type CreateInput struct {
	LaboratoryID string
	Name         string
	Email        string
	Phone        string
	Address      client.Address
}

// CreateClient creates a new client
func (s *Service) CreateClient(ctx context.Context, input CreateInput) (*client.Client, error) {
	// Validate laboratory exists
	_, err := s.labRepo.GetByID(ctx, input.LaboratoryID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check if email already exists within the laboratory
	existing, err := s.clientRepo.GetByEmail(ctx, input.LaboratoryID, input.Email)
	if err != nil && err != errors.ErrNotFound {
		return nil, errors.ErrInternal
	}
	if existing != nil {
		return nil, errors.ErrDuplicateEmail
	}

	// Create new client
	id := s.idGen.Generate()
	c, err := client.NewClient(id, input.LaboratoryID, input.Name, input.Email, input.Phone, input.Address)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.clientRepo.Create(ctx, c); err != nil {
		return nil, errors.ErrInternal
	}

	return c, nil
}

// GetClient retrieves a client by ID (laboratory-scoped)
func (s *Service) GetClient(ctx context.Context, id, laboratoryID string) (*client.Client, error) {
	c, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if c.LaboratoryID != laboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	return c, nil
}

// UpdateInput represents the input for updating a client
type UpdateInput struct {
	ID           string
	LaboratoryID string
	Name         string
	Email        string
	Phone        string
	Address      client.Address
}

// UpdateClient updates an existing client
func (s *Service) UpdateClient(ctx context.Context, input UpdateInput) (*client.Client, error) {
	// Get existing client
	c, err := s.clientRepo.GetByID(ctx, input.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if c.LaboratoryID != input.LaboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	// Check if email changed and already exists
	if c.Email != input.Email {
		existing, err := s.clientRepo.GetByEmail(ctx, input.LaboratoryID, input.Email)
		if err != nil && err != errors.ErrNotFound {
			return nil, errors.ErrInternal
		}
		if existing != nil && existing.ID != c.ID {
			return nil, errors.ErrDuplicateEmail
		}
	}

	// Update client
	if err := c.Update(input.Name, input.Email, input.Phone, input.Address); err != nil {
		return nil, err
	}

	// Persist
	if err := s.clientRepo.Update(ctx, c); err != nil {
		return nil, errors.ErrInternal
	}

	return c, nil
}

// ListClients retrieves all active clients for a laboratory
func (s *Service) ListClients(ctx context.Context, laboratoryID string) ([]*client.Client, error) {
	clients, err := s.clientRepo.List(ctx, laboratoryID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return clients, nil
}

// DeleteClient performs a soft delete on a client (laboratory-scoped)
func (s *Service) DeleteClient(ctx context.Context, id, laboratoryID string) error {
	// Check if client exists and belongs to the laboratory
	c, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.ErrInternal
	}

	// Check laboratory scope
	if c.LaboratoryID != laboratoryID {
		return errors.ErrNotFound // Security: don't reveal existence
	}

	// Delete
	if err := s.clientRepo.Delete(ctx, id); err != nil {
		return errors.ErrInternal
	}

	return nil
}
