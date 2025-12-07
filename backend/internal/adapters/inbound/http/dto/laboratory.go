package dto

import (
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

// CreateLaboratoryRequest represents the request body for creating a laboratory
type CreateLaboratoryRequest struct {
	Name    string         `json:"name" binding:"required"`
	Email   string         `json:"email" binding:"required,email"`
	Phone   string         `json:"phone" binding:"required"`
	Address AddressRequest `json:"address" binding:"required"`
}

// AddressRequest represents the address in request body
type AddressRequest struct {
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	Country    string `json:"country" binding:"required"`
}

// UpdateLaboratoryRequest represents the request body for updating a laboratory
type UpdateLaboratoryRequest struct {
	Name    string         `json:"name" binding:"required"`
	Email   string         `json:"email" binding:"required,email"`
	Phone   string         `json:"phone" binding:"required"`
	Address AddressRequest `json:"address" binding:"required"`
}

// LaboratoryResponse represents the response body for a laboratory
type LaboratoryResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Phone     string          `json:"phone"`
	Address   AddressResponse `json:"address"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// AddressResponse represents the address in response body
type AddressResponse struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// ToLaboratoryResponse converts a domain laboratory to response DTO
func ToLaboratoryResponse(lab *laboratory.Laboratory) LaboratoryResponse {
	return LaboratoryResponse{
		ID:        lab.ID,
		Name:      lab.Name,
		Email:     lab.Email,
		Phone:     lab.Phone,
		Address: AddressResponse{
			Street:     lab.Address.Street,
			City:       lab.Address.City,
			State:      lab.Address.State,
			PostalCode: lab.Address.PostalCode,
			Country:    lab.Address.Country,
		},
		CreatedAt: lab.CreatedAt,
		UpdatedAt: lab.UpdatedAt,
	}
}

// ToLaboratoryResponseList converts a list of domain laboratories to response DTOs
func ToLaboratoryResponseList(labs []*laboratory.Laboratory) []LaboratoryResponse {
	responses := make([]LaboratoryResponse, len(labs))
	for i, lab := range labs {
		responses[i] = ToLaboratoryResponse(lab)
	}
	return responses
}

// ToAddress converts address request to domain address
func (r *AddressRequest) ToAddress() laboratory.Address {
	return laboratory.Address{
		Street:     r.Street,
		City:       r.City,
		State:      r.State,
		PostalCode: r.PostalCode,
		Country:    r.Country,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

