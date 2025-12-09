package order

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/errors"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
)

// mockIDGenerator is a mock ID generator for testing
type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) Generate() string {
	return m.id
}

// mockOrderRepository is a mock order repository for testing
type mockOrderRepository struct {
	orders        map[string]*order.Order
	getByIDErr    error
	createErr     error
	updateErr     error
	deleteErr     error
	listErr       error
	listByClientErr error
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[string]*order.Order),
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, o *order.Order) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id string) (*order.Order, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	o, exists := m.orders[id]
	if !exists || o.IsDeleted() {
		return nil, errors.ErrNotFound
	}
	return o, nil
}

func (m *mockOrderRepository) Update(ctx context.Context, o *order.Order) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.orders[o.ID]; !exists {
		return errors.ErrNotFound
	}
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, id string, status order.Status) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	o, exists := m.orders[id]
	if !exists {
		return errors.ErrNotFound
	}
	o.Status = status
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	o, exists := m.orders[id]
	if !exists {
		return errors.ErrNotFound
	}
	o.Delete()
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, laboratoryID string) ([]*order.Order, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var orders []*order.Order
	for _, o := range m.orders {
		if o.LaboratoryID == laboratoryID && !o.IsDeleted() {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

func (m *mockOrderRepository) ListByClientID(ctx context.Context, clientID string) ([]*order.Order, error) {
	if m.listByClientErr != nil {
		return nil, m.listByClientErr
	}
	var orders []*order.Order
	for _, o := range m.orders {
		if o.ClientID == clientID && !o.IsDeleted() {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

// mockClientRepository is a mock client repository for testing
type mockClientRepository struct {
	clients    map[string]*client.Client
	getByIDErr error
}

func newMockClientRepository() *mockClientRepository {
	return &mockClientRepository{
		clients: make(map[string]*client.Client),
	}
}

func (m *mockClientRepository) Create(ctx context.Context, c *client.Client) error {
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
	return nil, errors.ErrNotFound
}

func (m *mockClientRepository) Update(ctx context.Context, c *client.Client) error {
	return nil
}

func (m *mockClientRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockClientRepository) List(ctx context.Context, laboratoryID string) ([]*client.Client, error) {
	return nil, nil
}

func TestService_CreateOrder(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateInput
		mockID    string
		setupRepo func(*mockOrderRepository, *mockClientRepository)
		wantErr   error
	}{
		{
			name: "successfully create order",
			input: CreateInput{
				ClientID:     "client-123",
				LaboratoryID: "lab-123",
				Prosthesis: []order.ProsthesisItem{
					{
						Type:     "crown",
						Material: "zirconia",
						Shade:    "A1",
						Quantity: 1,
					},
				},
			},
			mockID: "order-123",
			setupRepo: func(or *mockOrderRepository, cr *mockClientRepository) {
				cr.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: nil,
		},
		{
			name: "client not found",
			input: CreateInput{
				ClientID:     "non-existent",
				LaboratoryID: "lab-123",
				Prosthesis: []order.ProsthesisItem{
					{
						Type:     "crown",
						Material: "zirconia",
						Quantity: 1,
					},
				},
			},
			mockID:    "order-123",
			setupRepo: func(or *mockOrderRepository, cr *mockClientRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "client belongs to different laboratory",
			input: CreateInput{
				ClientID:     "client-123",
				LaboratoryID: "lab-456",
				Prosthesis: []order.ProsthesisItem{
					{
						Type:     "crown",
						Material: "zirconia",
						Quantity: 1,
					},
				},
			},
			mockID: "order-123",
			setupRepo: func(or *mockOrderRepository, cr *mockClientRepository) {
				cr.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: errors.ErrNotFound,
		},
		{
			name: "invalid prosthesis items",
			input: CreateInput{
				ClientID:     "client-123",
				LaboratoryID: "lab-123",
				Prosthesis:   []order.ProsthesisItem{},
			},
			mockID: "order-123",
			setupRepo: func(or *mockOrderRepository, cr *mockClientRepository) {
				cr.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr: errors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newMockOrderRepository()
			clientRepo := newMockClientRepository()
			tt.setupRepo(orderRepo, clientRepo)
			idGen := &mockIDGenerator{id: tt.mockID}
			svc := NewService(orderRepo, clientRepo, idGen)

			o, err := svc.CreateOrder(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CreateOrder() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					// Check if it's a validation error
					var ve errors.ValidationErrors
					if stderrors.As(err, &ve) && tt.wantErr == errors.ErrInvalidInput {
						return // ValidationErrors is considered ErrInvalidInput
					}
					t.Errorf("CreateOrder() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateOrder() unexpected error = %v", err)
				return
			}

			if o.ID != tt.mockID {
				t.Errorf("CreateOrder() ID = %v, want %v", o.ID, tt.mockID)
			}
			if o.Status != order.StatusReceived {
				t.Errorf("CreateOrder() Status = %v, want %v", o.Status, order.StatusReceived)
			}
		})
	}
}

func TestService_GetOrder(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		laboratoryID string
		setupRepo    func(*mockOrderRepository)
		wantErr      error
	}{
		{
			name:         "successfully get order",
			id:           "order-123",
			laboratoryID: "lab-123",
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: nil,
		},
		{
			name:         "order not found",
			id:           "non-existent",
			laboratoryID: "lab-123",
			setupRepo:    func(r *mockOrderRepository) {},
			wantErr:      errors.ErrNotFound,
		},
		{
			name:         "order belongs to different laboratory",
			id:           "order-123",
			laboratoryID: "lab-456",
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newMockOrderRepository()
			tt.setupRepo(orderRepo)
			svc := NewService(orderRepo, newMockClientRepository(), &mockIDGenerator{})

			o, err := svc.GetOrder(context.Background(), tt.id, tt.laboratoryID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetOrder() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("GetOrder() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetOrder() unexpected error = %v", err)
				return
			}

			if o.ID != tt.id {
				t.Errorf("GetOrder() ID = %v, want %v", o.ID, tt.id)
			}
		})
	}
}

func TestService_UpdateOrder(t *testing.T) {
	tests := []struct {
		name      string
		input     UpdateInput
		setupRepo func(*mockOrderRepository)
		wantErr   error
	}{
		{
			name: "successfully update order",
			input: UpdateInput{
				ID:           "order-123",
				LaboratoryID: "lab-123",
				Prosthesis: []order.ProsthesisItem{
					{
						Type:     "bridge",
						Material: "porcelain",
						Shade:    "A2",
						Quantity: 2,
					},
				},
			},
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: nil,
		},
		{
			name: "order not found",
			input: UpdateInput{
				ID:           "non-existent",
				LaboratoryID: "lab-123",
				Prosthesis: []order.ProsthesisItem{
					{
						Type:     "bridge",
						Material: "porcelain",
						Quantity: 2,
					},
				},
			},
			setupRepo: func(r *mockOrderRepository) {},
			wantErr:   errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newMockOrderRepository()
			tt.setupRepo(orderRepo)
			svc := NewService(orderRepo, newMockClientRepository(), &mockIDGenerator{})

			o, err := svc.UpdateOrder(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdateOrder() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("UpdateOrder() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateOrder() unexpected error = %v", err)
				return
			}

			if len(o.Prosthesis) != len(tt.input.Prosthesis) {
				t.Errorf("UpdateOrder() Prosthesis length = %v, want %v", len(o.Prosthesis), len(tt.input.Prosthesis))
			}
		})
	}
}

func TestService_UpdateOrderStatus(t *testing.T) {
	tests := []struct {
		name      string
		input     UpdateStatusInput
		setupRepo func(*mockOrderRepository)
		wantErr   error
	}{
		{
			name: "successfully update status - valid transition",
			input: UpdateStatusInput{
				ID:           "order-123",
				LaboratoryID: "lab-123",
				Status:       order.StatusInProduction,
			},
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: nil,
		},
		{
			name: "invalid status transition",
			input: UpdateStatusInput{
				ID:           "order-123",
				LaboratoryID: "lab-123",
				Status:       order.StatusReady,
			},
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: errors.ErrInvalidStatusTransition,
		},
		{
			name: "order not found",
			input: UpdateStatusInput{
				ID:           "non-existent",
				LaboratoryID: "lab-123",
				Status:       order.StatusInProduction,
			},
			setupRepo: func(r *mockOrderRepository) {},
			wantErr:   errors.ErrNotFound,
		},
		{
			name: "order belongs to different laboratory",
			input: UpdateStatusInput{
				ID:           "order-123",
				LaboratoryID: "lab-456",
				Status:       order.StatusInProduction,
			},
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newMockOrderRepository()
			tt.setupRepo(orderRepo)
			svc := NewService(orderRepo, newMockClientRepository(), &mockIDGenerator{})

			o, err := svc.UpdateOrderStatus(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("UpdateOrderStatus() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("UpdateOrderStatus() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateOrderStatus() unexpected error = %v", err)
				return
			}

			if o.Status != tt.input.Status {
				t.Errorf("UpdateOrderStatus() Status = %v, want %v", o.Status, tt.input.Status)
			}
		})
	}
}

func TestService_ListOrders(t *testing.T) {
	orderRepo := newMockOrderRepository()
	orderRepo.orders["order-1"] = &order.Order{
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
	}
	orderRepo.orders["order-2"] = &order.Order{
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
	}
	orderRepo.orders["order-3"] = &order.Order{
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
	}

	svc := NewService(orderRepo, newMockClientRepository(), &mockIDGenerator{})

	orders, err := svc.ListOrders(context.Background(), "lab-123")
	if err != nil {
		t.Errorf("ListOrders() unexpected error = %v", err)
		return
	}

	if len(orders) != 2 {
		t.Errorf("ListOrders() got %d orders, want 2", len(orders))
	}
}

func TestService_ListOrdersByClient(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		laboratoryID string
		setupRepo    func(*mockOrderRepository, *mockClientRepository)
		wantErr      error
		wantCount    int
	}{
		{
			name:         "successfully list orders by client",
			clientID:     "client-123",
			laboratoryID: "lab-123",
			setupRepo: func(or *mockOrderRepository, cr *mockClientRepository) {
				cr.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
				or.orders["order-1"] = &order.Order{
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
				}
				or.orders["order-2"] = &order.Order{
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
				}
			},
			wantErr:   nil,
			wantCount: 2,
		},
		{
			name:         "client not found",
			clientID:     "non-existent",
			laboratoryID: "lab-123",
			setupRepo:    func(or *mockOrderRepository, cr *mockClientRepository) {},
			wantErr:      errors.ErrNotFound,
			wantCount:    0,
		},
		{
			name:         "client belongs to different laboratory",
			clientID:     "client-123",
			laboratoryID: "lab-456",
			setupRepo: func(or *mockOrderRepository, cr *mockClientRepository) {
				cr.clients["client-123"] = &client.Client{
					ID:           "client-123",
					LaboratoryID: "lab-123",
					Name:         "Test Client",
				}
			},
			wantErr:   errors.ErrNotFound,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newMockOrderRepository()
			clientRepo := newMockClientRepository()
			tt.setupRepo(orderRepo, clientRepo)
			svc := NewService(orderRepo, clientRepo, &mockIDGenerator{})

			orders, err := svc.ListOrdersByClient(context.Background(), tt.clientID, tt.laboratoryID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ListOrdersByClient() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("ListOrdersByClient() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ListOrdersByClient() unexpected error = %v", err)
				return
			}

			if len(orders) != tt.wantCount {
				t.Errorf("ListOrdersByClient() got %d orders, want %d", len(orders), tt.wantCount)
			}
		})
	}
}

func TestService_DeleteOrder(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		laboratoryID string
		setupRepo    func(*mockOrderRepository)
		wantErr      error
	}{
		{
			name:         "successfully delete order",
			id:           "order-123",
			laboratoryID: "lab-123",
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: nil,
		},
		{
			name:         "order not found",
			id:           "non-existent",
			laboratoryID: "lab-123",
			setupRepo:    func(r *mockOrderRepository) {},
			wantErr:      errors.ErrNotFound,
		},
		{
			name:         "order belongs to different laboratory",
			id:           "order-123",
			laboratoryID: "lab-456",
			setupRepo: func(r *mockOrderRepository) {
				r.orders["order-123"] = &order.Order{
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
				}
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := newMockOrderRepository()
			tt.setupRepo(orderRepo)
			svc := NewService(orderRepo, newMockClientRepository(), &mockIDGenerator{})

			err := svc.DeleteOrder(context.Background(), tt.id, tt.laboratoryID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("DeleteOrder() expected error %v, got nil", tt.wantErr)
					return
				}
				if !stderrors.Is(err, tt.wantErr) {
					t.Errorf("DeleteOrder() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteOrder() unexpected error = %v", err)
			}
		})
	}
}
