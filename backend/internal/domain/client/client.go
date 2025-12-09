package client

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

// Client represents a dental clinic or dentist who orders prosthetic work
type Client struct {
	ID           string
	LaboratoryID string
	Name         string
	Email        string
	Phone        string
	Address      Address
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// Address represents a client's address
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// NewClient creates a new Client with validation
func NewClient(id, laboratoryID, name, email, phone string, address Address) (*Client, error) {
	client := &Client{
		ID:           id,
		LaboratoryID: laboratoryID,
		Name:         name,
		Email:        email,
		Phone:        phone,
		Address:      address,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := client.Validate(); err != nil {
		return nil, err
	}

	return client, nil
}

// Validate validates the client fields
func (c *Client) Validate() error {
	var validationErrors errors.ValidationErrors

	// Validate laboratory_id
	if strings.TrimSpace(c.LaboratoryID) == "" {
		validationErrors = append(validationErrors, errors.ValidationError{
			Field:   "laboratory_id",
			Message: "laboratory_id is required",
		})
	}

	// Validate name
	name := strings.TrimSpace(c.Name)
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
	email := strings.TrimSpace(c.Email)
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
	phone := strings.TrimSpace(c.Phone)
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
	if err := c.Address.Validate(); err != nil {
		if ve, ok := err.(errors.ValidationErrors); ok {
			validationErrors = append(validationErrors, ve...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// Update updates the client fields and sets UpdatedAt
func (c *Client) Update(name, email, phone string, address Address) error {
	c.Name = name
	c.Email = email
	c.Phone = phone
	c.Address = address
	c.UpdatedAt = time.Now().UTC()

	return c.Validate()
}

// Delete performs a soft delete by setting DeletedAt
func (c *Client) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}

// IsDeleted returns true if the client has been soft-deleted
func (c *Client) IsDeleted() bool {
	return c.DeletedAt != nil
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
