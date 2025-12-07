package laboratory

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

// mockIDGenerator is a mock ID generator for testing
type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) Generate() string {
	return m.id
}

// mockRepository is a mock repository for testing
type mockRepository struct {
	labs          map[string]*laboratory.Laboratory
	getByIDErr    error
	getByEmailErr error
	createErr     error
	updateErr     error
	deleteErr     error
	listErr       error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		labs: make(map[string]*laboratory.Laboratory),
	}
}

func (m *mockRepository) Create(ctx context.Context, lab *laboratory.Laboratory) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.labs[lab.ID] = lab
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*laboratory.Laboratory, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	lab, exists := m.labs[id]
	if !exists || lab.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return lab, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*laboratory.Laboratory, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	for _, lab := range m.labs {
		if lab.Email == email && !lab.IsDeleted() {
			return lab, nil
		}
	}
	return nil, errors.ErrNotFound
}

func (m *mockRepository) Update(ctx context.Context, lab *laboratory.Laboratory) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.labs[lab.ID]; !exists {
		return errors.ErrNotFound
	}
	m.labs[lab.ID] = lab
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	lab, exists := m.labs[id]
	if !exists {
		return errors.ErrNotFound
	}
	lab.Delete()
	return nil
}

func (m *mockRepository) List(ctx context.Context) ([]*laboratory.Laboratory, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var labs []*laboratory.Laboratory
	for _, lab := range m.labs {
		if !lab.IsDeleted() {
			labs = append(labs, lab)
		}
	}
	return labs, nil
}

func TestService_CreateLaboratory(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateInput
		mockID    string
		setupRepo func(*mockRepository)
		wantErr   error
	}{
		{
			name: "successfully create laboratory",
			input: CreateInput{
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
			},
			mockID:    "lab-123",
			setupRepo: func(r *mockRepository) {},
			wantErr:   nil,
		},
		{
			name: "duplicate email",
			input: CreateInput{
				Name:  "Test Lab",
				Email: "existing@lab.com",
				Phone: "+5511999999999",
				Address: laboratory.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID: "lab-123",
			setupRepo: func(r *mockRepository) {
				r.labs["existing"] = &laboratory.Laboratory{
					ID:    "existing",
					Email: "existing@lab.com",
				}
			},
			wantErr: errors.ErrDuplicateEmail,
		},
		{
			name: "invalid input - empty name",
			input: CreateInput{
				Name:  "",
				Email: "test@lab.com",
				Phone: "+5511999999999",
				Address: laboratory.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID:    "lab-123",
			setupRepo: func(r *mockRepository) {},
			wantErr:   errors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)
			idGen := &mockIDGenerator{id: tt.mockID}
			svc := NewService(repo, idGen)

			lab, err := svc.CreateLaboratory(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateLaboratory() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					// Check if it's a validation error
					var ve errors.ValidationErrors
					if stderrors.As(err, &ve) && tt.wantErr == errors.ErrInvalidInput {
						return // ValidationErrors is considered ErrInvalidInput
					}
					t.Errorf("CreateLaboratory() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateLaboratory() unexpected error = %v", err)
				return
			}

			if lab.ID != tt.mockID {
				t.Errorf("CreateLaboratory() ID = %v, want %v", lab.ID, tt.mockID)
			}
			if lab.Name != tt.input.Name {
				t.Errorf("CreateLaboratory() Name = %v, want %v", lab.Name, tt.input.Name)
			}
		})
	}
}

func TestService_GetLaboratory(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupRepo func(*mockRepository)
		wantErr   error
	}{
		{
			name: "successfully get laboratory",
			id:   "lab-123",
			setupRepo: func(r *mockRepository) {
				r.labs["lab-123"] = &laboratory.Laboratory{
					ID:   "lab-123",
					Name: "Test Lab",
				}
			},
			wantErr: nil,
		},
		{
			name:      "laboratory not found",
			id:        "non-existent",
			setupRepo: func(r *mockRepository) {},
			wantErr:   errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)
			svc := NewService(repo, &mockIDGenerator{})

			lab, err := svc.GetLaboratory(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetLaboratory() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("GetLaboratory() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetLaboratory() unexpected error = %v", err)
				return
			}

			if lab.ID != tt.id {
				t.Errorf("GetLaboratory() ID = %v, want %v", lab.ID, tt.id)
			}
		})
	}
}

func TestService_UpdateLaboratory(t *testing.T) {
	tests := []struct {
		name      string
		input     UpdateInput
		setupRepo func(*mockRepository)
		wantErr   error
	}{
		{
			name: "successfully update laboratory",
			input: UpdateInput{
				ID:    "lab-123",
				Name:  "Updated Lab",
				Email: "updated@lab.com",
				Phone: "+5511888888888",
				Address: laboratory.Address{
					Street:     "Updated Street",
					City:       "Updated City",
					State:      "RJ",
					PostalCode: "98765-432",
					Country:    "Brazil",
				},
			},
			setupRepo: func(r *mockRepository) {
				r.labs["lab-123"] = &laboratory.Laboratory{
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
				}
			},
			wantErr: nil,
		},
		{
			name: "laboratory not found",
			input: UpdateInput{
				ID:    "non-existent",
				Name:  "Updated Lab",
				Email: "updated@lab.com",
				Phone: "+5511888888888",
				Address: laboratory.Address{
					Street:     "Updated Street",
					City:       "Updated City",
					State:      "RJ",
					PostalCode: "98765-432",
					Country:    "Brazil",
				},
			},
			setupRepo: func(r *mockRepository) {},
			wantErr:   errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)
			svc := NewService(repo, &mockIDGenerator{})

			lab, err := svc.UpdateLaboratory(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdateLaboratory() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("UpdateLaboratory() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateLaboratory() unexpected error = %v", err)
				return
			}

			if lab.Name != tt.input.Name {
				t.Errorf("UpdateLaboratory() Name = %v, want %v", lab.Name, tt.input.Name)
			}
		})
	}
}

func TestService_ListLaboratories(t *testing.T) {
	repo := newMockRepository()
	repo.labs["lab-1"] = &laboratory.Laboratory{ID: "lab-1", Name: "Lab 1"}
	repo.labs["lab-2"] = &laboratory.Laboratory{ID: "lab-2", Name: "Lab 2"}

	svc := NewService(repo, &mockIDGenerator{})

	labs, err := svc.ListLaboratories(context.Background())
	if err != nil {
		t.Errorf("ListLaboratories() unexpected error = %v", err)
		return
	}

	if len(labs) != 2 {
		t.Errorf("ListLaboratories() got %d labs, want 2", len(labs))
	}
}

func TestService_DeleteLaboratory(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupRepo func(*mockRepository)
		wantErr   error
	}{
		{
			name: "successfully delete laboratory",
			id:   "lab-123",
			setupRepo: func(r *mockRepository) {
				r.labs["lab-123"] = &laboratory.Laboratory{
					ID:   "lab-123",
					Name: "Test Lab",
				}
			},
			wantErr: nil,
		},
		{
			name:      "laboratory not found",
			id:        "non-existent",
			setupRepo: func(r *mockRepository) {},
			wantErr:   errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)
			svc := NewService(repo, &mockIDGenerator{})

			err := svc.DeleteLaboratory(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeleteLaboratory() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("DeleteLaboratory() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteLaboratory() unexpected error = %v", err)
			}
		})
	}
}
