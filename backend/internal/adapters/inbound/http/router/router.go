package router

import (
	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/handler"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/auth"
)

// Config holds router configuration
type Config struct {
	LaboratoryHandler *handler.LaboratoryHandler
	ClientHandler     *handler.ClientHandler
	OrderHandler      *handler.OrderHandler
	ProsthesisHandler *handler.ProsthesisHandler
	TechnicianHandler *handler.TechnicianHandler
	ClerkMiddleware   *auth.ClerkMiddleware
}

// New creates a new Gin router with all routes configured
func New(cfg Config) *gin.Engine {
	r := gin.Default()

	// Health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Laboratory routes (protected)
	laboratories := v1.Group("/laboratories")
	if cfg.ClerkMiddleware != nil {
		laboratories.Use(cfg.ClerkMiddleware.Authenticate())
	}
	{
		laboratories.POST("", cfg.LaboratoryHandler.Create)
		laboratories.GET("", cfg.LaboratoryHandler.List)
		laboratories.GET("/:id", cfg.LaboratoryHandler.Get)
		laboratories.PUT("/:id", cfg.LaboratoryHandler.Update)
		laboratories.DELETE("/:id", cfg.LaboratoryHandler.Delete)
	}

	// Client routes (protected)
	if cfg.ClientHandler != nil {
		clients := v1.Group("/clients")
		if cfg.ClerkMiddleware != nil {
			clients.Use(cfg.ClerkMiddleware.Authenticate())
		}
		{
			clients.POST("", cfg.ClientHandler.Create)
			clients.GET("", cfg.ClientHandler.List)
			clients.GET("/:id", cfg.ClientHandler.Get)
			clients.PUT("/:id", cfg.ClientHandler.Update)
			clients.DELETE("/:id", cfg.ClientHandler.Delete)
		}

		// Nested route: GET /api/v1/clients/:id/orders
		if cfg.OrderHandler != nil {
			clients.GET("/:id/orders", cfg.OrderHandler.ListByClient)
		}
	}

	// Order routes (protected)
	if cfg.OrderHandler != nil {
		orders := v1.Group("/orders")
		if cfg.ClerkMiddleware != nil {
			orders.Use(cfg.ClerkMiddleware.Authenticate())
		}
		{
			orders.POST("", cfg.OrderHandler.Create)
			orders.GET("", cfg.OrderHandler.List)
			orders.GET("/:id", cfg.OrderHandler.Get)
			orders.PUT("/:id", cfg.OrderHandler.Update)
			orders.PATCH("/:id/status", cfg.OrderHandler.UpdateStatus)
			orders.DELETE("/:id", cfg.OrderHandler.Delete)
		}
	}

	// Prosthesis routes (protected)
	if cfg.ProsthesisHandler != nil {
		prostheses := v1.Group("/prostheses")
		if cfg.ClerkMiddleware != nil {
			prostheses.Use(cfg.ClerkMiddleware.Authenticate())
		}
		{
			prostheses.POST("", cfg.ProsthesisHandler.Create)
			prostheses.GET("", cfg.ProsthesisHandler.List)
			prostheses.GET("/:id", cfg.ProsthesisHandler.Get)
			prostheses.PUT("/:id", cfg.ProsthesisHandler.Update)
			prostheses.DELETE("/:id", cfg.ProsthesisHandler.Delete)
		}
	}

	// Technician routes (protected)
	if cfg.TechnicianHandler != nil {
		technicians := v1.Group("/technicians")
		if cfg.ClerkMiddleware != nil {
			technicians.Use(cfg.ClerkMiddleware.Authenticate())
		}
		{
			technicians.POST("", cfg.TechnicianHandler.Create)
			technicians.GET("", cfg.TechnicianHandler.List)
			technicians.GET("/:id", cfg.TechnicianHandler.Get)
			technicians.PUT("/:id", cfg.TechnicianHandler.Update)
			technicians.DELETE("/:id", cfg.TechnicianHandler.Delete)
		}
	}

	return r
}
