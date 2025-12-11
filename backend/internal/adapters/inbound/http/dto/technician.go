package dto

import (
	"errors"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
)

// CreateTechnicianRequest represents the request body for creating a technician
type CreateTechnicianRequest struct {
	Name            string   `json:"name" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	Phone           string   `json:"phone" binding:"required"`
	Role            string   `json:"role" binding:"required"`
	Specializations []string `json:"specializations,omitempty"`
}

// UpdateTechnicianRequest represents the request body for updating a technician
type UpdateTechnicianRequest struct {
	Name            string   `json:"name" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	Phone           string   `json:"phone" binding:"required"`
	Role            string   `json:"role" binding:"required"`
	Specializations []string `json:"specializations,omitempty"`
}

// TechnicianResponse represents the response body for a technician
type TechnicianResponse struct {
	ID              string    `json:"id"`
	LaboratoryID    string    `json:"laboratory_id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Role            string    `json:"role"`
	Specializations []string  `json:"specializations"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToTechnicianResponse converts a domain technician to response DTO
func ToTechnicianResponse(tech *technician.Technician) TechnicianResponse {
	return TechnicianResponse{
		ID:              tech.ID,
		LaboratoryID:    tech.LaboratoryID,
		Name:            tech.Name,
		Email:           tech.Email,
		Phone:           tech.Phone,
		Role:            tech.Role.String(),
		Specializations: tech.Specializations,
		CreatedAt:       tech.CreatedAt,
		UpdatedAt:       tech.UpdatedAt,
	}
}

// ToTechnicianResponseList converts a list of domain technicians to response DTOs
func ToTechnicianResponseList(techs []*technician.Technician) []TechnicianResponse {
	responses := make([]TechnicianResponse, len(techs))
	for i, tech := range techs {
		responses[i] = ToTechnicianResponse(tech)
	}
	return responses
}

// ToRole converts a string to a technician Role
func ToRole(roleStr string) (technician.Role, error) {
	role := technician.Role(roleStr)
	if !role.IsValid() {
		return technician.Role(""), errors.New("invalid role")
	}
	return role, nil
}
