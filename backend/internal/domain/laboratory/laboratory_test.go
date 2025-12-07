package laboratory

import (
	"testing"
	"time"
)

func TestNewLaboratory(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		labName     string
		email       string
		phone       string
		address     Address
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid laboratory",
			id:      "lab-123",
			labName: "Dental Lab",
			email:   "contact@dentallab.com",
			phone:   "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr: false,
		},
		{
			name:    "empty name",
			id:      "lab-123",
			labName: "",
			email:   "contact@dentallab.com",
			phone:   "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name:    "name too long",
			id:      "lab-123",
			labName: string(make([]byte, 201)),
			email:   "contact@dentallab.com",
			phone:   "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "name must be at most 200 characters",
		},
		{
			name:    "invalid email format",
			id:      "lab-123",
			labName: "Dental Lab",
			email:   "invalid-email",
			phone:   "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "invalid email format",
		},
		{
			name:    "empty email",
			id:      "lab-123",
			labName: "Dental Lab",
			email:   "",
			phone:   "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "email is required",
		},
		{
			name:    "empty phone",
			id:      "lab-123",
			labName: "Dental Lab",
			email:   "contact@dentallab.com",
			phone:   "",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "phone is required",
		},
		{
			name:    "missing street",
			id:      "lab-123",
			labName: "Dental Lab",
			email:   "contact@dentallab.com",
			phone:   "+5511999999999",
			address: Address{
				Street:     "",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "street is required",
		},
		{
			name:    "missing city",
			id:      "lab-123",
			labName: "Dental Lab",
			email:   "contact@dentallab.com",
			phone:   "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "city is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lab, err := NewLaboratory(tt.id, tt.labName, tt.email, tt.phone, tt.address)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewLaboratory() expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("NewLaboratory() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewLaboratory() unexpected error = %v", err)
				return
			}

			if lab.ID != tt.id {
				t.Errorf("NewLaboratory() ID = %v, want %v", lab.ID, tt.id)
			}
			if lab.Name != tt.labName {
				t.Errorf("NewLaboratory() Name = %v, want %v", lab.Name, tt.labName)
			}
			if lab.CreatedAt.IsZero() {
				t.Errorf("NewLaboratory() CreatedAt should not be zero")
			}
		})
	}
}

func TestLaboratory_Update(t *testing.T) {
	lab := &Laboratory{
		ID:        "lab-123",
		Name:      "Original Name",
		Email:     "original@lab.com",
		Phone:     "+5511999999999",
		Address: Address{
			Street:     "Original Street",
			City:       "Original City",
			State:      "SP",
			PostalCode: "01234-567",
			Country:    "Brazil",
		},
		CreatedAt: time.Now().UTC().Add(-time.Hour),
		UpdatedAt: time.Now().UTC().Add(-time.Hour),
	}

	originalUpdatedAt := lab.UpdatedAt

	newAddress := Address{
		Street:     "New Street",
		City:       "New City",
		State:      "RJ",
		PostalCode: "98765-432",
		Country:    "Brazil",
	}

	err := lab.Update("New Name", "new@lab.com", "+5511888888888", newAddress)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	if lab.Name != "New Name" {
		t.Errorf("Update() Name = %v, want %v", lab.Name, "New Name")
	}
	if lab.Email != "new@lab.com" {
		t.Errorf("Update() Email = %v, want %v", lab.Email, "new@lab.com")
	}
	if !lab.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Update() UpdatedAt should be after original")
	}
}

func TestLaboratory_Update_Invalid(t *testing.T) {
	lab := &Laboratory{
		ID:        "lab-123",
		Name:      "Original Name",
		Email:     "original@lab.com",
		Phone:     "+5511999999999",
		Address: Address{
			Street:     "Original Street",
			City:       "Original City",
			State:      "SP",
			PostalCode: "01234-567",
			Country:    "Brazil",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := lab.Update("", "invalid-email", "", Address{})
	if err == nil {
		t.Errorf("Update() expected error for invalid input, got nil")
	}
}

func TestLaboratory_Delete(t *testing.T) {
	lab := &Laboratory{
		ID:        "lab-123",
		Name:      "Test Lab",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if lab.IsDeleted() {
		t.Errorf("IsDeleted() should be false before delete")
	}

	lab.Delete()

	if !lab.IsDeleted() {
		t.Errorf("IsDeleted() should be true after delete")
	}
	if lab.DeletedAt == nil {
		t.Errorf("DeletedAt should not be nil after delete")
	}
}

func TestAddress_Validate(t *testing.T) {
	tests := []struct {
		name        string
		address     Address
		wantErr     bool
		errContains string
	}{
		{
			name: "valid address",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr: false,
		},
		{
			name: "missing state",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "state is required",
		},
		{
			name: "missing postal code",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "postal code is required",
		},
		{
			name: "missing country",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "",
			},
			wantErr:     true,
			errContains: "country is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Validate()

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

