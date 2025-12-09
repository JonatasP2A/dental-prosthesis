package memory

import (
	"context"
	"testing"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
)

func TestProsthesisRepository_Create(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p, err := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "A1", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}

	err = repo.Create(ctx, p)
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	// Verify it was stored
	stored, err := repo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if stored.ID != p.ID {
		t.Errorf("GetByID() ID = %v, want %v", stored.ID, p.ID)
	}
}

func TestProsthesisRepository_GetByID(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByID(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	p, err := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}
	_ = repo.Create(ctx, p)

	found, err := repo.GetByID(ctx, "prosthesis-123")
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if found.Type != p.Type {
		t.Errorf("GetByID() Type = %v, want %v", found.Type, p.Type)
	}
}

func TestProsthesisRepository_GetByID_SoftDeleted(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p, err := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}
	_ = repo.Create(ctx, p)

	// Delete the prosthesis
	_ = repo.Delete(ctx, p.ID)

	// Should not be found
	_, err = repo.GetByID(ctx, p.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() for deleted prosthesis error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestProsthesisRepository_Update(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p, err := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}
	_ = repo.Create(ctx, p)

	// Update
	p.Update(prosthesis.ProsthesisTypeBridge, "porcelain", "A2", "3-unit bridge", "Updated")
	err = repo.Update(ctx, p)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	// Verify update
	updated, err := repo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if updated.Type != prosthesis.ProsthesisTypeBridge {
		t.Errorf("Update() Type = %v, want %v", updated.Type, prosthesis.ProsthesisTypeBridge)
	}
	if updated.Material != "porcelain" {
		t.Errorf("Update() Material = %v, want %v", updated.Material, "porcelain")
	}
}

func TestProsthesisRepository_Update_NotFound(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p, err := prosthesis.NewProsthesis("non-existent", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}

	err = repo.Update(ctx, p)
	if err != errors.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestProsthesisRepository_Delete(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p, err := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}
	_ = repo.Create(ctx, p)

	err = repo.Delete(ctx, p.ID)
	if err != nil {
		t.Fatalf("Delete() unexpected error = %v", err)
	}

	// Verify soft delete
	_, err = repo.GetByID(ctx, p.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestProsthesisRepository_Delete_NotFound(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestProsthesisRepository_List(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	// Create prostheses for lab-123
	p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
	_ = repo.Create(ctx, p1)
	_ = repo.Create(ctx, p2)

	// Create prosthesis for different lab
	p3, _ := prosthesis.NewProsthesis("prosthesis-3", "lab-999", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	_ = repo.Create(ctx, p3)

	// List prostheses for lab-123
	prostheses, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}

	if len(prostheses) != 2 {
		t.Errorf("List() count = %v, want %v", len(prostheses), 2)
	}

	// Verify laboratory scoping
	for _, p := range prostheses {
		if p.LaboratoryID != "lab-123" {
			t.Errorf("List() LaboratoryID = %v, want %v", p.LaboratoryID, "lab-123")
		}
	}
}

func TestProsthesisRepository_List_Empty(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	prostheses, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}

	if len(prostheses) != 0 {
		t.Errorf("List() count = %v, want %v", len(prostheses), 0)
	}
}

func TestProsthesisRepository_List_ExcludesSoftDeleted(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
	_ = repo.Create(ctx, p1)
	_ = repo.Create(ctx, p2)

	// Delete one prosthesis
	_ = repo.Delete(ctx, p1.ID)

	// List should only return non-deleted
	prostheses, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}

	if len(prostheses) != 1 {
		t.Errorf("List() count = %v, want %v", len(prostheses), 1)
	}
	if prostheses[0].ID != p2.ID {
		t.Errorf("List() ID = %v, want %v", prostheses[0].ID, p2.ID)
	}
}

func TestProsthesisRepository_FindByType(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	// Create prostheses with different types
	p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
	p3, _ := prosthesis.NewProsthesis("prosthesis-3", "lab-123", prosthesis.ProsthesisTypeCrown, "metal", "", "", "")
	_ = repo.Create(ctx, p1)
	_ = repo.Create(ctx, p2)
	_ = repo.Create(ctx, p3)

	// Find by type
	prostheses, err := repo.FindByType(ctx, "lab-123", prosthesis.ProsthesisTypeCrown)
	if err != nil {
		t.Fatalf("FindByType() unexpected error = %v", err)
	}

	if len(prostheses) != 2 {
		t.Errorf("FindByType() count = %v, want %v", len(prostheses), 2)
	}

	for _, p := range prostheses {
		if p.Type != prosthesis.ProsthesisTypeCrown {
			t.Errorf("FindByType() Type = %v, want %v", p.Type, prosthesis.ProsthesisTypeCrown)
		}
		if p.LaboratoryID != "lab-123" {
			t.Errorf("FindByType() LaboratoryID = %v, want %v", p.LaboratoryID, "lab-123")
		}
	}
}

func TestProsthesisRepository_FindByMaterial(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	// Create prostheses with different materials
	p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
	p3, _ := prosthesis.NewProsthesis("prosthesis-3", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
	_ = repo.Create(ctx, p1)
	_ = repo.Create(ctx, p2)
	_ = repo.Create(ctx, p3)

	// Find by material
	prostheses, err := repo.FindByMaterial(ctx, "lab-123", "zirconia")
	if err != nil {
		t.Fatalf("FindByMaterial() unexpected error = %v", err)
	}

	if len(prostheses) != 2 {
		t.Errorf("FindByMaterial() count = %v, want %v", len(prostheses), 2)
	}

	for _, p := range prostheses {
		if p.Material != "zirconia" {
			t.Errorf("FindByMaterial() Material = %v, want %v", p.Material, "zirconia")
		}
		if p.LaboratoryID != "lab-123" {
			t.Errorf("FindByMaterial() LaboratoryID = %v, want %v", p.LaboratoryID, "lab-123")
		}
	}
}

func TestProsthesisRepository_FindByMaterial_CaseInsensitive(t *testing.T) {
	repo := NewProsthesisRepository()
	ctx := context.Background()

	p, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "Zirconia", "", "", "")
	_ = repo.Create(ctx, p)

	// Find with lowercase
	prostheses, err := repo.FindByMaterial(ctx, "lab-123", "zirconia")
	if err != nil {
		t.Fatalf("FindByMaterial() unexpected error = %v", err)
	}

	if len(prostheses) != 1 {
		t.Errorf("FindByMaterial() count = %v, want %v", len(prostheses), 1)
	}
}
