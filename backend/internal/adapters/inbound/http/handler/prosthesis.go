package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	prosthesisapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/prosthesis"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
	domainerrors "github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

// ProsthesisHandler handles HTTP requests for prosthesis operations
type ProsthesisHandler struct {
	service *prosthesisapp.Service
}

// NewProsthesisHandler creates a new prosthesis handler
func NewProsthesisHandler(service *prosthesisapp.Service) *ProsthesisHandler {
	return &ProsthesisHandler{service: service}
}

// getLaboratoryID extracts laboratory_id from query parameter
func (h *ProsthesisHandler) getLaboratoryID(c *gin.Context) (string, error) {
	laboratoryID := c.Query("laboratory_id")
	if laboratoryID == "" {
		return "", errors.New("laboratory_id query parameter is required")
	}
	return laboratoryID, nil
}

// Create handles POST /api/v1/prostheses
func (h *ProsthesisHandler) Create(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req dto.CreateProsthesisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	prosthesisType := prosthesis.ProsthesisType(req.Type)
	if !prosthesisType.IsValid() {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid prosthesis type",
		})
		return
	}

	input := prosthesisapp.CreateInput{
		LaboratoryID:   laboratoryID,
		Type:           prosthesisType,
		Material:       req.Material,
		Shade:          req.Shade,
		Specifications: req.Specifications,
		Notes:          req.Notes,
	}

	p, err := h.service.CreateProsthesis(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToProsthesisResponse(p))
}

// Get handles GET /api/v1/prostheses/:id
func (h *ProsthesisHandler) Get(c *gin.Context) {
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
			Error: "prosthesis id is required",
		})
		return
	}

	p, err := h.service.GetProsthesis(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToProsthesisResponse(p))
}

// Update handles PUT /api/v1/prostheses/:id
func (h *ProsthesisHandler) Update(c *gin.Context) {
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
			Error: "prosthesis id is required",
		})
		return
	}

	var req dto.UpdateProsthesisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	prosthesisType := prosthesis.ProsthesisType(req.Type)
	if !prosthesisType.IsValid() {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid prosthesis type",
		})
		return
	}

	input := prosthesisapp.UpdateInput{
		ID:             id,
		LaboratoryID:   laboratoryID,
		Type:           prosthesisType,
		Material:       req.Material,
		Shade:          req.Shade,
		Specifications: req.Specifications,
		Notes:          req.Notes,
	}

	p, err := h.service.UpdateProsthesis(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToProsthesisResponse(p))
}

// List handles GET /api/v1/prostheses
func (h *ProsthesisHandler) List(c *gin.Context) {
	// Get laboratory ID from query parameter
	laboratoryID, err := h.getLaboratoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Parse optional filters
	var prosthesisType *prosthesis.ProsthesisType
	if typeParam := c.Query("type"); typeParam != "" {
		pt := prosthesis.ProsthesisType(typeParam)
		if !pt.IsValid() {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "invalid prosthesis type",
			})
			return
		}
		prosthesisType = &pt
	}

	var material *string
	if materialParam := c.Query("material"); materialParam != "" {
		material = &materialParam
	}

	prostheses, err := h.service.ListProstheses(c.Request.Context(), laboratoryID, prosthesisType, material)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToProsthesisResponseList(prostheses))
}

// Delete handles DELETE /api/v1/prostheses/:id
func (h *ProsthesisHandler) Delete(c *gin.Context) {
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
			Error: "prosthesis id is required",
		})
		return
	}

	err = h.service.DeleteProsthesis(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError converts domain errors to HTTP responses
func (h *ProsthesisHandler) handleError(c *gin.Context, err error) {
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
			Error: "prosthesis not found",
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
