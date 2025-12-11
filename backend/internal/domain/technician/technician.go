package technician

import (
	"regexp"
	"strings"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

// Role represents a technician's role
type Role string

const (
	RoleSeniorTechnician Role = "senior_technician"
	RoleTechnician       Role = "technician"
	RoleApprentice       Role = "apprentice"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	return r == RoleSeniorTechnician || r == RoleTechnician || r == RoleApprentice
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

// Technician represents a laboratory technician
type Technician struct {
	ID             string
	LaboratoryID   string
	Name           string
	Email          string
	Phone          string
	Role           Role
	Specializations []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// NewTechnician creates a new Technician with validation
func NewTechnician(id, laboratoryID, name, email, phone string, role Role, specializations []string) (*Technician, error) {
	tech := &Technician{
		ID:              id,
		LaboratoryID:    laboratoryID,
		Name:            name,
		Email:           email,
		Phone:           phone,
		Role:            role,
		Specializations: specializations,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := tech.Validate(); err != nil {
		return nil, err
	}

	return tech, nil
}

// Validate validates the technician fields
func (t *Technician) Validate() error {
	var validationErrors errors.ValidationErrors

	// Validate laboratory_id
	if strings.TrimSpace(t.LaboratoryID) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "laboratory_id",
			Message: "laboratory_id is required",
		})
	}

	// Validate name
	name := strings.TrimSpace(t.Name)
	if name == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	} else if len(name) > 200 {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "name",
			Message: "name must be at most 200 characters",
		})
	}

	// Validate email
	email := strings.TrimSpace(t.Email)
	if email == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "email",
			Message: "email is required",
		})
	} else if !emailRegex.MatchString(email) {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "email",
			Message: "invalid email format",
		})
	}

	// Validate phone
	phone := strings.TrimSpace(t.Phone)
	if phone == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "phone",
			Message: "phone is required",
		})
	} else if !phoneRegex.MatchString(strings.ReplaceAll(phone, " ", "")) {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "phone",
			Message: "invalid phone format",
		})
	}

	// Validate role
	if strings.TrimSpace(string(t.Role)) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "role",
			Message: "role is required",
		})
	} else if !t.Role.IsValid() {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "role",
			Message: "invalid role. Must be one of: senior_technician, technician, apprentice",
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// Update updates the technician fields and sets UpdatedAt
func (t *Technician) Update(name, email, phone string, role Role, specializations []string) error {
	t.Name = name
	t.Email = email
	t.Phone = phone
	t.Role = role
	t.Specializations = specializations
	t.UpdatedAt = time.Now().UTC()

	return t.Validate()
}

// Delete performs a soft delete by setting DeletedAt
func (t *Technician) Delete() {
	now := time.Now().UTC()
	t.DeletedAt = &now
}

// IsDeleted returns true if the technician has been soft-deleted
func (t *Technician) IsDeleted() bool {
	return t.DeletedAt != nil
}
