package client

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
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

// mockClientRepository is a mock client repository for testing
type mockClientRepository struct {
	clients        map[string]*client.Client
	getByIDErr     error
	getByEmailErr  error
	createErr      error
	updateErr      error
	deleteErr      error
	listErr        error
}

func newMockClientRepository() *mockClientRepository {
	return &mockClientRepository{
		clients: make(map[string]*client.Client),
	}
}

func (m *mockClientRepository) Create(ctx context.Context, c *client.Client) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.clients[c.ID] = c
	return nil
}

func (m *mockClientRepository) GetByID(ctx context.Context, id string) (*client.Client, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	c, exists := m.clients[id]
	if !exists || c.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return c, nil
}

func (m *mockClientRepository) GetByEmail(ctx context.Context, laboratoryID, email string) (*client.Client, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	for _, c := range m.clients {
		if c.Email == email && c.LaboratoryID == laboratoryID && !c.IsDeleted() {
			return c, nil
		}
	}
	return nil, errors.ErrNotFound
}

func (m *mockClientRepository) Update(ctx context.Context, c *client.Client) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.clients[c.ID]; !exists {
		return errors.ErrNotFound
	}
	m.clients[c.ID] = c
	return nil
}

func (m *mockClientRepository) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	c, exists := m.clients[id]
	if !exists {
		return errors.ErrNotFound
	}
	c.Delete()
	return nil
}

func (m *mockClientRepository) List(ctx context.Context, laboratoryID string) ([]*client.Client, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var clients []*client.Client
	for _, c := range m.clients {
		if c.LaboratoryID == laboratoryID && !c.IsDeleted() {
			clients = append(clients, c)
		}
	}
	return clients, nil
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

func TestService_CreateClient(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateInput
		mockID    string
		setupRepo func(*mockClientRepository, *mockLaboratoryRepository)
		wantErr   error
	}{
		{
			name: "successfully create client",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "Test Client",
				Email:        "test@example.com",
				Phone:        "+5511999999999",
				Address: client.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID: "client-123",
			setupRepo: func(cr *mockClientRepository, lr *mockLaboratoryRepository) {
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
				Name:         "Test Client",
				Email:        "test@example.com",
				Phone:        "+5511999999999",
				Address: client.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID:    "client-123",
			setupRepo: func(cr *mockClientRepository, lr *mockLaboratoryRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "duplicate email within laboratory",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "Test Client",
				Email:        "existing@example.com",
				Phone:        "+5511999999999",
				Address: client.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID: "client-123",
			setupRepo: func(cr *mockClientRepository, lr *mockLaboratoryRepository) {
				lr.labs["lab-123"] = &laboratory.Laboratory{
					ID:   "lab-123",
					Name: "Test Lab",
				}
				cr.clients["existing"] = &client.Client{
					ID:           "existing",
					LaboratoryID: "lab-123",
					Email:        "existing@example.com",
				}
			},
			wantErr: errors.ErrDuplicateEmail,
		},
		{
			name: "duplicate email across laboratories allowed",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "Test Client",
				Email:        "existing@example.com",
				Phone:        "+5511999999999",
				Address: client.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID: "client-123",
			setupRepo: func(cr *mockClientRepository, lr *mockLaboratoryRepository) {
				lr.labs["lab-123"] = &laboratory.Laboratory{
					ID:   "lab-123",
					Name: "Test Lab",
				}
				// Client exists in different laboratory
				cr.clients["existing"] = &client.Client{
					ID:           "existing",
					LaboratoryID: "lab-456",
					Email:        "existing@example.com",
				}
			},
			wantErr: nil,
		},
		{
			name: "invalid input - empty name",
			input: CreateInput{
				LaboratoryID: "lab-123",
				Name:         "",
				Email:        "test@example.com",
				Phone:        "+5511999999999",
				Address: client.Address{
					Street:     "Test Street",
					City:       "Test City",
					State:      "SP",
					PostalCode: "01234-567",
					Country:    "Brazil",
				},
			},
			mockID: "client-123",
			setupRepo: func(cr *mockClientRepository, lr *mockLaboratoryRepository) {
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
			clientRepo := newMockClientRepository()
			labRepo := newMockLaboratoryRepository()
			tt.setupRepo(clientRepo, labRepo)
			idGen := &mockIDGenerator{id: tt.mockID}
			svc := NewService(clientRepo, labRepo, idGen)

			c, err := svc.CreateClient(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateClient() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					// Check if it's a validation error
					var ve errors.ValidationErrors
					if stderrors.As(err, &ve) && tt.wantErr == errors.ErrInvalidInput {
						return // ValidationErrors is considered ErrInvalidInput
					}
					t.Errorf("CreateClient() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateClient() unexpected error = %v", err)
				return
			}

			if c.ID != tt.mockID {
				t.Errorf("CreateClient() ID = %v, want %v", c.ID, tt.mockID)
			}
			if c.Name != tt.input.Name {
				t.Errorf("CreateClient() Name = %v, want %v", c.Name, tt.input.Name)
			}
		})
	}
}

func TestService_GetClient(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		laboratoryID string
		setupRepo    func(*mockClientRepository)
		wantErr      error
	}{
		{
			name:         "successfully get client",
			id:           "client-123",
			laboratoryID: "lab-123",
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: nil,
		},
		{
			name:         "client not found",
			id:           "non-existent",
			laboratoryID: "lab-123",
			setupRepo:    func(r *mockClientRepository) {},
			wantErr:      errors.ErrNotFound,
		},
		{
			name:         "client belongs to different laboratory",
			id:           "client-123",
			laboratoryID: "lab-456",
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientRepo := newMockClientRepository()
			tt.setupRepo(clientRepo)
			svc := NewService(clientRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

			c, err := svc.GetClient(context.Background(), tt.id, tt.laboratoryID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetClient() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("GetClient() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetClient() unexpected error = %v", err)
				return
			}

			if c.ID != tt.id {
				t.Errorf("GetClient() ID = %v, want %v", c.ID, tt.id)
			}
		})
	}
}

func TestService_UpdateClient(t *testing.T) {
	tests := []struct {
		name      string
		input     UpdateInput
		setupRepo func(*mockClientRepository)
		wantErr   error
	}{
		{
			name: "successfully update client",
			input: UpdateInput{
				ID:           "client-123",
				LaboratoryID: "lab-123",
				Name:         "Updated Client",
				Email:        "updated@example.com",
				Phone:        "+5511888888888",
				Address: client.Address{
					Street:     "Updated Street",
					City:       "Updated City",
					State:      "RJ",
					PostalCode: "98765-432",
					Country:    "Brazil",
				},
			},
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
					Email:        "test@example.com",
					Phone:        "+5511999999999",
					Address: client.Address{
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
			name: "client not found",
			input: UpdateInput{
				ID:           "non-existent",
				LaboratoryID: "lab-123",
				Name:         "Updated Client",
				Email:        "updated@example.com",
				Phone:        "+5511888888888",
				Address: client.Address{
					Street:     "Updated Street",
					City:       "Updated City",
					State:      "RJ",
					PostalCode: "98765-432",
					Country:    "Brazil",
				},
			},
			setupRepo: func(r *mockClientRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "duplicate email on update",
			input: UpdateInput{
				ID:           "client-123",
				LaboratoryID: "lab-123",
				Name:         "Updated Client",
				Email:        "existing@example.com",
				Phone:        "+5511888888888",
				Address: client.Address{
					Street:     "Updated Street",
					City:       "Updated City",
					State:      "RJ",
					PostalCode: "98765-432",
					Country:    "Brazil",
				},
			},
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
					Email:        "test@example.com",
				}
				r.clients["existing"] = &client.Client{
					ID:           "existing",
					LaboratoryID: "lab-123",
					Email:        "existing@example.com",
				}
			},
			wantErr: errors.ErrDuplicateEmail,
		},
		{
			name: "same email on update allowed",
			input: UpdateInput{
				ID:           "client-123",
				LaboratoryID: "lab-123",
				Name:         "Updated Client",
				Email:        "test@example.com", // Same email
				Phone:        "+5511888888888",
				Address: client.Address{
					Street:     "Updated Street",
					City:       "Updated City",
					State:      "RJ",
					PostalCode: "98765-432",
					Country:    "Brazil",
				},
			},
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
					Email:        "test@example.com",
					Phone:        "+5511999999999",
					Address: client.Address{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientRepo := newMockClientRepository()
			tt.setupRepo(clientRepo)
			svc := NewService(clientRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

			c, err := svc.UpdateClient(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdateClient() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("UpdateClient() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateClient() unexpected error = %v", err)
				return
			}

			if c.Name != tt.input.Name {
				t.Errorf("UpdateClient() Name = %v, want %v", c.Name, tt.input.Name)
			}
		})
	}
}

func TestService_ListClients(t *testing.T) {
	clientRepo := newMockClientRepository()
	clientRepo.clients["client-1"] = &client.Client{
		ID:           "client-1",
		LaboratoryID: "lab-123",
		Name:         "Client 1",
	}
	clientRepo.clients["client-2"] = &client.Client{
		ID:           "client-2",
		LaboratoryID: "lab-123",
		Name:         "Client 2",
	}
	clientRepo.clients["client-3"] = &client.Client{
		ID:           "client-3",
		LaboratoryID: "lab-456", // Different laboratory
		Name:         "Client 3",
	}

	svc := NewService(clientRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

	clients, err := svc.ListClients(context.Background(), "lab-123")
	if err != nil {
		t.Errorf("ListClients() unexpected error = %v", err)
		return
	}

	if len(clients) != 2 {
		t.Errorf("ListClients() got %d clients, want 2", len(clients))
	}
}

func TestService_DeleteClient(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		laboratoryID string
		setupRepo    func(*mockClientRepository)
		wantErr      error
	}{
		{
			name:         "successfully delete client",
			id:           "client-123",
			laboratoryID: "lab-123",
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: nil,
		},
		{
			name:         "client not found",
			id:           "non-existent",
			laboratoryID: "lab-123",
			setupRepo:    func(r *mockClientRepository) {},
			wantErr:      errors.ErrNotFound,
		},
		{
			name:         "client belongs to different laboratory",
			id:           "client-123",
			laboratoryID: "lab-456",
			setupRepo: func(r *mockClientRepository) {
				r.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientRepo := newMockClientRepository()
			tt.setupRepo(clientRepo)
			svc := NewService(clientRepo, newMockLaboratoryRepository(), &mockIDGenerator{})

			err := svc.DeleteClient(context.Background(), tt.id, tt.laboratoryID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeleteClient() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("DeleteClient() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteClient() unexpected error = %v", err)
			}
		})
	}
}
