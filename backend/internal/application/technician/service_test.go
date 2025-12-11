package technician

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
)

// mockIDGenerator is a mock ID generator for testing
type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) Generate() string {
	return m.id
}

// mockTechnicianRepository is a mock repository for testing
type mockTechnicianRepository struct {
	techs          map[string]*technician.Technician
	getByIDErr     error
	getByEmailErr  error
	createErr      error
	updateErr      error
	deleteErr      error
	listErr        error
	listByRoleErr  error
}

func newMockTechnicianRepository() *mockTechnicianRepository {
	return &mockTechnicianRepository{
		techs: make(map[string]*technician.Technician),
	}
}

func (m *mockTechnicianRepository) Create(ctx context.Context, tech *technician.Technician) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.techs[tech.ID] = tech
	return nil
}

func (m *mockTechnicianRepository) GetByID(ctx context.Context, id string) (*technician.Technician, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	tech, exists := m.techs[id]
	if !exists || tech.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return tech, nil
}

func (m *mockTechnicianRepository) GetByEmail(ctx context.Context, laboratoryID, email string) (*technician.Technician, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	for _, tech := range m.techs {
		if tech.Email == email && tech.LaboratoryID == laboratoryID && !tech.IsDeleted() {
			return tech, nil
		}
	}
	return nil, errors.ErrNotFound
}

func (m *mockTechnicianRepository) Update(ctx context.Context, tech *technician.Technician) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.techs[tech.ID]; !exists {
		return errors.ErrNotFound
	}
	m.techs[tech.ID] = tech
	return nil
}

func (m *mockTechnicianRepository) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	tech, exists := m.techs[id]
	if !exists {
		return errors.ErrNotFound
	}
	tech.Delete()
	return nil
}

func (m *mockTechnicianRepository) List(ctx context.Context, laboratoryID string) ([]*technician.Technician, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var techs []*technician.Technician
	for _, tech := range m.techs {
		if tech.LaboratoryID == laboratoryID && !tech.IsDeleted() {
			techs = append(techs, tech)
		}
	}
	return techs, nil
}

func (m *mockTechnicianRepository) ListByRole(ctx context.Context, laboratoryID string, role technician.Role) ([]*technician.Technician, error) {
	if m.listByRoleErr != nil {
		return nil, m.listByRoleErr
	}
	var techs []*technician.Technician
	for _, tech := range m.techs {
		if tech.LaboratoryID == laboratoryID && tech.Role == role && !tech.IsDeleted() {
			techs = append(techs, tech)
		}
	}
	return techs, nil
}

// mockLaboratoryRepository is a mock laboratory repository for testing
type mockLaboratoryRepository struct {
	labs map[string]*laboratory.Laboratory
}

func newMockLaboratoryRepository() *mockLaboratoryRepository {
	return &mockLaboratoryRepository{
		labs: make(map[string]*laboratory.Laboratory),
	}
}

func (m *mockLaboratoryRepository) Create(ctx context.Context, lab *laboratory.Laboratory) error {
	m.labs[lab.ID] = lab
	return nil
}

func (m *mockLaboratoryRepository) GetByID(ctx context.Context, id string) (*laboratory.Laboratory, error) {
	lab, exists := m.labs[id]
	if !exists || lab.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return lab, nil
}

func (m *mockLaboratoryRepository) GetByEmail(ctx context.Context, email string) (*laboratory.Laboratory, error) {
	for _, lab := range m.labs {
		if lab.Email == email && !lab.IsDeleted() {
			return lab, nil
		}
	}
	return nil, errors.ErrNotFound
}

func (m *mockLaboratoryRepository) Update(ctx context.Context, lab *laboratory.Laboratory) error {
	if _, exists := m.labs[lab.ID]; !exists {
		return errors.ErrNotFound
	}
	m.labs[lab.ID] = lab
	return nil
}

func (m *mockLaboratoryRepository) Delete(ctx context.Context, id string) error {
	lab, exists := m.labs[id]
	if !exists {
		return errors.ErrNotFound
	}
	lab.Delete()
	return nil
}

func (m *mockLaboratoryRepository) List(ctx context.Context) ([]*laboratory.Laboratory, error) {
	var labs []*laboratory.Laboratory
	for _, lab := range m.labs {
		if !lab.IsDeleted() {
			labs = append(labs, lab)
		}
	}
	return labs, nil
}

func TestService_CreateTechnician(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateInput
		mockID    string
		setupRepo func(*mockTechnicianRepository, *mockLaboratoryRepository)
		wantErr   error
	}{
		{
			name: "successfully create technician",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "John Doe",
				Email:        "john@lab.com",
				Phone:        "+5511999999999",
				Role:         technician.RoleTechnician,
			},
			mockID: "tech-123",
			setupRepo: func(techRepo *mockTechnicianRepository, labRepo *mockLaboratoryRepository) {
				labRepo.labs["lab-123"] = &laboratory.Laboratory{ID: "lab-123"}
			},
			wantErr: nil,
		},
		{
			name: "laboratory not found",
			input: CreateInput{
				LaboratoryID: "non-existent",
				Name:         "John Doe",
				Email:        "john@lab.com",
				Phone:        "+5511999999999",
				Role:         technician.RoleTechnician,
			},
			mockID:    "tech-123",
			setupRepo: func(techRepo *mockTechnicianRepository, labRepo *mockLaboratoryRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "duplicate email within laboratory",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "John Doe",
				Email:        "existing@lab.com",
				Phone:        "+5511999999999",
				Role:         technician.RoleTechnician,
			},
			mockID: "tech-123",
			setupRepo: func(techRepo *mockTechnicianRepository, labRepo *mockLaboratoryRepository) {
				labRepo.labs["lab-123"] = &laboratory.Laboratory{ID: "lab-123"}
				techRepo.techs["existing"] = &technician.Technician{
					ID:           "existing",
					LaboratoryID: "lab-123",
					Email:        "existing@lab.com",
				}
			},
			wantErr: errors.ErrDuplicateEmail,
		},
		{
			name: "invalid input - empty name",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "",
				Email:        "john@lab.com",
				Phone:        "+5511999999999",
				Role:         technician.RoleTechnician,
			},
			mockID: "tech-123",
			setupRepo: func(techRepo *mockTechnicianRepository, labRepo *mockLaboratoryRepository) {
				labRepo.labs["lab-123"] = &laboratory.Laboratory{ID: "lab-123"}
			},
			wantErr: errors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			techRepo := newMockTechnicianRepository()
			labRepo := newMockLaboratoryRepository()
			tt.setupRepo(techRepo, labRepo)
			idGen := &mockIDGenerator{id: tt.mockID}
			svc := NewService(techRepo, labRepo, idGen)

			tech, err := svc.CreateTechnician(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateTechnician() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					// Check if it's a validation error
					var ve errors.ValidationErrors
					if stderrors.As(err, &ve) && tt.wantErr == errors.ErrInvalidInput {
						return // ValidationErrors is considered ErrInvalidInput
					}
					t.Errorf("CreateTechnician() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateTechnician() unexpected error = %v", err)
				return
			}

			if tech.ID != tt.mockID {
				t.Errorf("CreateTechnician() ID = %v, want %v", tech.ID, tt.mockID)
			}
			if tech.Name != tt.input.Name {
				t.Errorf("CreateTechnician() Name = %v, want %v", tech.Name, tt.input.Name)
			}
		})
	}
}

func TestService_GetTechnician(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		labID     string
		setupRepo func(*mockTechnicianRepository)
		wantErr   error
	}{
		{
			name: "successfully get technician",
			id:   "tech-123",
			labID: "lab-123",
			setupRepo: func(r *mockTechnicianRepository) {
				r.techs["tech-123"] = &technician.Technician{
					ID:           "tech-123",
					LaboratoryID: "lab-123",
					Name:         "John Doe",
				}
			},
			wantErr: nil,
		},
		{
			name:      "technician not found",
			id:        "non-existent",
			labID:     "lab-123",
			setupRepo: func(r *mockTechnicianRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "cross-laboratory access",
			id:   "tech-123",
			labID: "lab-999",
			setupRepo: func(r *mockTechnicianRepository) {
				r.techs["tech-123"] = &technician.Technician{
					ID:           "tech-123",
					LaboratoryID: "lab-123",
					Name:         "John Doe",
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			techRepo := newMockTechnicianRepository()
			tt.setupRepo(techRepo)
			svc := NewService(techRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

			tech, err := svc.GetTechnician(context.Background(), tt.id, tt.labID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetTechnician() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("GetTechnician() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetTechnician() unexpected error = %v", err)
				return
			}

			if tech.ID != tt.id {
				t.Errorf("GetTechnician() ID = %v, want %v", tech.ID, tt.id)
			}
		})
	}
}

func TestService_ListTechnicians(t *testing.T) {
	techRepo := newMockTechnicianRepository()
	techRepo.techs["tech-1"] = &technician.Technician{
		ID:           "tech-1",
		LaboratoryID: "lab-123",
		Name:         "Tech 1",
		Role:         technician.RoleTechnician,
	}
	techRepo.techs["tech-2"] = &technician.Technician{
		ID:           "tech-2",
		LaboratoryID: "lab-123",
		Name:         "Tech 2",
		Role:         technician.RoleSeniorTechnician,
	}
	techRepo.techs["tech-3"] = &technician.Technician{
		ID:           "tech-3",
		LaboratoryID: "lab-999",
		Name:         "Tech 3",
		Role:         technician.RoleTechnician,
	}

	svc := NewService(techRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

	// List all technicians for lab-123
	techs, err := svc.ListTechnicians(context.Background(), "lab-123", nil)
	if err != nil {
		t.Errorf("ListTechnicians() unexpected error = %v", err)
		return
	}

	if len(techs) != 2 {
		t.Errorf("ListTechnicians() got %d techs, want 2", len(techs))
	}

	// List technicians filtered by role
	role := technician.RoleSeniorTechnician
	techs, err = svc.ListTechnicians(context.Background(), "lab-123", &role)
	if err != nil {
		t.Errorf("ListTechnicians() unexpected error = %v", err)
		return
	}

	if len(techs) != 1 {
		t.Errorf("ListTechnicians() with role filter got %d techs, want 1", len(techs))
	}
	if techs[0].Role != technician.RoleSeniorTechnician {
		t.Errorf("ListTechnicians() Role = %v, want %v", techs[0].Role, technician.RoleSeniorTechnician)
	}
}

func TestService_DeleteTechnician(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		labID     string
		setupRepo func(*mockTechnicianRepository)
		wantErr   error
	}{
		{
			name: "successfully delete technician",
			id:   "tech-123",
			labID: "lab-123",
			setupRepo: func(r *mockTechnicianRepository) {
				r.techs["tech-123"] = &technician.Technician{
					ID:           "tech-123",
					LaboratoryID: "lab-123",
					Name:         "John Doe",
				}
			},
			wantErr: nil,
		},
		{
			name:      "technician not found",
			id:        "non-existent",
			labID:     "lab-123",
			setupRepo: func(r *mockTechnicianRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "cross-laboratory delete",
			id:   "tech-123",
			labID: "lab-999",
			setupRepo: func(r *mockTechnicianRepository) {
				r.techs["tech-123"] = &technician.Technician{
					ID:           "tech-123",
					LaboratoryID: "lab-123",
					Name:         "John Doe",
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			techRepo := newMockTechnicianRepository()
			tt.setupRepo(techRepo)
			svc := NewService(techRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

			err := svc.DeleteTechnician(context.Background(), tt.id, tt.labID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeleteTechnician() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("DeleteTechnician() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteTechnician() unexpected error = %v", err)
			}
		})
	}
}
