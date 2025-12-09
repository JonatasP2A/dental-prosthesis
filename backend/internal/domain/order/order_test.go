package order

import (
	"testing"
	"time"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		clientID     string
		laboratoryID string
		prosthesis   []ProsthesisItem
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid order",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis: []ProsthesisItem{
				{
					Type:     "crown",
					Material: "zirconia",
					Shade:    "A1",
					Quantity: 1,
					Notes:    "Test notes",
				},
			},
			wantErr: false,
		},
		{
			name:         "empty client_id",
			id:           "order-123",
			clientID:     "",
			laboratoryID: "lab-123",
			prosthesis: []ProsthesisItem{
				{
					Type:     "crown",
					Material: "zirconia",
					Shade:    "A1",
					Quantity: 1,
				},
			},
			wantErr:     true,
			errContains: "client_id is required",
		},
		{
			name:         "empty laboratory_id",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "",
			prosthesis: []ProsthesisItem{
				{
					Type:     "crown",
					Material: "zirconia",
					Shade:    "A1",
					Quantity: 1,
				},
			},
			wantErr:     true,
			errContains: "laboratory_id is required",
		},
		{
			name:         "empty prosthesis items",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis:   []ProsthesisItem{},
			wantErr:     true,
			errContains: "at least one prosthesis item is required",
		},
		{
			name:         "nil prosthesis items",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis:   nil,
			wantErr:     true,
			errContains: "at least one prosthesis item is required",
		},
		{
			name:         "prosthesis item missing type",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis: []ProsthesisItem{
				{
					Type:     "",
					Material: "zirconia",
					Shade:    "A1",
					Quantity: 1,
				},
			},
			wantErr:     true,
			errContains: "type is required",
		},
		{
			name:         "prosthesis item missing material",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis: []ProsthesisItem{
				{
					Type:     "crown",
					Material: "",
					Shade:    "A1",
					Quantity: 1,
				},
			},
			wantErr:     true,
			errContains: "material is required",
		},
		{
			name:         "prosthesis item invalid quantity",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis: []ProsthesisItem{
				{
					Type:     "crown",
					Material: "zirconia",
					Shade:    "A1",
					Quantity: 0,
				},
			},
			wantErr:     true,
			errContains: "quantity must be greater than 0",
		},
		{
			name:         "prosthesis item negative quantity",
			id:           "order-123",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			prosthesis: []ProsthesisItem{
				{
					Type:     "crown",
					Material: "zirconia",
					Shade:    "A1",
					Quantity: -1,
				},
			},
			wantErr:     true,
			errContains: "quantity must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.id, tt.clientID, tt.laboratoryID, tt.prosthesis)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewOrder() expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("NewOrder() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewOrder() unexpected error = %v", err)
				return
			}

			if order.ID != tt.id {
				t.Errorf("NewOrder() ID = %v, want %v", order.ID, tt.id)
			}
			if order.ClientID != tt.clientID {
				t.Errorf("NewOrder() ClientID = %v, want %v", order.ClientID, tt.clientID)
			}
			if order.LaboratoryID != tt.laboratoryID {
				t.Errorf("NewOrder() LaboratoryID = %v, want %v", order.LaboratoryID, tt.laboratoryID)
			}
			if order.Status != StatusReceived {
				t.Errorf("NewOrder() Status = %v, want %v", order.Status, StatusReceived)
			}
			if order.CreatedAt.IsZero() {
				t.Errorf("NewOrder() CreatedAt should not be zero")
			}
			if order.UpdatedAt.IsZero() {
				t.Errorf("NewOrder() UpdatedAt should not be zero")
			}
		})
	}
}

func TestOrder_Update(t *testing.T) {
	order := &Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       StatusReceived,
		Prosthesis: []ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC().Add(-time.Hour),
		UpdatedAt: time.Now().UTC().Add(-time.Hour),
	}

	originalUpdatedAt := order.UpdatedAt

	newItems := []ProsthesisItem{
		{
			Type:     "bridge",
			Material: "porcelain",
			Shade:    "A2",
			Quantity: 2,
			Notes:    "Updated notes",
		},
	}

	err := order.Update(newItems)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	if len(order.Prosthesis) != 1 {
		t.Errorf("Update() Prosthesis length = %v, want 1", len(order.Prosthesis))
	}
	if order.Prosthesis[0].Type != "bridge" {
		t.Errorf("Update() Prosthesis[0].Type = %v, want bridge", order.Prosthesis[0].Type)
	}
	if !order.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Update() UpdatedAt should be after original")
	}
}

func TestOrder_Update_Invalid(t *testing.T) {
	order := &Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       StatusReceived,
		Prosthesis: []ProsthesisItem{
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

	err := order.Update([]ProsthesisItem{})
	if err == nil {
		t.Errorf("Update() expected error for invalid input, got nil")
	}
}

func TestOrder_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name       string
		current    Status
		newStatus  Status
		wantResult bool
	}{
		{"received -> in_production", StatusReceived, StatusInProduction, true},
		{"in_production -> quality_check", StatusInProduction, StatusQualityCheck, true},
		{"quality_check -> ready", StatusQualityCheck, StatusReady, true},
		{"quality_check -> revision", StatusQualityCheck, StatusRevision, true},
		{"ready -> delivered", StatusReady, StatusDelivered, true},
		{"ready -> revision", StatusReady, StatusRevision, true},
		{"revision -> in_production", StatusRevision, StatusInProduction, true},
		{"received -> ready", StatusReceived, StatusReady, false},
		{"received -> delivered", StatusReceived, StatusDelivered, false},
		{"in_production -> received", StatusInProduction, StatusReceived, false},
		{"delivered -> ready", StatusDelivered, StatusReady, false},
		{"delivered -> any", StatusDelivered, StatusInProduction, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{
				ID:           "order-123",
				ClientID:     "client-123",
				LaboratoryID: "lab-123",
				Status:       tt.current,
				Prosthesis: []ProsthesisItem{
					{
						Type:     "crown",
						Material: "zirconia",
						Quantity: 1,
					},
				},
			}

			result := order.CanTransitionTo(tt.newStatus)
			if result != tt.wantResult {
				t.Errorf("CanTransitionTo() = %v, want %v", result, tt.wantResult)
			}
		})
	}
}

func TestOrder_UpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		current   Status
		newStatus Status
		wantErr   bool
	}{
		{"valid transition: received -> in_production", StatusReceived, StatusInProduction, false},
		{"valid transition: in_production -> quality_check", StatusInProduction, StatusQualityCheck, false},
		{"valid transition: quality_check -> ready", StatusQualityCheck, StatusReady, false},
		{"valid transition: quality_check -> revision", StatusQualityCheck, StatusRevision, false},
		{"valid transition: ready -> delivered", StatusReady, StatusDelivered, false},
		{"valid transition: revision -> in_production", StatusRevision, StatusInProduction, false},
		{"invalid transition: received -> ready", StatusReceived, StatusReady, true},
		{"invalid transition: delivered -> ready", StatusDelivered, StatusReady, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{
				ID:           "order-123",
				ClientID:     "client-123",
				LaboratoryID: "lab-123",
				Status:       tt.current,
				Prosthesis: []ProsthesisItem{
					{
						Type:     "crown",
						Material: "zirconia",
						Quantity: 1,
					},
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC().Add(-time.Hour),
			}

			originalUpdatedAt := order.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			err := order.UpdateStatus(tt.newStatus)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateStatus() expected error, got nil")
					return
				}
				if err != errors.ErrInvalidStatusTransition {
					t.Errorf("UpdateStatus() error = %v, want %v", err, errors.ErrInvalidStatusTransition)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error = %v", err)
				return
			}

			if order.Status != tt.newStatus {
				t.Errorf("UpdateStatus() Status = %v, want %v", order.Status, tt.newStatus)
			}
			if !order.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateStatus() UpdatedAt should be after original")
			}
		})
	}
}

func TestOrder_Delete(t *testing.T) {
	order := &Order{
		ID:           "order-123",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       StatusReceived,
		Prosthesis: []ProsthesisItem{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if order.IsDeleted() {
		t.Errorf("IsDeleted() should be false before delete")
	}

	order.Delete()

	if !order.IsDeleted() {
		t.Errorf("IsDeleted() should be true after delete")
	}
	if order.DeletedAt == nil {
		t.Errorf("DeletedAt should not be nil after delete")
	}
}

func TestProsthesisItem_Validate(t *testing.T) {
	tests := []struct {
		name        string
		item        ProsthesisItem
		index       int
		wantErr     bool
		errContains string
	}{
		{
			name: "valid item",
			item: ProsthesisItem{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 1,
				Notes:    "Test notes",
			},
			index:   0,
			wantErr: false,
		},
		{
			name: "missing type",
			item: ProsthesisItem{
				Type:     "",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 1,
			},
			index:       0,
			wantErr:     true,
			errContains: "type is required",
		},
		{
			name: "missing material",
			item: ProsthesisItem{
				Type:     "crown",
				Material: "",
				Shade:    "A1",
				Quantity: 1,
			},
			index:       0,
			wantErr:     true,
			errContains: "material is required",
		},
		{
			name: "zero quantity",
			item: ProsthesisItem{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 0,
			},
			index:       0,
			wantErr:     true,
			errContains: "quantity must be greater than 0",
		},
		{
			name: "negative quantity",
			item: ProsthesisItem{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: -1,
			},
			index:       0,
			wantErr:     true,
			errContains: "quantity must be greater than 0",
		},
		{
			name: "valid item with shade optional",
			item: ProsthesisItem{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "",
				Quantity: 1,
			},
			index:   0,
			wantErr: false,
		},
		{
			name: "valid item with notes optional",
			item: ProsthesisItem{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 1,
				Notes:    "",
			},
			index:   0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate(tt.index)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid: received", "received", true},
		{"valid: in_production", "in_production", true},
		{"valid: quality_check", "quality_check", true},
		{"valid: ready", "ready", true},
		{"valid: delivered", "delivered", true},
		{"valid: revision", "revision", true},
		{"invalid: unknown", "unknown", false},
		{"invalid: empty", "", false},
		{"invalid: mixed case", "Received", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidStatus(tt.status)
			if result != tt.want {
				t.Errorf("IsValidStatus(%v) = %v, want %v", tt.status, result, tt.want)
			}
		})
	}
}

func TestAllStatuses(t *testing.T) {
	statuses := AllStatuses()
	expected := []Status{
		StatusReceived,
		StatusInProduction,
		StatusQualityCheck,
		StatusReady,
		StatusDelivered,
		StatusRevision,
	}

	if len(statuses) != len(expected) {
		t.Errorf("AllStatuses() length = %v, want %v", len(statuses), len(expected))
	}

	for _, expectedStatus := range expected {
		found := false
		for _, status := range statuses {
			if status == expectedStatus {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AllStatuses() missing status %v", expectedStatus)
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
