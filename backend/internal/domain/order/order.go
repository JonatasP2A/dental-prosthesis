package order

import (
	"strings"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

// Order represents a prosthesis work order from a client
type Order struct {
	ID           string
	ClientID     string
	LaboratoryID string
	Status       Status
	Prosthesis   []ProsthesisItem
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// Status represents the workflow state of an order
type Status string

const (
	StatusReceived     Status = "received"
	StatusInProduction Status = "in_production"
	StatusQualityCheck Status = "quality_check"
	StatusReady        Status = "ready"
	StatusDelivered    Status = "delivered"
	StatusRevision     Status = "revision"
)

// validTransitions defines valid status transitions
var validTransitions = map[Status][]Status{
	StatusReceived:     {StatusInProduction},
	StatusInProduction: {StatusQualityCheck},
	StatusQualityCheck: {StatusReady, StatusRevision},
	StatusReady:        {StatusDelivered, StatusRevision},
	StatusRevision:     {StatusInProduction},
	StatusDelivered:    {}, // Terminal state
}

// AllStatuses returns all valid status values
func AllStatuses() []Status {
	return []Status{
		StatusReceived,
		StatusInProduction,
		StatusQualityCheck,
		StatusReady,
		StatusDelivered,
		StatusRevision,
	}
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(s string) bool {
	for _, status := range AllStatuses() {
		if string(status) == s {
			return true
		}
	}
	return false
}

// ProsthesisItem represents a prosthesis item in an order
type ProsthesisItem struct {
	Type     string
	Material string
	Shade    string
	Quantity int
	Notes    string
}

// NewOrder creates a new Order with validation
func NewOrder(id, clientID, laboratoryID string, items []ProsthesisItem) (*Order, error) {
	order := &Order{
		ID:           id,
		ClientID:     clientID,
		LaboratoryID: laboratoryID,
		Status:       StatusReceived, // Initial status is always "received"
		Prosthesis:   items,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}

// Validate validates the order fields
func (o *Order) Validate() error {
	var validationErrors errors.ValidationErrors

	// Validate client_id
	if strings.TrimSpace(o.ClientID) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "client_id",
			Message: "client_id is required",
		})
	}

	// Validate laboratory_id
	if strings.TrimSpace(o.LaboratoryID) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "laboratory_id",
			Message: "laboratory_id is required",
		})
	}

	// Validate prosthesis items
	if len(o.Prosthesis) == 0 {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "prosthesis",
			Message: "at least one prosthesis item is required",
		})
	}

	// Validate each prosthesis item
	for i, item := range o.Prosthesis {
		if err := item.Validate(i); err != nil {
			if ve, ok := err.(errors.ValidationErrors); ok {
				validationErrors = append(validationErrors, ve...)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// Update updates the order fields (except status) and sets UpdatedAt
func (o *Order) Update(items []ProsthesisItem) error {
	o.Prosthesis = items
	o.UpdatedAt = time.Now().UTC()

	return o.Validate()
}

// CanTransitionTo checks if a status transition is valid
func (o *Order) CanTransitionTo(newStatus Status) bool {
	allowedTransitions, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

// UpdateStatus updates the order status with workflow validation
func (o *Order) UpdateStatus(newStatus Status) error {
	if !o.CanTransitionTo(newStatus) {
		return errors.ErrInvalidStatusTransition
	}

	o.Status = newStatus
	o.UpdatedAt = time.Now().UTC()
	return nil
}

// Delete performs a soft delete by setting DeletedAt
func (o *Order) Delete() {
	now := time.Now().UTC()
	o.DeletedAt = &now
}

// IsDeleted returns true if the order has been soft-deleted
func (o *Order) IsDeleted() bool {
	return o.DeletedAt != nil
}

// Validate validates the prosthesis item fields
func (p *ProsthesisItem) Validate(index int) error {
	var validationErrors errors.ValidationErrors

	if strings.TrimSpace(p.Type) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "prosthesis[" + string(rune('0'+index)) + "].type",
			Message: "type is required",
		})
	}

	if strings.TrimSpace(p.Material) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "prosthesis[" + string(rune('0'+index)) + "].material",
			Message: "material is required",
		})
	}

	if p.Quantity <= 0 {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "prosthesis[" + string(rune('0'+index)) + "].quantity",
			Message: "quantity must be greater than 0",
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
