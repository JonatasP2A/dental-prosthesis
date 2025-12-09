package order

import (
	"context"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/ports/outbound"
)

// Service provides order use cases
type Service struct {
	orderRepo  outbound.OrderRepository
	clientRepo outbound.ClientRepository
	idGen      IDGenerator
}

// IDGenerator generates unique IDs
type IDGenerator interface {
	Generate() string
}

// NewService creates a new order service
func NewService(orderRepo outbound.OrderRepository, clientRepo outbound.ClientRepository, idGen IDGenerator) *Service {
	return &Service{
		orderRepo:  orderRepo,
		clientRepo: clientRepo,
		idGen:      idGen,
	}
}

// CreateInput represents the input for creating an order
type CreateInput struct {
	ClientID     string
	LaboratoryID string // Used for validation
	Prosthesis   []order.ProsthesisItem
}

// CreateOrder creates a new order
func (s *Service) CreateOrder(ctx context.Context, input CreateInput) (*order.Order, error) {
	// Validate client exists and belongs to the laboratory
	client, err := s.clientRepo.GetByID(ctx, input.ClientID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if client.LaboratoryID != input.LaboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	// Create new order with laboratory_id derived from client
	id := s.idGen.Generate()
	o, err := order.NewOrder(id, input.ClientID, client.LaboratoryID, input.Prosthesis)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.orderRepo.Create(ctx, o); err != nil {
		return nil, errors.ErrInternal
	}

	return o, nil
}

// GetOrder retrieves an order by ID (laboratory-scoped)
func (s *Service) GetOrder(ctx context.Context, id, laboratoryID string) (*order.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if o.LaboratoryID != laboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	return o, nil
}

// UpdateInput represents the input for updating an order
type UpdateInput struct {
	ID           string
	LaboratoryID string
	Prosthesis   []order.ProsthesisItem
}

// UpdateOrder updates an existing order (excluding status)
func (s *Service) UpdateOrder(ctx context.Context, input UpdateInput) (*order.Order, error) {
	// Get existing order
	o, err := s.orderRepo.GetByID(ctx, input.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if o.LaboratoryID != input.LaboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	// Update order
	if err := o.Update(input.Prosthesis); err != nil {
		return nil, err
	}

	// Persist
	if err := s.orderRepo.Update(ctx, o); err != nil {
		return nil, errors.ErrInternal
	}

	return o, nil
}

// UpdateStatusInput represents the input for updating an order's status
type UpdateStatusInput struct {
	ID           string
	LaboratoryID string
	Status       order.Status
}

// UpdateOrderStatus updates an order's status with workflow validation
func (s *Service) UpdateOrderStatus(ctx context.Context, input UpdateStatusInput) (*order.Order, error) {
	// Get existing order
	o, err := s.orderRepo.GetByID(ctx, input.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if o.LaboratoryID != input.LaboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	// Update status with workflow validation
	if err := o.UpdateStatus(input.Status); err != nil {
		return nil, err
	}

	// Persist
	if err := s.orderRepo.Update(ctx, o); err != nil {
		return nil, errors.ErrInternal
	}

	return o, nil
}

// ListOrders retrieves all active orders for a laboratory
func (s *Service) ListOrders(ctx context.Context, laboratoryID string) ([]*order.Order, error) {
	orders, err := s.orderRepo.List(ctx, laboratoryID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return orders, nil
}

// ListOrdersByClient retrieves all active orders for a specific client (laboratory-scoped)
func (s *Service) ListOrdersByClient(ctx context.Context, clientID, laboratoryID string) ([]*order.Order, error) {
	// Validate client exists and belongs to the laboratory
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	// Check laboratory scope
	if client.LaboratoryID != laboratoryID {
		return nil, errors.ErrNotFound // Security: don't reveal existence
	}

	orders, err := s.orderRepo.ListByClientID(ctx, clientID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return orders, nil
}

// DeleteOrder performs a soft delete on an order (laboratory-scoped)
func (s *Service) DeleteOrder(ctx context.Context, id, laboratoryID string) error {
	// Check if order exists and belongs to the laboratory
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.ErrInternal
	}

	// Check laboratory scope
	if o.LaboratoryID != laboratoryID {
		return errors.ErrNotFound // Security: don't reveal existence
	}

	// Delete
	if err := s.orderRepo.Delete(ctx, id); err != nil {
		return errors.ErrInternal
	}

	return nil
}
