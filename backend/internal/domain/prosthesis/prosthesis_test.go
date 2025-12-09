package prosthesis

import (
	"strings"
	"testing"
	"time"
)

func TestNewProsthesis(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		laboratoryID    string
		prosthesisType  ProsthesisType
		material        string
		shade           string
		specifications  string
		notes           string
		wantErr         bool
		errContains     string
	}{
		{
			name:           "valid prosthesis - crown",
			id:             "prosthesis-123",
			laboratoryID:   "lab-123",
			prosthesisType: ProsthesisTypeCrown,
			material:       "zirconia",
			shade:          "A1",
			specifications: "Full coverage",
			notes:          "High priority",
			wantErr:        false,
		},
		{
			name:           "valid prosthesis - bridge",
			id:             "prosthesis-124",
			laboratoryID:   "lab-123",
			prosthesisType: ProsthesisTypeBridge,
			material:       "porcelain",
			shade:          "A2",
			wantErr:        false,
		},
		{
			name:           "valid prosthesis - complete denture",
			id:             "prosthesis-125",
			laboratoryID:   "lab-123",
			prosthesisType: ProsthesisTypeCompleteDenture,
			material:       "acrylic",
			wantErr:        false,
		},
		{
			name:           "empty laboratory_id",
			id:             "prosthesis-123",
			laboratoryID:   "",
			prosthesisType: ProsthesisTypeCrown,
			material:       "zirconia",
			wantErr:        true,
			errContains:    "laboratory_id is required",
		},
		{
			name:           "empty type",
			id:             "prosthesis-123",
			laboratoryID:   "lab-123",
			prosthesisType: "",
			material:       "zirconia",
			wantErr:        true,
			errContains:    "type is required",
		},
		{
			name:           "invalid type",
			id:             "prosthesis-123",
			laboratoryID:   "lab-123",
			prosthesisType: ProsthesisType("invalid_type"),
			material:       "zirconia",
			wantErr:        true,
			errContains:    "invalid prosthesis type",
		},
		{
			name:           "empty material",
			id:             "prosthesis-123",
			laboratoryID:   "lab-123",
			prosthesisType: ProsthesisTypeCrown,
			material:       "",
			wantErr:        true,
			errContains:    "material is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewProsthesis(tt.id, tt.laboratoryID, tt.prosthesisType, tt.material, tt.shade, tt.specifications, tt.notes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProsthesis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewProsthesis() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}
			if p == nil {
				t.Fatal("NewProsthesis() returned nil prosthesis")
			}
			if p.ID != tt.id {
				t.Errorf("NewProsthesis() ID = %v, want %v", p.ID, tt.id)
			}
			if p.LaboratoryID != tt.laboratoryID {
				t.Errorf("NewProsthesis() LaboratoryID = %v, want %v", p.LaboratoryID, tt.laboratoryID)
			}
			if p.Type != tt.prosthesisType {
				t.Errorf("NewProsthesis() Type = %v, want %v", p.Type, tt.prosthesisType)
			}
			if p.Material != tt.material {
				t.Errorf("NewProsthesis() Material = %v, want %v", p.Material, tt.material)
			}
			if p.Shade != tt.shade {
				t.Errorf("NewProsthesis() Shade = %v, want %v", p.Shade, tt.shade)
			}
			if p.Specifications != tt.specifications {
				t.Errorf("NewProsthesis() Specifications = %v, want %v", p.Specifications, tt.specifications)
			}
			if p.Notes != tt.notes {
				t.Errorf("NewProsthesis() Notes = %v, want %v", p.Notes, tt.notes)
			}
			if p.CreatedAt.IsZero() {
				t.Error("NewProsthesis() CreatedAt should be set")
			}
			if p.UpdatedAt.IsZero() {
				t.Error("NewProsthesis() UpdatedAt should be set")
			}
		})
	}
}

func TestProsthesisType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		pt    ProsthesisType
		valid bool
	}{
		{"crown", ProsthesisTypeCrown, true},
		{"bridge", ProsthesisTypeBridge, true},
		{"complete_denture", ProsthesisTypeCompleteDenture, true},
		{"partial_denture", ProsthesisTypePartialDenture, true},
		{"implant", ProsthesisTypeImplant, true},
		{"veneer", ProsthesisTypeVeneer, true},
		{"inlay", ProsthesisTypeInlay, true},
		{"onlay", ProsthesisTypeOnlay, true},
		{"invalid", ProsthesisType("invalid"), false},
		{"empty", ProsthesisType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pt.IsValid(); got != tt.valid {
				t.Errorf("ProsthesisType.IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestProsthesis_Update(t *testing.T) {
	p, err := NewProsthesis("prosthesis-123", "lab-123", ProsthesisTypeCrown, "zirconia", "A1", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}

	originalUpdatedAt := p.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

	err = p.Update(ProsthesisTypeBridge, "porcelain", "A2", "3-unit bridge", "Urgent")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if p.Type != ProsthesisTypeBridge {
		t.Errorf("Update() Type = %v, want %v", p.Type, ProsthesisTypeBridge)
	}
	if p.Material != "porcelain" {
		t.Errorf("Update() Material = %v, want %v", p.Material, "porcelain")
	}
	if p.Shade != "A2" {
		t.Errorf("Update() Shade = %v, want %v", p.Shade, "A2")
	}
	if p.Specifications != "3-unit bridge" {
		t.Errorf("Update() Specifications = %v, want %v", p.Specifications, "3-unit bridge")
	}
	if p.Notes != "Urgent" {
		t.Errorf("Update() Notes = %v, want %v", p.Notes, "Urgent")
	}
	if !p.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Update() UpdatedAt should be updated")
	}

	// Test invalid update
	err = p.Update(ProsthesisType("invalid"), "material", "", "", "")
	if err == nil {
		t.Error("Update() with invalid type should return error")
	}
}

func TestProsthesis_Delete(t *testing.T) {
	p, err := NewProsthesis("prosthesis-123", "lab-123", ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}

	if p.IsDeleted() {
		t.Error("Prosthesis should not be deleted initially")
	}

	p.Delete()

	if !p.IsDeleted() {
		t.Error("Prosthesis should be deleted")
	}
	if p.DeletedAt == nil {
		t.Error("DeletedAt should be set")
	}
	if p.DeletedAt.IsZero() {
		t.Error("DeletedAt should not be zero")
	}
}

func TestProsthesis_IsDeleted(t *testing.T) {
	p, err := NewProsthesis("prosthesis-123", "lab-123", ProsthesisTypeCrown, "zirconia", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create prosthesis: %v", err)
	}

	if p.IsDeleted() {
		t.Error("IsDeleted() = true, want false")
	}

	p.Delete()

	if !p.IsDeleted() {
		t.Error("IsDeleted() = false, want true")
	}
}

func TestProsthesis_Validate(t *testing.T) {
	tests := []struct {
		name         string
		prosthesis   *Prosthesis
		wantErr      bool
		errContains  string
	}{
		{
			name: "valid prosthesis",
			prosthesis: &Prosthesis{
				ID:           "prosthesis-123",
				LaboratoryID: "lab-123",
				Type:         ProsthesisTypeCrown,
				Material:     "zirconia",
			},
			wantErr: false,
		},
		{
			name: "empty laboratory_id",
			prosthesis: &Prosthesis{
				ID:       "prosthesis-123",
				Type:     ProsthesisTypeCrown,
				Material: "zirconia",
			},
			wantErr:     true,
			errContains: "laboratory_id is required",
		},
		{
			name: "empty type",
			prosthesis: &Prosthesis{
				ID:           "prosthesis-123",
				LaboratoryID: "lab-123",
				Material:     "zirconia",
			},
			wantErr:     true,
			errContains: "type is required",
		},
		{
			name: "invalid type",
			prosthesis: &Prosthesis{
				ID:           "prosthesis-123",
				LaboratoryID: "lab-123",
				Type:         ProsthesisType("invalid"),
				Material:     "zirconia",
			},
			wantErr:     true,
			errContains: "invalid prosthesis type",
		},
		{
			name: "empty material",
			prosthesis: &Prosthesis{
				ID:           "prosthesis-123",
				LaboratoryID: "lab-123",
				Type:         ProsthesisTypeCrown,
			},
			wantErr:     true,
			errContains: "material is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prosthesis.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}
