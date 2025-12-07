package laboratory

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

// Laboratory represents a dental prosthesis lab (tenant)
type Laboratory struct {
	ID        string
	Name      string
	Email     string
	Phone     string
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Address represents a laboratory's address
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// NewLaboratory creates a new Laboratory with validation
func NewLaboratory(id, name, email, phone string, address Address) (*Laboratory, error) {
	lab := &Laboratory{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := lab.Validate(); err != nil {
		return nil, err
	}

	return lab, nil
}

// Validate validates the laboratory fields
func (l *Laboratory) Validate() error {
	var validationErrors errors.ValidationErrors

	// Validate name
	name := strings.TrimSpace(l.Name)
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
	email := strings.TrimSpace(l.Email)
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
	phone := strings.TrimSpace(l.Phone)
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

	// Validate address
	if err := l.Address.Validate(); err != nil {
		if ve, ok := err.(errors.ValidationErrors); ok {
			validationErrors = append(validationErrors, ve...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// Update updates the laboratory fields and sets UpdatedAt
func (l *Laboratory) Update(name, email, phone string, address Address) error {
	l.Name = name
	l.Email = email
	l.Phone = phone
	l.Address = address
	l.UpdatedAt = time.Now().UTC()

	return l.Validate()
}

// Delete performs a soft delete by setting DeletedAt
func (l *Laboratory) Delete() {
	now := time.Now().UTC()
	l.DeletedAt = &now
}

// IsDeleted returns true if the laboratory has been soft-deleted
func (l *Laboratory) IsDeleted() bool {
	return l.DeletedAt != nil
}

// Validate validates the address fields
func (a *Address) Validate() error {
	var validationErrors errors.ValidationErrors

	if strings.TrimSpace(a.Street) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "address.street",
			Message: "street is required",
		})
	}

	if strings.TrimSpace(a.City) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "address.city",
			Message: "city is required",
		})
	}

	if strings.TrimSpace(a.State) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "address.state",
			Message: "state is required",
		})
	}

	if strings.TrimSpace(a.PostalCode) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "address.postal_code",
			Message: "postal code is required",
		})
	}

	if strings.TrimSpace(a.Country) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "address.country",
			Message: "country is required",
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
