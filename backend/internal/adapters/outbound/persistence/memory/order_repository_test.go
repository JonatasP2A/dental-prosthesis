package memory

import (
	"context"
	"testing"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
)

func TestOrderRepository_Create(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.Create(ctx, o)
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	// Verify it was stored
	stored, err := repo.GetByID(ctx, o.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if stored.ID != o.ID {
		t.Errorf("GetByID() ID = %v, want %v", stored.ID, o.ID)
	}
}

func TestOrderRepository_GetByID(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByID(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, errors.ErrNotFound)
	}

	// Create and get
	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, o)

	found, err := repo.GetByID(ctx, "order-123")
	if err != nil {
		t.Fatalf("GetByID() unexpected error = %v", err)
	}

	if found.Status != o.Status {
		t.Errorf("GetByID() Status = %v, want %v", found.Status, o.Status)
	}
}

func TestOrderRepository_GetByID_SoftDeleted(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, o)

	// Delete the order
	_ = repo.Delete(ctx, o.ID)

	// Should not be found
	_, err := repo.GetByID(ctx, o.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() for deleted order error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestOrderRepository_Update(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, o)

	// Update
	o.Prosthesis = []order.ProsthesisItem{
		{
			Type:     "bridge",
			Material: "porcelain",
			Quantity: 2,
		},
	}
	err := repo.Update(ctx, o)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	// Verify update
	found, _ := repo.GetByID(ctx, o.ID)
	if len(found.Prosthesis) != 1 {
		t.Errorf("Update() Prosthesis length = %v, want 1", len(found.Prosthesis))
	}
	if found.Prosthesis[0].Type != "bridge" {
		t.Errorf("Update() Prosthesis[0].Type = %v, want bridge", found.Prosthesis[0].Type)
	}
}

func TestOrderRepository_Update_NotFound(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:   "non-existent",
		Status: order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
	}

	err := repo.Update(ctx, o)
	if err != errors.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, o)

	err := repo.UpdateStatus(ctx, "order-123", order.StatusInProduction)
	if err != nil {
		t.Fatalf("UpdateStatus() unexpected error = %v", err)
	}

	// Verify update
	found, _ := repo.GetByID(ctx, o.ID)
	if found.Status != order.StatusInProduction {
		t.Errorf("UpdateStatus() Status = %v, want %v", found.Status, order.StatusInProduction)
	}
}

func TestOrderRepository_UpdateStatus_NotFound(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	err := repo.UpdateStatus(ctx, "non-existent", order.StatusInProduction)
	if err != errors.ErrNotFound {
		t.Errorf("UpdateStatus() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestOrderRepository_Delete(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, o)

	err := repo.Delete(ctx, o.ID)
	if err != nil {
		t.Fatalf("Delete() unexpected error = %v", err)
	}

	// Should not be found anymore
	_, err = repo.GetByID(ctx, o.ID)
	if err != errors.ErrNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestOrderRepository_Delete_NotFound(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent")
	if err != errors.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, errors.ErrNotFound)
	}
}

func TestOrderRepository_List(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// Empty list
	orders, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("List() got %d orders, want 0", len(orders))
	}

	// Add orders
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-1",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-2",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusInProduction,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-3",
		ClientID:     "client-456",
		LaboratoryID: "lab-456", // Different laboratory
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	orders, err = repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("List() got %d orders, want 2", len(orders))
	}
}

func TestOrderRepository_List_ExcludesDeleted(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	_ = repo.Create(ctx, &order.Order{
		ID:           "order-1",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-2",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusInProduction,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	// Delete one
	_ = repo.Delete(ctx, "order-1")

	orders, err := repo.List(ctx, "lab-123")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}
	if len(orders) != 1 {
		t.Errorf("List() got %d orders, want 1", len(orders))
	}
	if orders[0].ID != "order-2" {
		t.Errorf("List() remaining order ID = %v, want order-2", orders[0].ID)
	}
}

func TestOrderRepository_ListByClientID(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// Empty list
	orders, err := repo.ListByClientID(ctx, "client-123")
	if err != nil {
		t.Fatalf("ListByClientID() unexpected error = %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("ListByClientID() got %d orders, want 0", len(orders))
	}

	// Add orders
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-1",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-2",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusInProduction,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-3",
		ClientID:     "client-456", // Different client
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	orders, err = repo.ListByClientID(ctx, "client-123")
	if err != nil {
		t.Fatalf("ListByClientID() unexpected error = %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("ListByClientID() got %d orders, want 2", len(orders))
	}
}

func TestOrderRepository_ListByClientID_ExcludesDeleted(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	_ = repo.Create(ctx, &order.Order{
		ID:           "order-1",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	_ = repo.Create(ctx, &order.Order{
		ID:           "order-2",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusInProduction,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	// Delete one
	_ = repo.Delete(ctx, "order-1")

	orders, err := repo.ListByClientID(ctx, "client-123")
	if err != nil {
		t.Fatalf("ListByClientID() unexpected error = %v", err)
	}
	if len(orders) != 1 {
		t.Errorf("ListByClientID() got %d orders, want 1", len(orders))
	}
	if orders[0].ID != "order-2" {
		t.Errorf("ListByClientID() remaining order ID = %v, want order-2", orders[0].ID)
	}
}

func TestOrderRepository_ImmutableClone(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	o := &order.Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       order.StatusReceived,
		Prosthesis: []order.ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, o)

	// Get and modify
	found, _ := repo.GetByID(ctx, o.ID)
	found.Status = order.StatusDelivered
	found.Prosthesis[0].Type = "Modified"

	// Original should be unchanged
	original, _ := repo.GetByID(ctx, o.ID)
	if original.Status != order.StatusReceived {
		t.Errorf("Repository data was mutated, got %v, want %v", original.Status, order.StatusReceived)
	}
	if original.Prosthesis[0].Type != "crown" {
		t.Errorf("Repository prosthesis was mutated, got %v, want crown", original.Prosthesis[0].Type)
	}
}
