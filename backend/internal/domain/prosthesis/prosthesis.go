package prosthesis

import (
	"strings"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

// ProsthesisType represents the type of dental prosthesis
type ProsthesisType string

const (
	ProsthesisTypeCrown            ProsthesisType = "crown"
	ProsthesisTypeBridge           ProsthesisType = "bridge"
	ProsthesisTypeCompleteDenture  ProsthesisType = "complete_denture"
	ProsthesisTypePartialDenture   ProsthesisType = "partial_denture"
	ProsthesisTypeImplant          ProsthesisType = "implant"
	ProsthesisTypeVeneer           ProsthesisType = "veneer"
	ProsthesisTypeInlay            ProsthesisType = "inlay"
	ProsthesisTypeOnlay            ProsthesisType = "onlay"
)

// ValidProsthesisTypes returns a map of valid prosthesis types
func ValidProsthesisTypes() map[ProsthesisType]bool {
	return map[ProsthesisType]bool{
		ProsthesisTypeCrown:           true,
		ProsthesisTypeBridge:          true,
		ProsthesisTypeCompleteDenture: true,
		ProsthesisTypePartialDenture:  true,
		ProsthesisTypeImplant:         true,
		ProsthesisTypeVeneer:          true,
		ProsthesisTypeInlay:           true,
		ProsthesisTypeOnlay:           true,
	}
}

// IsValid checks if the prosthesis type is valid
func (pt ProsthesisType) IsValid() bool {
	return ValidProsthesisTypes()[pt]
}

// Prosthesis represents a dental prosthetic item
type Prosthesis struct {
	ID            string
	LaboratoryID  string
	Type          ProsthesisType
	Material      string
	Shade         string
	Specifications string
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

// NewProsthesis creates a new Prosthesis with validation
func NewProsthesis(id, laboratoryID string, prosthesisType ProsthesisType, material, shade, specifications, notes string) (*Prosthesis, error) {
	prosthesis := &Prosthesis{
		ID:             id,
		LaboratoryID:   laboratoryID,
		Type:           prosthesisType,
		Material:       material,
		Shade:          shade,
		Specifications: specifications,
		Notes:          notes,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := prosthesis.Validate(); err != nil {
		return nil, err
	}

	return prosthesis, nil
}

// Validate validates the prosthesis fields
func (p *Prosthesis) Validate() error {
	var validationErrors errors.ValidationErrors

	// Validate laboratory_id
	if strings.TrimSpace(p.LaboratoryID) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "laboratory_id",
			Message: "laboratory_id is required",
		})
	}

	// Validate type
	if strings.TrimSpace(string(p.Type)) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "type",
			Message: "type is required",
		})
	} else if !p.Type.IsValid() {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "type",
			Message: "invalid prosthesis type",
		})
	}

	// Validate material
	if strings.TrimSpace(p.Material) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "material",
			Message: "material is required",
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// Update updates the prosthesis fields and sets UpdatedAt
func (p *Prosthesis) Update(prosthesisType ProsthesisType, material, shade, specifications, notes string) error {
	p.Type = prosthesisType
	p.Material = material
	p.Shade = shade
	p.Specifications = specifications
	p.Notes = notes
	p.UpdatedAt = time.Now().UTC()

	return p.Validate()
}

// Delete performs a soft delete by setting DeletedAt
func (p *Prosthesis) Delete() {
	now := time.Now().UTC()
	p.DeletedAt = &now
}

// IsDeleted returns true if the prosthesis has been soft-deleted
func (p *Prosthesis) IsDeleted() bool {
	return p.DeletedAt != nil
}
