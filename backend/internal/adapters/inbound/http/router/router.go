package router

import (
	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/handler"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/auth"
)

// Config holds router configuration
type Config struct {
	LaboratoryHandler *handler.LaboratoryHandler
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

	return r
}

