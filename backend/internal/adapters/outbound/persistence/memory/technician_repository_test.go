package memory

import (
	"context"
	"testing"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
)

func TestTechnicianRepository_Create(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err := repo.Create(ctx, tech)
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	// Verify it was stored
	stored, err := repo.GetByID(ctx, tech.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if stored.ID != tech.ID {
		t.Errorf("GetByID() ID = %v, want %v", stored.ID, tech.ID)
	}
}

func TestTechnicianRepository_GetByID(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByID(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, tech)

	found, err := repo.GetByID(ctx, "tech-123")
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if found.Name != tech.Name {
		t.Errorf("GetByID() Name = %v, want %v", found.Name, tech.Name)
	}
}

func TestTechnicianRepository_GetByID_SoftDeleted(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, tech)

	// Delete the technician
	_ = repo.Delete(ctx, tech.ID)

	// Should not be found
	_, err := repo.GetByID(ctx, tech.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() for deleted tech error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestTechnicianRepository_GetByEmail(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByEmail(ctx, "lab-123", "nonexistent@lab.com")
	if err != errors.ErrNotFound {
		t.Errorf("GetByEmail() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, tech)

	found, err := repo.GetByEmail(ctx, "lab-123", "john@lab.com")
	if err != nil {
		t.Fatalf("GetByEmail() unexpected error = %v", err)
	}

	if found.ID != tech.ID {
		t.Errorf("GetByEmail() ID = %v, want %v", found.ID, tech.ID)
	}

	// Test email exists in different laboratory
	_, err = repo.GetByEmail(ctx, "lab-999", "john@lab.com")
	if err != errors.ErrNotFound {
		t.Errorf("GetByEmail() should not find email from different lab, error = %v", err)
	}
}

func TestTechnicianRepository_Update(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, tech)

	// Update
	tech.Name = "Updated Name"
	tech.Role = technician.RoleSeniorTechnician
	err := repo.Update(ctx, tech)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	// Verify update
	found, _ := repo.GetByID(ctx, tech.ID)
	if found.Name != "Updated Name" {
		t.Errorf("Update() Name = %v, want %v", found.Name, "Updated Name")
	}
	if found.Role != technician.RoleSeniorTechnician {
		t.Errorf("Update() Role = %v, want %v", found.Role, technician.RoleSeniorTechnician)
	}
}

func TestTechnicianRepository_Update_NotFound(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	tech := &technician.Technician{
		ID:   "non-existent",
		Name: "John Doe",
	}

	err := repo.Update(ctx, tech)
	if err != errors.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestTechnicianRepository_Delete(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, tech)

	err := repo.Delete(ctx, tech.ID)
	if err != nil {
		t.Fatalf("Delete() unexpected error = %v", err)
	}

	// Should not be found anymore
	_, err = repo.GetByID(ctx, tech.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestTechnicianRepository_List(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	// Empty list
	techs, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(techs) != 0 {
		t.Errorf("List() got %d techs, want 0", len(techs))
	}

	// Add technicians
	_ = repo.Create(ctx, &technician.Technician{
		ID:           "tech-1",
		LaboratoryID: "lab-123",
		Name:         "Tech 1",
		Email:        "tech1@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	_ = repo.Create(ctx, &technician.Technician{
		ID:           "tech-2",
		LaboratoryID: "lab-123",
		Name:         "Tech 2",
		Email:        "tech2@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	_ = repo.Create(ctx, &technician.Technician{
		ID:           "tech-3",
		LaboratoryID: "lab-999",
		Name:         "Tech 3",
		Email:        "tech3@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})

	techs, err = repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(techs) != 2 {
		t.Errorf("List() got %d techs, want 2", len(techs))
	}
}

func TestTechnicianRepository_ListByRole(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	_ = repo.Create(ctx, &technician.Technician{
		ID:           "tech-1",
		LaboratoryID: "lab-123",
		Name:         "Tech 1",
		Email:        "tech1@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	_ = repo.Create(ctx, &technician.Technician{
		ID:           "tech-2",
		LaboratoryID: "lab-123",
		Name:         "Senior Tech",
		Email:        "senior@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleSeniorTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})

	techs, err := repo.ListByRole(ctx, "lab-123", technician.RoleSeniorTechnician)
	if err != nil {
		t.Fatalf("ListByRole() unexpected error = %v", err)
	}
	if len(techs) != 1 {
		t.Errorf("ListByRole() got %d techs, want 1", len(techs))
	}
	if techs[0].Role != technician.RoleSeniorTechnician {
		t.Errorf("ListByRole() Role = %v, want %v", techs[0].Role, technician.RoleSeniorTechnician)
	}
}

func TestTechnicianRepository_ImmutableClone(t *testing.T) {
	repo := NewTechnicianRepository()
	ctx := context.Background()

	tech := &technician.Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "Original",
		Email:        "original@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(ctx, tech)

	// Get and modify
	found, _ := repo.GetByID(ctx, tech.ID)
	found.Name = "Modified"

	// Original should be unchanged
	original, _ := repo.GetByID(ctx, tech.ID)
	if original.Name != "Original" {
		t.Errorf("Repository data was mutated, got %v, want Original", original.Name)
	}
}
