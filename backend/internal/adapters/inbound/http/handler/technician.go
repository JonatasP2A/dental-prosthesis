package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	techapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/technician"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
	domainerrors "github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

// TechnicianHandler handles HTTP requests for technician operations
type TechnicianHandler struct {
	service *techapp.Service
}

// NewTechnicianHandler creates a new technician handler
func NewTechnicianHandler(service *techapp.Service) *TechnicianHandler {
	return &TechnicianHandler{service: service}
}

// getLaboratoryID extracts laboratory_id from query parameter
func (h *TechnicianHandler) getLaboratoryID(c *gin.Context) (string, error) {
	laboratoryID := c.Query("laboratory_id")
	if laboratoryID == "" {
		return "", errors.New("laboratory_id query parameter is required")
	}
	return laboratoryID, nil
}

// Create handles POST /api/v1/technicians?laboratory_id=xxx
func (h *TechnicianHandler) Create(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req dto.CreateTechnicianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	// Convert role string to Role enum
	role, err := dto.ToRole(req.Role)
	if err != nil || !role.IsValid() {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid role. Must be one of: senior_technician, technician, apprentice",
		})
		return
	}

	input := techapp.CreateInput{
		LaboratoryID:    laboratoryID,
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		Role:            role,
		Specializations: req.Specializations,
	}

	tech, err := h.service.CreateTechnician(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToTechnicianResponse(tech))
}

// Get handles GET /api/v1/technicians/:id?laboratory_id=xxx
func (h *TechnicianHandler) Get(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "technician id is required",
		})
		return
	}

	tech, err := h.service.GetTechnician(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTechnicianResponse(tech))
}

// Update handles PUT /api/v1/technicians/:id?laboratory_id=xxx
func (h *TechnicianHandler) Update(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "technician id is required",
		})
		return
	}

	var req dto.UpdateTechnicianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	// Convert role string to Role enum
	role, err := dto.ToRole(req.Role)
	if err != nil || !role.IsValid() {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid role. Must be one of: senior_technician, technician, apprentice",
		})
		return
	}

	input := techapp.UpdateInput{
		ID:              id,
		LaboratoryID:    laboratoryID,
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		Role:            role,
		Specializations: req.Specializations,
	}

	tech, err := h.service.UpdateTechnician(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTechnicianResponse(tech))
}

// List handles GET /api/v1/technicians?laboratory_id=xxx&role=xxx
func (h *TechnicianHandler) List(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Get optional role filter
	var roleFilter *technician.Role
	roleStr := c.Query("role")
	if roleStr != "" {
		role, err := dto.ToRole(roleStr)
		if err != nil || !role.IsValid() {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "invalid role filter. Must be one of: senior_technician, technician, apprentice",
			})
			return
		}
		roleFilter = &role
	}

	techs, err := h.service.ListTechnicians(c.Request.Context(), laboratoryID, roleFilter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTechnicianResponseList(techs))
}

// Delete handles DELETE /api/v1/technicians/:id?laboratory_id=xxx
func (h *TechnicianHandler) Delete(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "technician id is required",
		})
		return
	}

	err = h.service.DeleteTechnician(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError converts domain errors to HTTP responses
func (h *TechnicianHandler) handleError(c *gin.Context, err error) {
	var validationErrors domainerrors.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make(map[string]string)
		for _, ve := range validationErrors {
			details[ve.Field] = ve.Message
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation failed",
			Details: details,
		})
		return
	}

	switch {
	case errors.Is(err, domainerrors.ErrNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "technician not found",
		})
	case errors.Is(err, domainerrors.ErrDuplicateEmail):
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error: "email already exists",
		})
	case errors.Is(err, domainerrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "unauthorized",
		})
	case errors.Is(err, domainerrors.ErrForbidden):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "forbidden",
		})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal server error",
		})
	}
}
