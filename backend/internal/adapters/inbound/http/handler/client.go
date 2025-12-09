package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	clientapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/client"
	domainerrors "github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/auth"
)

// ClientHandler handles HTTP requests for client operations
type ClientHandler struct {
	service *clientapp.Service
}

// NewClientHandler creates a new client handler
func NewClientHandler(service *clientapp.Service) *ClientHandler {
	return &ClientHandler{service: service}
}

// Create handles POST /api/v1/clients
func (h *ClientHandler) Create(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	var req dto.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	input := clientapp.CreateInput{
		LaboratoryID: laboratoryID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address.ToClientAddress(),
	}

	client, err := h.service.CreateClient(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToClientResponse(client))
}

// Get handles GET /api/v1/clients/:id
func (h *ClientHandler) Get(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "client id is required",
		})
		return
	}

	client, err := h.service.GetClient(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToClientResponse(client))
}

// Update handles PUT /api/v1/clients/:id
func (h *ClientHandler) Update(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "client id is required",
		})
		return
	}

	var req dto.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	input := clientapp.UpdateInput{
		ID:           id,
		LaboratoryID: laboratoryID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address.ToClientAddress(),
	}

	client, err := h.service.UpdateClient(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToClientResponse(client))
}

// List handles GET /api/v1/clients
func (h *ClientHandler) List(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	clients, err := h.service.ListClients(c.Request.Context(), laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToClientResponseList(clients))
}

// Delete handles DELETE /api/v1/clients/:id
func (h *ClientHandler) Delete(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "client id is required",
		})
		return
	}

	err := h.service.DeleteClient(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError converts domain errors to HTTP responses
func (h *ClientHandler) handleError(c *gin.Context, err error) {
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
			Error: "client not found",
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
