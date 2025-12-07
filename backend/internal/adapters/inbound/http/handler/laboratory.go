package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	labapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/laboratory"
	domainerrors "github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

// LaboratoryHandler handles HTTP requests for laboratory operations
type LaboratoryHandler struct {
	service *labapp.Service
}

// NewLaboratoryHandler creates a new laboratory handler
func NewLaboratoryHandler(service *labapp.Service) *LaboratoryHandler {
	return &LaboratoryHandler{service: service}
}

// Create handles POST /api/v1/laboratories
func (h *LaboratoryHandler) Create(c *gin.Context) {
	var req dto.CreateLaboratoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	input := labapp.CreateInput{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address.ToAddress(),
	}

	lab, err := h.service.CreateLaboratory(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToLaboratoryResponse(lab))
}

// Get handles GET /api/v1/laboratories/:id
func (h *LaboratoryHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory id is required",
		})
		return
	}

	lab, err := h.service.GetLaboratory(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToLaboratoryResponse(lab))
}

// Update handles PUT /api/v1/laboratories/:id
func (h *LaboratoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory id is required",
		})
		return
	}

	var req dto.UpdateLaboratoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	input := labapp.UpdateInput{
		ID:      id,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address.ToAddress(),
	}

	lab, err := h.service.UpdateLaboratory(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToLaboratoryResponse(lab))
}

// List handles GET /api/v1/laboratories
func (h *LaboratoryHandler) List(c *gin.Context) {
	labs, err := h.service.ListLaboratories(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToLaboratoryResponseList(labs))
}

// Delete handles DELETE /api/v1/laboratories/:id
func (h *LaboratoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory id is required",
		})
		return
	}

	err := h.service.DeleteLaboratory(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError converts domain errors to HTTP responses
func (h *LaboratoryHandler) handleError(c *gin.Context, err error) {
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
			Error: "laboratory not found",
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

