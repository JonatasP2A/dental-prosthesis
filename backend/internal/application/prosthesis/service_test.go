package prosthesis

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
)

// mockIDGenerator is a mock ID generator for testing
type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) Generate() string {
	return m.id
}

// mockProsthesisRepository is a mock prosthesis repository for testing
type mockProsthesisRepository struct {
	prostheses map[string]*prosthesis.Prosthesis
	getByIDErr error
	createErr  error
	updateErr  error
	deleteErr  error
	listErr    error
	findByTypeErr error
	findByMaterialErr error
}

func newMockProsthesisRepository() *mockProsthesisRepository {
	return &mockProsthesisRepository{
		prostheses: make(map[string]*prosthesis.Prosthesis),
	}
}

func (m *mockProsthesisRepository) Create(ctx context.Context, p *prosthesis.Prosthesis) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.prostheses[p.ID] = p
	return nil
}

func (m *mockProsthesisRepository) GetByID(ctx context.Context, id string) (*prosthesis.Prosthesis, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	p, exists := m.prostheses[id]
	if !exists || p.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return p, nil
}

func (m *mockProsthesisRepository) Update(ctx context.Context, p *prosthesis.Prosthesis) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.prostheses[p.ID]; !exists {
		return errors.ErrNotFound
	}
	m.prostheses[p.ID] = p
	return nil
}

func (m *mockProsthesisRepository) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	p, exists := m.prostheses[id]
	if !exists {
		return errors.ErrNotFound
	}
	p.Delete()
	return nil
}

func (m *mockProsthesisRepository) List(ctx context.Context, laboratoryID string) ([]*prosthesis.Prosthesis, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var prostheses []*prosthesis.Prosthesis
	for _, p := range m.prostheses {
		if p.LaboratoryID == laboratoryID && !p.IsDeleted() {
			prostheses = append(prostheses, p)
		}
	}
	return prostheses, nil
}

func (m *mockProsthesisRepository) FindByType(ctx context.Context, laboratoryID string, prosthesisType prosthesis.ProsthesisType) ([]*prosthesis.Prosthesis, error) {
	if m.findByTypeErr != nil {
		return nil, m.findByTypeErr
	}
	var prostheses []*prosthesis.Prosthesis
	for _, p := range m.prostheses {
		if p.LaboratoryID == laboratoryID && p.Type == prosthesisType && !p.IsDeleted() {
			prostheses = append(prostheses, p)
		}
	}
	return prostheses, nil
}

func (m *mockProsthesisRepository) FindByMaterial(ctx context.Context, laboratoryID string, material string) ([]*prosthesis.Prosthesis, error) {
	if m.findByMaterialErr != nil {
		return nil, m.findByMaterialErr
	}
	var prostheses []*prosthesis.Prosthesis
	for _, p := range m.prostheses {
		if p.LaboratoryID == laboratoryID && p.Material == material && !p.IsDeleted() {
			prostheses = append(prostheses, p)
		}
	}
	return prostheses, nil
}

// mockLaboratoryRepository is a mock laboratory repository for testing
type mockLaboratoryRepository struct {
	labs       map[string]*laboratory.Laboratory
	getByIDErr error
}

func newMockLaboratoryRepository() *mockLaboratoryRepository {
	return &mockLaboratoryRepository{
		labs: make(map[string]*laboratory.Laboratory),
	}
}

func (m *mockLaboratoryRepository) GetByID(ctx context.Context, id string) (*laboratory.Laboratory, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	lab, exists := m.labs[id]
	if !exists || lab.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return lab, nil
}

func (m *mockLaboratoryRepository) Create(ctx context.Context, lab *laboratory.Laboratory) error {
	m.labs[lab.ID] = lab
	return nil
}

func (m *mockLaboratoryRepository) GetByEmail(ctx context.Context, email string) (*laboratory.Laboratory, error) {
	return nil, errors.ErrNotFound
}

func (m *mockLaboratoryRepository) Update(ctx context.Context, lab *laboratory.Laboratory) error {
	return nil
}

func (m *mockLaboratoryRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockLaboratoryRepository) List(ctx context.Context) ([]*laboratory.Laboratory, error) {
	return nil, nil
}

func TestService_CreateProsthesis(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateInput
		mockID    string
		setupRepo func(*mockProsthesisRepository, *mockLaboratoryRepository)
		wantErr   error
	}{
		{
			name: "successfully create prosthesis",
			input: CreateInput{
				LaboratoryID:   "lab-123",
				Type:           prosthesis.ProsthesisTypeCrown,
				Material:       "zirconia",
				Shade:          "A1",
				Specifications: "Full coverage",
				Notes:          "High priority",
			},
			mockID: "prosthesis-123",
			setupRepo: func(pr *mockProsthesisRepository, lr *mockLaboratoryRepository) {
				lr.labs["lab-123"] = &laboratory.Laboratory{
					ID:   "lab-123",
					Name: "Test Lab",
				}
			},
			wantErr: nil,
		},
		{
			name: "laboratory not found",
			input: CreateInput{
				LaboratoryID: "non-existent",
				Type:         prosthesis.ProsthesisTypeCrown,
				Material:     "zirconia",
			},
			mockID:    "prosthesis-123",
			setupRepo: func(pr *mockProsthesisRepository, lr *mockLaboratoryRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "invalid prosthesis type",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Type:         prosthesis.ProsthesisType("invalid"),
				Material:     "zirconia",
			},
			mockID: "prosthesis-123",
			setupRepo: func(pr *mockProsthesisRepository, lr *mockLaboratoryRepository) {
				lr.labs["lab-123"] = &laboratory.Laboratory{
					ID:   "lab-123",
					Name: "Test Lab",
				}
			},
			wantErr: errors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newMockProsthesisRepository()
			lr := newMockLaboratoryRepository()
			tt.setupRepo(pr, lr)

			idGen := &mockIDGenerator{id: tt.mockID}
			svc := NewService(pr, lr, idGen)

			ctx := context.Background()
			p, err := svc.CreateProsthesis(ctx, tt.input)

			if !stderrors.Is(err, tt.wantErr) {
				t.Errorf("CreateProsthesis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if p == nil {
					t.Fatal("CreateProsthesis() returned nil prosthesis")
				}
				if p.ID != tt.mockID {
					t.Errorf("CreateProsthesis() ID = %v, want %v", p.ID, tt.mockID)
				}
				if p.LaboratoryID != tt.input.LaboratoryID {
					t.Errorf("CreateProsthesis() LaboratoryID = %v, want %v", p.LaboratoryID, tt.input.LaboratoryID)
				}
				if p.Type != tt.input.Type {
					t.Errorf("CreateProsthesis() Type = %v, want %v", p.Type, tt.input.Type)
				}
				if p.Material != tt.input.Material {
					t.Errorf("CreateProsthesis() Material = %v, want %v", p.Material, tt.input.Material)
				}
			}
		})
	}
}

func TestService_GetProsthesis(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		labID     string
		setupRepo func(*mockProsthesisRepository)
		wantErr   error
	}{
		{
			name: "successfully get prosthesis",
			id:   "prosthesis-123",
			labID: "lab-123",
			setupRepo: func(pr *mockProsthesisRepository) {
				p, _ := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				pr.prostheses["prosthesis-123"] = p
			},
			wantErr: nil,
		},
		{
			name: "prosthesis not found",
			id:   "non-existent",
			labID: "lab-123",
			setupRepo: func(pr *mockProsthesisRepository) {},
			wantErr: errors.ErrNotFound,
		},
		{
			name: "cross-laboratory access",
			id:   "prosthesis-123",
			labID: "lab-999",
			setupRepo: func(pr *mockProsthesisRepository) {
				p, _ := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				pr.prostheses["prosthesis-123"] = p
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newMockProsthesisRepository()
			lr := newMockLaboratoryRepository()
			tt.setupRepo(pr)

			idGen := &mockIDGenerator{}
			svc := NewService(pr, lr, idGen)

			ctx := context.Background()
			p, err := svc.GetProsthesis(ctx, tt.id, tt.labID)

			if !stderrors.Is(err, tt.wantErr) {
				t.Errorf("GetProsthesis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if p == nil {
					t.Fatal("GetProsthesis() returned nil prosthesis")
				}
				if p.ID != tt.id {
					t.Errorf("GetProsthesis() ID = %v, want %v", p.ID, tt.id)
				}
			}
		})
	}
}

func TestService_UpdateProsthesis(t *testing.T) {
	tests := []struct {
		name      string
		input     UpdateInput
		setupRepo func(*mockProsthesisRepository)
		wantErr   error
	}{
		{
			name: "successfully update prosthesis",
			input: UpdateInput{
				ID:             "prosthesis-123",
				LaboratoryID:   "lab-123",
				Type:           prosthesis.ProsthesisTypeBridge,
				Material:       "porcelain",
				Shade:          "A2",
				Specifications: "3-unit bridge",
				Notes:          "Updated",
			},
			setupRepo: func(pr *mockProsthesisRepository) {
				p, _ := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				pr.prostheses["prosthesis-123"] = p
			},
			wantErr: nil,
		},
		{
			name: "prosthesis not found",
			input: UpdateInput{
				ID:           "non-existent",
				LaboratoryID: "lab-123",
				Type:         prosthesis.ProsthesisTypeCrown,
				Material:     "zirconia",
			},
			setupRepo: func(pr *mockProsthesisRepository) {},
			wantErr: errors.ErrNotFound,
		},
		{
			name: "cross-laboratory update",
			input: UpdateInput{
				ID:           "prosthesis-123",
				LaboratoryID: "lab-999",
				Type:         prosthesis.ProsthesisTypeCrown,
				Material:     "zirconia",
			},
			setupRepo: func(pr *mockProsthesisRepository) {
				p, _ := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				pr.prostheses["prosthesis-123"] = p
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newMockProsthesisRepository()
			lr := newMockLaboratoryRepository()
			tt.setupRepo(pr)

			idGen := &mockIDGenerator{}
			svc := NewService(pr, lr, idGen)

			ctx := context.Background()
			p, err := svc.UpdateProsthesis(ctx, tt.input)

			if !stderrors.Is(err, tt.wantErr) {
				t.Errorf("UpdateProsthesis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if p == nil {
					t.Fatal("UpdateProsthesis() returned nil prosthesis")
				}
				if p.Type != tt.input.Type {
					t.Errorf("UpdateProsthesis() Type = %v, want %v", p.Type, tt.input.Type)
				}
				if p.Material != tt.input.Material {
					t.Errorf("UpdateProsthesis() Material = %v, want %v", p.Material, tt.input.Material)
				}
			}
		})
	}
}

func TestService_ListProstheses(t *testing.T) {
	tests := []struct {
		name           string
		laboratoryID   string
		prosthesisType *prosthesis.ProsthesisType
		material       *string
		setupRepo      func(*mockProsthesisRepository)
		wantCount      int
		wantErr        error
	}{
		{
			name:         "list all prostheses",
			laboratoryID: "lab-123",
			setupRepo: func(pr *mockProsthesisRepository) {
				p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
				pr.prostheses["prosthesis-1"] = p1
				pr.prostheses["prosthesis-2"] = p2
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:         "list filtered by type",
			laboratoryID: "lab-123",
			prosthesisType: func() *prosthesis.ProsthesisType {
				t := prosthesis.ProsthesisTypeCrown
				return &t
			}(),
			setupRepo: func(pr *mockProsthesisRepository) {
				p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
				pr.prostheses["prosthesis-1"] = p1
				pr.prostheses["prosthesis-2"] = p2
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:         "list filtered by material",
			laboratoryID: "lab-123",
			material:     func() *string { m := "zirconia"; return &m }(),
			setupRepo: func(pr *mockProsthesisRepository) {
				p1, _ := prosthesis.NewProsthesis("prosthesis-1", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
				pr.prostheses["prosthesis-1"] = p1
				pr.prostheses["prosthesis-2"] = p2
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:         "empty list",
			laboratoryID: "lab-123",
			setupRepo:    func(pr *mockProsthesisRepository) {},
			wantCount:    0,
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newMockProsthesisRepository()
			lr := newMockLaboratoryRepository()
			tt.setupRepo(pr)

			idGen := &mockIDGenerator{}
			svc := NewService(pr, lr, idGen)

			ctx := context.Background()
			prostheses, err := svc.ListProstheses(ctx, tt.laboratoryID, tt.prosthesisType, tt.material)

			if !stderrors.Is(err, tt.wantErr) {
				t.Errorf("ListProstheses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(prostheses) != tt.wantCount {
				t.Errorf("ListProstheses() count = %v, want %v", len(prostheses), tt.wantCount)
			}
		})
	}
}

func TestService_DeleteProsthesis(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		labID     string
		setupRepo func(*mockProsthesisRepository)
		wantErr   error
	}{
		{
			name: "successfully delete prosthesis",
			id:   "prosthesis-123",
			labID: "lab-123",
			setupRepo: func(pr *mockProsthesisRepository) {
				p, _ := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				pr.prostheses["prosthesis-123"] = p
			},
			wantErr: nil,
		},
		{
			name: "prosthesis not found",
			id:   "non-existent",
			labID: "lab-123",
			setupRepo: func(pr *mockProsthesisRepository) {},
			wantErr: errors.ErrNotFound,
		},
		{
			name: "cross-laboratory delete",
			id:   "prosthesis-123",
			labID: "lab-999",
			setupRepo: func(pr *mockProsthesisRepository) {
				p, _ := prosthesis.NewProsthesis("prosthesis-123", "lab-123", prosthesis.ProsthesisTypeCrown, "zirconia", "", "", "")
				pr.prostheses["prosthesis-123"] = p
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newMockProsthesisRepository()
			lr := newMockLaboratoryRepository()
			tt.setupRepo(pr)

			idGen := &mockIDGenerator{}
			svc := NewService(pr, lr, idGen)

			ctx := context.Background()
			err := svc.DeleteProsthesis(ctx, tt.id, tt.labID)

			if !stderrors.Is(err, tt.wantErr) {
				t.Errorf("DeleteProsthesis() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				p := pr.prostheses[tt.id]
				if p == nil || !p.IsDeleted() {
					t.Error("DeleteProsthesis() prosthesis should be soft-deleted")
				}
			}
		})
	}
}
