package technician

import (
	"testing"
	"time"
)

func TestNewTechnician(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		laboratoryID    string
		techName        string
		email           string
		phone           string
		role            Role
		specializations []string
		wantErr         bool
		errContains     string
	}{
		{
			name:         "valid technician",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         RoleTechnician,
			wantErr:      false,
		},
		{
			name:         "valid technician with specializations",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         RoleSeniorTechnician,
			specializations: []string{"crowns", "dentures"},
			wantErr:      false,
		},
		{
			name:         "empty laboratory_id",
			id:           "tech-123",
			laboratoryID: "",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "laboratory_id is required",
		},
		{
			name:         "empty name",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "",
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "name is required",
		},
		{
			name:         "name too long",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     string(make([]byte, 201)),
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "name must be at most 200 characters",
		},
		{
			name:         "invalid email format",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "invalid-email",
			phone:        "+5511999999999",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "invalid email format",
		},
		{
			name:         "empty email",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "",
			phone:        "+5511999999999",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "email is required",
		},
		{
			name:         "empty phone",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "phone is required",
		},
		{
			name:         "invalid phone format",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "invalid",
			role:         RoleTechnician,
			wantErr:      true,
			errContains:  "invalid phone format",
		},
		{
			name:         "empty role",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         Role(""),
			wantErr:      true,
			errContains:  "role is required",
		},
		{
			name:         "invalid role",
			id:           "tech-123",
			laboratoryID: "lab-123",
			techName:     "John Doe",
			email:        "john@lab.com",
			phone:        "+5511999999999",
			role:         Role("invalid_role"),
			wantErr:      true,
			errContains:  "invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tech, err := NewTechnician(tt.id, tt.laboratoryID, tt.techName, tt.email, tt.phone, tt.role, tt.specializations)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTechnician() expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("NewTechnician() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTechnician() unexpected error = %v", err)
				return
			}

			if tech.ID != tt.id {
				t.Errorf("NewTechnician() ID = %v, want %v", tech.ID, tt.id)
			}
			if tech.Name != tt.techName {
				t.Errorf("NewTechnician() Name = %v, want %v", tech.Name, tt.techName)
			}
			if tech.Role != tt.role {
				t.Errorf("NewTechnician() Role = %v, want %v", tech.Role, tt.role)
			}
			if tech.CreatedAt.IsZero() {
				t.Errorf("NewTechnician() CreatedAt should not be zero")
			}
		})
	}
}

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{
			name: "valid senior_technician",
			role: RoleSeniorTechnician,
			want: true,
		},
		{
			name: "valid technician",
			role: RoleTechnician,
			want: true,
		},
		{
			name: "valid apprentice",
			role: RoleApprentice,
			want: true,
		},
		{
			name: "invalid role",
			role: Role("invalid"),
			want: false,
		},
		{
			name: "empty role",
			role: Role(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("Role.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechnician_Update(t *testing.T) {
	tech := &Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "Original Name",
		Email:        "original@lab.com",
		Phone:        "+5511999999999",
		Role:         RoleTechnician,
		CreatedAt:    time.Now().UTC().Add(-time.Hour),
		UpdatedAt:    time.Now().UTC().Add(-time.Hour),
	}

	originalUpdatedAt := tech.UpdatedAt

	err := tech.Update("New Name", "new@lab.com", "+5511888888888", RoleSeniorTechnician, []string{"crowns"})
	if err != nil {
		t.Fatalf("Update() unexpected error = %v", err)
	}

	if tech.Name != "New Name" {
		t.Errorf("Update() Name = %v, want %v", tech.Name, "New Name")
	}
	if tech.Email != "new@lab.com" {
		t.Errorf("Update() Email = %v, want %v", tech.Email, "new@lab.com")
	}
	if tech.Role != RoleSeniorTechnician {
		t.Errorf("Update() Role = %v, want %v", tech.Role, RoleSeniorTechnician)
	}
	if !tech.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Update() UpdatedAt should be after original")
	}
}

func TestTechnician_Update_Invalid(t *testing.T) {
	tech := &Technician{
		ID:           "tech-123",
		LaboratoryID: "lab-123",
		Name:         "Original Name",
		Email:        "original@lab.com",
		Phone:        "+5511999999999",
		Role:         RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err := tech.Update("", "invalid-email", "", Role("invalid"), nil)
	if err == nil {
		t.Errorf("Update() expected error for invalid input, got nil")
	}
}

func TestTechnician_Delete(t *testing.T) {
	tech := &Technician{
		ID:        "tech-123",
		Name:      "Test Tech",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if tech.IsDeleted() {
		t.Errorf("IsDeleted() should be false before delete")
	}

	tech.Delete()

	if !tech.IsDeleted() {
		t.Errorf("IsDeleted() should be true after delete")
	}
	if tech.DeletedAt == nil {
		t.Errorf("DeletedAt should not be nil after delete")
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
