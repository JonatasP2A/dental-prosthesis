package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	orderapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/order"
	domainerrors "github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/auth"
)

// OrderHandler handles HTTP requests for order operations
type OrderHandler struct {
	service *orderapp.Service
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(service *orderapp.Service) *OrderHandler {
	return &OrderHandler{service: service}
}

// Create handles POST /api/v1/orders
func (h *OrderHandler) Create(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	input := orderapp.CreateInput{
		ClientID:     req.ClientID,
		LaboratoryID: laboratoryID,
		Prosthesis:   dto.ToProsthesisItems(req.Prosthesis),
	}

	order, err := h.service.CreateOrder(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToOrderResponse(order))
}

// Get handles GET /api/v1/orders/:id
func (h *OrderHandler) Get(c *gin.Context) {
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
			Error: "order id is required",
		})
		return
	}

	order, err := h.service.GetOrder(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(order))
}

// Update handles PUT /api/v1/orders/:id
func (h *OrderHandler) Update(c *gin.Context) {
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
			Error: "order id is required",
		})
		return
	}

	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	input := orderapp.UpdateInput{
		ID:           id,
		LaboratoryID: laboratoryID,
		Prosthesis:   dto.ToProsthesisItems(req.Prosthesis),
	}

	order, err := h.service.UpdateOrder(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(order))
}

// UpdateStatus handles PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
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
			Error: "order id is required",
		})
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	// Validate status value
	if !order.IsValidStatus(req.Status) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid status value",
		})
		return
	}

	input := orderapp.UpdateStatusInput{
		ID:           id,
		LaboratoryID: laboratoryID,
		Status:       order.Status(req.Status),
	}

	o, err := h.service.UpdateOrderStatus(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(o))
}

// List handles GET /api/v1/orders
func (h *OrderHandler) List(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	orders, err := h.service.ListOrders(c.Request.Context(), laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponseList(orders))
}

// ListByClient handles GET /api/v1/clients/:id/orders
func (h *OrderHandler) ListByClient(c *gin.Context) {
	// Get laboratory ID from JWT claims
	laboratoryID := auth.GetLaboratoryID(c.Request.Context())
	if laboratoryID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "laboratory_id not found in token",
		})
		return
	}

	clientID := c.Param("id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "client id is required",
		})
		return
	}

	orders, err := h.service.ListOrdersByClient(c.Request.Context(), clientID, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponseList(orders))
}

// Delete handles DELETE /api/v1/orders/:id
func (h *OrderHandler) Delete(c *gin.Context) {
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
			Error: "order id is required",
		})
		return
	}

	err := h.service.DeleteOrder(c.Request.Context(), id, laboratoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError converts domain errors to HTTP responses
func (h *OrderHandler) handleError(c *gin.Context, err error) {
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
			Error: "order not found",
		})
	case errors.Is(err, domainerrors.ErrInvalidStatusTransition):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid status transition",
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
