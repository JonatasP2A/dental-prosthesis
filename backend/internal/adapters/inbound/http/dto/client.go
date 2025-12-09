package dto

import (
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
)

// CreateClientRequest represents the request body for creating a client
type CreateClientRequest struct {
	Name    string                `json:"name" binding:"required"`
	Email   string                `json:"email" binding:"required,email"`
	Phone   string                `json:"phone" binding:"required"`
	Address ClientAddressRequest  `json:"address" binding:"required"`
}

// ClientAddressRequest represents the address in client request body
type ClientAddressRequest struct {
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	Country    string `json:"country" binding:"required"`
}

// UpdateClientRequest represents the request body for updating a client
type UpdateClientRequest struct {
	Name    string                `json:"name" binding:"required"`
	Email   string                `json:"email" binding:"required,email"`
	Phone   string                `json:"phone" binding:"required"`
	Address ClientAddressRequest  `json:"address" binding:"required"`
}

// ClientResponse represents the response body for a client
type ClientResponse struct {
	ID           string                `json:"id"`
	LaboratoryID string                `json:"laboratory_id"`
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	Phone        string                `json:"phone"`
	Address      ClientAddressResponse `json:"address"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// ClientAddressResponse represents the address in client response body
type ClientAddressResponse struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// ToClientResponse converts a domain client to response DTO
func ToClientResponse(c *client.Client) ClientResponse {
	return ClientResponse{
		ID:           c.ID,
		LaboratoryID: c.LaboratoryID,
		Name:         c.Name,
		Email:        c.Email,
		Phone:        c.Phone,
		Address: ClientAddressResponse{
			Street:     c.Address.Street,
			City:       c.Address.City,
			State:      c.Address.State,
			PostalCode: c.Address.PostalCode,
			Country:    c.Address.Country,
		},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToClientResponseList converts a list of domain clients to response DTOs
func ToClientResponseList(clients []*client.Client) []ClientResponse {
	responses := make([]ClientResponse, len(clients))
	for i, c := range clients {
		responses[i] = ToClientResponse(c)
	}
	return responses
}

// ToClientAddress converts client address request to domain address
func (r *ClientAddressRequest) ToClientAddress() client.Address {
	return client.Address{
		Street:     r.Street,
		City:       r.City,
		State:      r.State,
		PostalCode: r.PostalCode,
		Country:    r.Country,
	}
}
