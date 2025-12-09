package dto

import (
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
)

// CreateProsthesisRequest represents the request body for creating a prosthesis
type CreateProsthesisRequest struct {
	Type           string `json:"type" binding:"required"`
	Material       string `json:"material" binding:"required"`
	Shade          string `json:"shade"`
	Specifications string `json:"specifications"`
	Notes          string `json:"notes"`
}

// UpdateProsthesisRequest represents the request body for updating a prosthesis
type UpdateProsthesisRequest struct {
	Type           string `json:"type" binding:"required"`
	Material       string `json:"material" binding:"required"`
	Shade          string `json:"shade"`
	Specifications string `json:"specifications"`
	Notes          string `json:"notes"`
}

// ProsthesisResponse represents the response body for a prosthesis
type ProsthesisResponse struct {
	ID             string    `json:"id"`
	LaboratoryID   string    `json:"laboratory_id"`
	Type           string    `json:"type"`
	Material       string    `json:"material"`
	Shade          string    `json:"shade"`
	Specifications string    `json:"specifications"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ToProsthesisResponse converts a domain prosthesis to response DTO
func ToProsthesisResponse(p *prosthesis.Prosthesis) ProsthesisResponse {
	return ProsthesisResponse{
		ID:             p.ID,
		LaboratoryID:   p.LaboratoryID,
		Type:           string(p.Type),
		Material:       p.Material,
		Shade:          p.Shade,
		Specifications: p.Specifications,
		Notes:          p.Notes,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

// ToProsthesisResponseList converts a list of domain prostheses to response DTOs
func ToProsthesisResponseList(prostheses []*prosthesis.Prosthesis) []ProsthesisResponse {
	responses := make([]ProsthesisResponse, len(prostheses))
	for i, p := range prostheses {
		responses[i] = ToProsthesisResponse(p)
	}
	return responses
}
