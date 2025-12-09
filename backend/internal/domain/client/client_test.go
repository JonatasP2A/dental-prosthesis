package client

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		laboratoryID string
		clientName  string
		email       string
		phone       string
		address     Address
		wantErr     bool
		errContains string
	}{
		{
			name:         "valid client",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			name:         "empty laboratory_id",
			id:           "client-123",
			laboratoryID: "",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "laboratory_id is required",
		},
		{
			name:         "empty name",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			name:         "name too long",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   string(make([]byte, 201)),
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			name:         "invalid email format",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "invalid-email",
			phone:        "+5511999999999",
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
			name:         "empty email",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "",
			phone:        "+5511999999999",
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
			name:         "empty phone",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "",
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
			name:         "invalid phone format",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "invalid-phone",
			address: Address{
				Street:     "Rua das Flores",
				City:       "São Paulo",
				State:      "SP",
				PostalCode: "01234-567",
				Country:    "Brazil",
			},
			wantErr:     true,
			errContains: "invalid phone format",
		},
		{
			name:         "missing street",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			name:         "missing city",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
		{
			name:         "missing state",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			name:         "missing postal code",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			name:         "missing country",
			id:           "client-123",
			laboratoryID: "lab-123",
			clientName:   "Dental Clinic",
			email:        "clinic@example.com",
			phone:        "+5511999999999",
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
			client, err := NewClient(tt.id, tt.laboratoryID, tt.clientName, tt.email, tt.phone, tt.address)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewClient() expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("NewClient() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewClient() unexpected error = %v", err)
				return
			}

			if client.ID != tt.id {
				t.Errorf("NewClient() ID = %v, want %v", client.ID, tt.id)
			}
			if client.LaboratoryID != tt.laboratoryID {
				t.Errorf("NewClient() LaboratoryID = %v, want %v", client.LaboratoryID, tt.laboratoryID)
			}
			if client.Name != tt.clientName {
				t.Errorf("NewClient() Name = %v, want %v", client.Name, tt.clientName)
			}
			if client.CreatedAt.IsZero() {
				t.Errorf("NewClient() CreatedAt should not be zero")
			}
			if client.UpdatedAt.IsZero() {
				t.Errorf("NewClient() UpdatedAt should not be zero")
			}
		})
	}
}

func TestClient_Update(t *testing.T) {
	client := &Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Original Name",
		Email:        "original@example.com",
		Phone:        "+5511999999999",
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

	originalUpdatedAt := client.UpdatedAt

	newAddress := Address{
		Street:     "New Street",
		City:       "New City",
		State:      "RJ",
		PostalCode: "98765-432",
		Country:    "Brazil",
	}

	err := client.Update("New Name", "new@example.com", "+5511888888888", newAddress)
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	if client.Name != "New Name" {
		t.Errorf("Update() Name = %v, want %v", client.Name, "New Name")
	}
	if client.Email != "new@example.com" {
		t.Errorf("Update() Email = %v, want %v", client.Email, "new@example.com")
	}
	if client.Phone != "+5511888888888" {
		t.Errorf("Update() Phone = %v, want %v", client.Phone, "+5511888888888")
	}
	if !client.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Update() UpdatedAt should be after original")
	}
}

func TestClient_Update_Invalid(t *testing.T) {
	client := &Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Original Name",
		Email:        "original@example.com",
		Phone:        "+5511999999999",
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

	err := client.Update("", "invalid-email", "", Address{})
	if err == nil {
		t.Errorf("Update() expected error for invalid input, got nil")
	}
}

func TestClient_Delete(t *testing.T) {
	client := &Client{
		ID:           "client-123",
		LaboratoryID: "lab-123",
		Name:         "Test Client",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if client.IsDeleted() {
		t.Errorf("IsDeleted() should be false before delete")
	}

	client.Delete()

	if !client.IsDeleted() {
		t.Errorf("IsDeleted() should be true after delete")
	}
	if client.DeletedAt == nil {
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
			name: "missing street",
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
			name: "missing city",
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

func TestClient_Validate_EmailFormats(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "test@sub.example.com", false},
		{"valid email with plus", "test+tag@example.com", false},
		{"valid email with dot", "test.name@example.com", false},
		{"invalid email - no @", "testexample.com", true},
		{"invalid email - no domain", "test@", true},
		{"invalid email - no tld", "test@example", true},
		{"invalid email - spaces", "test @example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				ID:           "client-123",
				LaboratoryID: "lab-123",
				Name:         "Test Client",
				Email:        tt.email,
				Phone:        "+5511999999999",
				Address: Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			}

			err := client.Validate()
			hasErr := err != nil

			if hasErr != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Validate_PhoneFormats(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid phone with +", "+5511999999999", false},
		{"valid phone without +", "5511999999999", false},
		{"valid phone short", "+1234567890", false},
		{"invalid phone - letters", "abc123", true},
		{"invalid phone - spaces", "+55 11 99999 9999", false}, // Spaces are removed
		{"invalid phone - special chars", "+55-11-99999-9999", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				ID:           "client-123",
				LaboratoryID: "lab-123",
				Name:         "Test Client",
				Email:        "test@example.com",
				Phone:        tt.phone,
				Address: Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			}

			err := client.Validate()
			hasErr := err != nil

			if hasErr != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
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
