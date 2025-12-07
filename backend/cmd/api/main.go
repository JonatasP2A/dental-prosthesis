package main

import (
	"log"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/handler"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/router"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/outbound/persistence/memory"
	labapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/config"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/auth"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/uuid"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize dependencies
	labRepo := memory.NewLaboratoryRepository()
	idGen := uuid.NewGenerator()
	labService := labapp.NewService(labRepo, idGen)
	labHandler := handler.NewLaboratoryHandler(labService)

	// Initialize Clerk middleware (optional - only if configured)
	var clerkMiddleware *auth.ClerkMiddleware
	if cfg.Clerk.JWKSURL != "" {
		clerkMiddleware = auth.NewClerkMiddleware(auth.ClerkConfig{
			JWKSURL: cfg.Clerk.JWKSURL,
		})
	} else {
		log.Println("Warning: Clerk JWKS URL not configured, authentication disabled")
	}

	// Create router
	r := router.New(router.Config{
		LaboratoryHandler: labHandler,
		ClerkMiddleware:   clerkMiddleware,
	})

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
