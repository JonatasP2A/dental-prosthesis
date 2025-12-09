package dto

import (
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
)

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	ClientID   string                   `json:"client_id" binding:"required"`
	Prosthesis []ProsthesisItemRequest  `json:"prosthesis" binding:"required,dive"`
}

// ProsthesisItemRequest represents a prosthesis item in the request body
type ProsthesisItemRequest struct {
	Type     string `json:"type" binding:"required"`
	Material string `json:"material" binding:"required"`
	Shade    string `json:"shade"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Notes    string `json:"notes"`
}

// UpdateOrderRequest represents the request body for updating an order
type UpdateOrderRequest struct {
	Prosthesis []ProsthesisItemRequest `json:"prosthesis" binding:"required,dive"`
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// OrderResponse represents the response body for an order
type OrderResponse struct {
	ID           string                    `json:"id"`
	ClientID     string                    `json:"client_id"`
	LaboratoryID string                    `json:"laboratory_id"`
	Status       string                    `json:"status"`
	Prosthesis   []ProsthesisItemResponse  `json:"prosthesis"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// ProsthesisItemResponse represents a prosthesis item in the response body
type ProsthesisItemResponse struct {
	Type     string `json:"type"`
	Material string `json:"material"`
	Shade    string `json:"shade,omitempty"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes,omitempty"`
}

// ToOrderResponse converts a domain order to response DTO
func ToOrderResponse(o *order.Order) OrderResponse {
	prosthesisResponses := make([]ProsthesisItemResponse, len(o.Prosthesis))
	for i, p := range o.Prosthesis {
		prosthesisResponses[i] = ProsthesisItemResponse{
			Type:     p.Type,
			Material: p.Material,
			Shade:    p.Shade,
			Quantity: p.Quantity,
			Notes:    p.Notes,
		}
	}

	return OrderResponse{
		ID:           o.ID,
		ClientID:     o.ClientID,
		LaboratoryID: o.LaboratoryID,
		Status:       string(o.Status),
		Prosthesis:   prosthesisResponses,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

// ToOrderResponseList converts a list of domain orders to response DTOs
func ToOrderResponseList(orders []*order.Order) []OrderResponse {
	responses := make([]OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = ToOrderResponse(o)
	}
	return responses
}

// ToProsthesisItems converts prosthesis item requests to domain prosthesis items
func ToProsthesisItems(items []ProsthesisItemRequest) []order.ProsthesisItem {
	result := make([]order.ProsthesisItem, len(items))
	for i, item := range items {
		result[i] = order.ProsthesisItem{
			Type:     item.Type,
			Material: item.Material,
			Shade:    item.Shade,
			Quantity: item.Quantity,
			Notes:    item.Notes,
		}
	}
	return result
}
