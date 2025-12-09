package main

import (
	"log"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/handler"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/router"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/outbound/persistence/memory"
	clientapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/client"
	labapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/laboratory"
	orderapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/order"
	prosthesisapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/prosthesis"
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
	idGen := uuid.NewGenerator()

	// Repositories
	labRepo := memory.NewLaboratoryRepository()
	clientRepo := memory.NewClientRepository()
	orderRepo := memory.NewOrderRepository()
	prosthesisRepo := memory.NewProsthesisRepository()

	// Services
	labService := labapp.NewService(labRepo, idGen)
	clientService := clientapp.NewService(clientRepo, labRepo, idGen)
	orderService := orderapp.NewService(orderRepo, clientRepo, idGen)
	prosthesisService := prosthesisapp.NewService(prosthesisRepo, labRepo, idGen)

	// Handlers
	labHandler := handler.NewLaboratoryHandler(labService)
	clientHandler := handler.NewClientHandler(clientService)
	orderHandler := handler.NewOrderHandler(orderService)
	prosthesisHandler := handler.NewProsthesisHandler(prosthesisService)

	// Initialize Clerk middleware (optional - only if configured)
	var clerkMiddleware *auth.ClerkMiddleware
	if cfg.Clerk.SecretKey != "" {
		log.Println("Clerk Secret Key configured, authentication enabled")
		clerkMiddleware = auth.NewClerkMiddleware(auth.ClerkConfig{
			SecretKey: cfg.Clerk.SecretKey,
		})
	} else {
		log.Println("Warning: Clerk Secret Key not configured, authentication disabled")
	}

	// Create router
	r := router.New(router.Config{
		LaboratoryHandler: labHandler,
		ClientHandler:     clientHandler,
		OrderHandler:      orderHandler,
		ProsthesisHandler: prosthesisHandler,
		ClerkMiddleware:   clerkMiddleware,
	})

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
