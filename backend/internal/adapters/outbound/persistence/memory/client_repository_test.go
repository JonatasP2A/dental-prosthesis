package memory

import (
	"context"
	"testing"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

func TestClientRepository_Create(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		Email:        "test@example.com",
		Phone:        "+5511999999999",
		Address: client.Address{
			Street:     "Test Street",
			City:       "Test City",
			State:      "SP",
			PostalCode: "01234-567",
			Country:    "Brazil",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.Create(ctx, c)
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	// Verify it was stored
	stored, err := repo.GetByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if stored.ID != c.ID {
		t.Errorf("GetByID() ID = %v, want %v", stored.ID, c.ID)
	}
}

func TestClientRepository_GetByID(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByID(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		Email:        "test@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c)

	found, err := repo.GetByID(ctx, "client-123")
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if found.Name != c.Name {
		t.Errorf("GetByID() Name = %v, want %v", found.Name, c.Name)
	}
}

func TestClientRepository_GetByID_SoftDeleted(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		Email:        "test@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c)

	// Delete the client
	_ = repo.Delete(ctx, c.ID)

	// Should not be found
	_, err := repo.GetByID(ctx, c.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() for deleted client error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestClientRepository_GetByEmail(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByEmail(ctx, "lab-123", "nonexistent@example.com")
	if err != errors.ErrNotFound {
		t.Errorf("GetByEmail() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		Email:        "test@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c)

	found, err := repo.GetByEmail(ctx, "lab-123", "test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() unexpected error = %v", err)
	}

	if found.ID != c.ID {
		t.Errorf("GetByEmail() ID = %v, want %v", found.ID, c.ID)
	}
}

func TestClientRepository_GetByEmail_LaboratoryScoped(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	// Create clients with same email in different laboratories
	c1 := &client.Client{
		ID:           "client-1",
		LaboratoryID: "lab-123",
		Name:         "Client 1",
		Email:        "same@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c1)

	c2 := &client.Client{
		ID:           "client-2",
		LaboratoryID: "lab-456",
		Name:         "Client 2",
		Email:        "same@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c2)

	// Should find client in lab-123
	found, err := repo.GetByEmail(ctx, "lab-123", "same@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() unexpected error = %v", err)
	}
	if found.ID != "client-1" {
		t.Errorf("GetByEmail() ID = %v, want client-1", found.ID)
	}

	// Should find client in lab-456
	found, err = repo.GetByEmail(ctx, "lab-456", "same@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() unexpected error = %v", err)
	}
	if found.ID != "client-2" {
		t.Errorf("GetByEmail() ID = %v, want client-2", found.ID)
	}
}

func TestClientRepository_Update(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		Email:        "test@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c)

	// Update
	c.Name = "Updated Client"
	err := repo.Update(ctx, c)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	// Verify update
	found, _ := repo.GetByID(ctx, c.ID)
	if found.Name != "Updated Client" {
		t.Errorf("Update() Name = %v, want %v", found.Name, "Updated Client")
	}
}

func TestClientRepository_Update_NotFound(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	c := &client.Client{
		ID:   "non-existent",
		Name: "Test Client",
	}

	err := repo.Update(ctx, c)
	if err != errors.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestClientRepository_Delete(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		Email:        "test@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c)

	err := repo.Delete(ctx, c.ID)
	if err != nil {
		t.Fatalf("Delete() unexpected error = %v", err)
	}

	// Should not be found anymore
	_, err = repo.GetByID(ctx, c.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestClientRepository_Delete_NotFound(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestClientRepository_List(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	// Empty list
	clients, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(clients) != 0 {
		t.Errorf("List() got %d clients, want 0", len(clients))
	}

	// Add clients
	_ = repo.Create(ctx, &client.Client{
		ID:           "client-1",
		LaboratoryID: "lab-123",
		Name:         "Client 1",
		Email:        "client1@test.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	_ = repo.Create(ctx, &client.Client{
		ID:           "client-2",
		LaboratoryID: "lab-123",
		Name:         "Client 2",
		Email:        "client2@test.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	_ = repo.Create(ctx, &client.Client{
		ID:           "client-3",
		LaboratoryID: "lab-456", // Different laboratory
		Name:         "Client 3",
		Email:        "client3@test.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})

	clients, err = repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(clients) != 2 {
		t.Errorf("List() got %d clients, want 2", len(clients))
	}
}

func TestClientRepository_List_ExcludesDeleted(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	_ = repo.Create(ctx, &client.Client{
		ID:           "client-1",
		LaboratoryID: "lab-123",
		Name:         "Client 1",
		Email:        "client1@test.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	_ = repo.Create(ctx, &client.Client{
		ID:           "client-2",
		LaboratoryID: "lab-123",
		Name:         "Client 2",
		Email:        "client2@test.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})

	// Delete one
	_ = repo.Delete(ctx, "client-1")

	clients, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(clients) != 1 {
		t.Errorf("List() got %d clients, want 1", len(clients))
	}
	if clients[0].ID != "client-2" {
		t.Errorf("List() remaining client ID = %v, want client-2", clients[0].ID)
	}
}

func TestClientRepository_ImmutableClone(t *testing.T) {
	repo := NewClientRepository()
	ctx := context.Background()

	c := &client.Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Original",
		Email:        "test@example.com",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, c)

	// Get and modify
	found, _ := repo.GetByID(ctx, c.ID)
	found.Name = "Modified"

	// Original should be unchanged
	original, _ := repo.GetByID(ctx, c.ID)
	if original.Name != "Original" {
		t.Errorf("Repository data was mutated, got %v, want Original", original.Name)
	}
}
