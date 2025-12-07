package memory

import (
	"context"
	"testing"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

func TestLaboratoryRepository_Create(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	lab := &laboratory.Laboratory{
		ID:    "lab-123",
		Name:  "Test Lab",
		Email: "test@lab.com",
		Phone: "+5511999999999",
		Address: laboratory.Address{
			Street:     "Test Street",
			City:       "Test City",
			State:      "SP",
			PostalCode: "01234-567",
			Country:    "Brazil",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.Create(ctx, lab)
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	// Verify it was stored
	stored, err := repo.GetByID(ctx, lab.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if stored.ID != lab.ID {
		t.Errorf("GetByID() ID = %v, want %v", stored.ID, lab.ID)
	}
}

func TestLaboratoryRepository_GetByID(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByID(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	lab := &laboratory.Laboratory{
		ID:        "lab-123",
		Name:      "Test Lab",
		Email:     "test@lab.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, lab)

	found, err := repo.GetByID(ctx, "lab-123")
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if found.Name != lab.Name {
		t.Errorf("GetByID() Name = %v, want %v", found.Name, lab.Name)
	}
}

func TestLaboratoryRepository_GetByID_SoftDeleted(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	lab := &laboratory.Laboratory{
		ID:        "lab-123",
		Name:      "Test Lab",
		Email:     "test@lab.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, lab)

	// Delete the lab
	_ = repo.Delete(ctx, lab.ID)

	// Should not be found
	_, err := repo.GetByID(ctx, lab.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() for deleted lab error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestLaboratoryRepository_GetByEmail(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByEmail(ctx, "nonexistent@lab.com")
	if err != errors.ErrNotFound {
		t.Errorf("GetByEmail() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	lab := &laboratory.Laboratory{
		ID:        "lab-123",
		Name:      "Test Lab",
		Email:     "test@lab.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, lab)

	found, err := repo.GetByEmail(ctx, "test@lab.com")
	if err != nil {
		t.Fatalf("GetByEmail() unexpected error = %v", err)
	}

	if found.ID != lab.ID {
		t.Errorf("GetByEmail() ID = %v, want %v", found.ID, lab.ID)
	}
}

func TestLaboratoryRepository_Update(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	lab := &laboratory.Laboratory{
		ID:        "lab-123",
		Name:      "Test Lab",
		Email:     "test@lab.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, lab)

	// Update
	lab.Name = "Updated Lab"
	err := repo.Update(ctx, lab)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	// Verify update
	found, _ := repo.GetByID(ctx, lab.ID)
	if found.Name != "Updated Lab" {
		t.Errorf("Update() Name = %v, want %v", found.Name, "Updated Lab")
	}
}

func TestLaboratoryRepository_Update_NotFound(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	lab := &laboratory.Laboratory{
		ID:   "non-existent",
		Name: "Test Lab",
	}

	err := repo.Update(ctx, lab)
	if err != errors.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestLaboratoryRepository_Delete(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	lab := &laboratory.Laboratory{
		ID:        "lab-123",
		Name:      "Test Lab",
		Email:     "test@lab.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, lab)

	err := repo.Delete(ctx, lab.ID)
	if err != nil {
		t.Fatalf("Delete() unexpected error = %v", err)
	}

	// Should not be found anymore
	_, err = repo.GetByID(ctx, lab.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestLaboratoryRepository_Delete_NotFound(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestLaboratoryRepository_List(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	// Empty list
	labs, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(labs) != 0 {
		t.Errorf("List() got %d labs, want 0", len(labs))
	}

	// Add labs
	_ = repo.Create(ctx, &laboratory.Laboratory{
		ID:        "lab-1",
		Name:      "Lab 1",
		Email:     "lab1@test.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &laboratory.Laboratory{
		ID:        "lab-2",
		Name:      "Lab 2",
		Email:     "lab2@test.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	labs, err = repo.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(labs) != 2 {
		t.Errorf("List() got %d labs, want 2", len(labs))
	}
}

func TestLaboratoryRepository_List_ExcludesDeleted(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	_ = repo.Create(ctx, &laboratory.Laboratory{
		ID:        "lab-1",
		Name:      "Lab 1",
		Email:     "lab1@test.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &laboratory.Laboratory{
		ID:        "lab-2",
		Name:      "Lab 2",
		Email:     "lab2@test.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	// Delete one
	_ = repo.Delete(ctx, "lab-1")

	labs, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(labs) != 1 {
		t.Errorf("List() got %d labs, want 1", len(labs))
	}
	if labs[0].ID != "lab-2" {
		t.Errorf("List() remaining lab ID = %v, want lab-2", labs[0].ID)
	}
}

func TestLaboratoryRepository_ImmutableClone(t *testing.T) {
	repo := NewLaboratoryRepository()
	ctx := context.Background()

	lab := &laboratory.Laboratory{
		ID:        "lab-123",
		Name:      "Original",
		Email:     "test@lab.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, lab)

	// Get and modify
	found, _ := repo.GetByID(ctx, lab.ID)
	found.Name = "Modified"

	// Original should be unchanged
	original, _ := repo.GetByID(ctx, lab.ID)
	if original.Name != "Original" {
		t.Errorf("Repository data was mutated, got %v, want Original", original.Name)
	}
}

