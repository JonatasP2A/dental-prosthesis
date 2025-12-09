package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/outbound/persistence/memory"
	orderapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/order"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
	ord "github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/order"
	"github.com/JonatasP2A/dental-prosthesis/backend/pkg/auth"
)

// mockOrderIDGenerator is a mock ID generator for testing
type mockOrderIDGenerator struct {
	id string
}

func (m *mockOrderIDGenerator) Generate() string {
	return m.id
}

func setupOrderTestRouter() (*gin.Engine, *orderapp.Service, *memory.OrderRepository, *memory.ClientRepository, *memory.LaboratoryRepository) {
	gin.SetMode(gin.TestMode)

	orderRepo := memory.NewOrderRepository()
	clientRepo := memory.NewClientRepository()
	labRepo := memory.NewLaboratoryRepository()
	idGen := &mockOrderIDGenerator{id: "test-id-123"}
	orderSvc := orderapp.NewService(orderRepo, clientRepo, idGen)
	orderHandler := NewOrderHandler(orderSvc)

	r := gin.New()
	r.POST("/orders", orderHandler.Create)
	r.GET("/orders/:id", orderHandler.Get)
	r.PUT("/orders/:id", orderHandler.Update)
	r.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
	r.GET("/orders", orderHandler.List)
	r.GET("/clients/:id/orders", orderHandler.ListByClient)
	r.DELETE("/orders/:id", orderHandler.Delete)

	return r, orderSvc, orderRepo, clientRepo, labRepo
}

func createTestLaboratoryForOrder(repo *memory.LaboratoryRepository, id string) {
	lab := &laboratory.Laboratory{
		ID:    id,
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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(nil, lab)
}

func createTestClientForOrder(repo *memory.ClientRepository, id, laboratoryID string) {
	c := &client.Client{
		ID:           id,
		LaboratoryID: laboratoryID,
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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(nil, c)
}

func createTestOrder(repo *memory.OrderRepository, id, clientID, laboratoryID string) {
	o := &ord.Order{
		ID:           id,
		ClientID:     clientID,
		LaboratoryID: laboratoryID,
		Status:       ord.StatusReceived,
		Prosthesis: []ord.ProsthesisItem{
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
	_ = repo.Create(nil, o)
}

func setLaboratoryIDInContextForOrder(req *http.Request, laboratoryID string) *http.Request {
	ctx := context.WithValue(req.Context(), auth.LaboratoryIDKey, laboratoryID)
	return req.WithContext(ctx)
}

func TestOrderHandler_Create_Success(t *testing.T) {
	router, _, _, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")

	reqBody := dto.CreateOrderRequest{
		ClientID: "client-123",
		Prosthesis: []dto.ProsthesisItemRequest{
			{
				Type:     "crown",
				Material: "zirconia",
				Shade:    "A1",
				Quantity: 1,
				Notes:    "Test notes",
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Create() status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp dto.OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ClientID != reqBody.ClientID {
		t.Errorf("Create() ClientID = %v, want %v", resp.ClientID, reqBody.ClientID)
	}
	if resp.Status != string(ord.StatusReceived) {
		t.Errorf("Create() Status = %v, want %v", resp.Status, ord.StatusReceived)
	}
}

func TestOrderHandler_Create_MissingLaboratoryID(t *testing.T) {
	router, _, _, _, _ := setupOrderTestRouter()

	reqBody := dto.CreateOrderRequest{
		ClientID: "client-123",
		Prosthesis: []dto.ProsthesisItemRequest{
			{
				Type:     "crown",
				Material: "zirconia",
				Quantity: 1,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No laboratory ID in context

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOrderHandler_Create_InvalidBody(t *testing.T) {
	router, _, _, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOrderHandler_Get_Success(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	req := httptest.NewRequest(http.MethodGet, "/orders/order-123", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp dto.OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ID != "order-123" {
		t.Errorf("Get() ID = %v, want %v", resp.ID, "order-123")
	}
}

func TestOrderHandler_Get_NotFound(t *testing.T) {
	router, _, _, _, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")

	req := httptest.NewRequest(http.MethodGet, "/orders/non-existent", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestOrderHandler_Update_Success(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	reqBody := dto.UpdateOrderRequest{
		Prosthesis: []dto.ProsthesisItemRequest{
			{
				Type:     "bridge",
				Material: "porcelain",
				Shade:    "A2",
				Quantity: 2,
				Notes:    "Updated notes",
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/orders/order-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Update() status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp dto.OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Prosthesis) != 1 {
		t.Errorf("Update() Prosthesis length = %v, want 1", len(resp.Prosthesis))
	}
	if resp.Prosthesis[0].Type != "bridge" {
		t.Errorf("Update() Prosthesis[0].Type = %v, want bridge", resp.Prosthesis[0].Type)
	}
}

func TestOrderHandler_Update_NotFound(t *testing.T) {
	router, _, _, _, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")

	reqBody := dto.UpdateOrderRequest{
		Prosthesis: []dto.ProsthesisItemRequest{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 2,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/orders/non-existent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Update() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestOrderHandler_UpdateStatus_Success(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	reqBody := dto.UpdateOrderStatusRequest{
		Status: string(ord.StatusInProduction),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/orders/order-123/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("UpdateStatus() status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp dto.OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != string(ord.StatusInProduction) {
		t.Errorf("UpdateStatus() Status = %v, want %v", resp.Status, ord.StatusInProduction)
	}
}

func TestOrderHandler_UpdateStatus_InvalidTransition(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	reqBody := dto.UpdateOrderStatusRequest{
		Status: string(ord.StatusReady), // Invalid transition from received
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/orders/order-123/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("UpdateStatus() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOrderHandler_UpdateStatus_InvalidStatus(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	reqBody := dto.UpdateOrderStatusRequest{
		Status: "invalid_status",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/orders/order-123/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("UpdateStatus() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOrderHandler_List_Success(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-1", "client-123", "lab-123")

	// Create a second order
	o2 := &ord.Order{
		ID:           "order-2",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       ord.StatusInProduction,
		Prosthesis: []ord.ProsthesisItem{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = orderRepo.Create(nil, o2)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("List() got %d orders, want 2", len(resp))
	}
}

func TestOrderHandler_ListByClient_Success(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-1", "client-123", "lab-123")

	// Create a second order for the same client
	o2 := &ord.Order{
		ID:           "order-2",
		ClientID:     "client-123",
		LaboratoryID: "lab-123",
		Status:       ord.StatusInProduction,
		Prosthesis: []ord.ProsthesisItem{
			{
				Type:     "bridge",
				Material: "porcelain",
				Quantity: 1,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = orderRepo.Create(nil, o2)

	req := httptest.NewRequest(http.MethodGet, "/clients/client-123/orders", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ListByClient() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("ListByClient() got %d orders, want 2", len(resp))
	}
}

func TestOrderHandler_Delete_Success(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	req := httptest.NewRequest(http.MethodDelete, "/orders/order-123", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Verify it's deleted
	req = httptest.NewRequest(http.MethodGet, "/orders/order-123", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() after Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestOrderHandler_Delete_NotFound(t *testing.T) {
	router, _, _, _, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")

	req := httptest.NewRequest(http.MethodDelete, "/orders/non-existent", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-123")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestOrderHandler_LaboratoryScopedAccess(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestLaboratoryForOrder(labRepo, "lab-456")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	// Try to access order from lab-123 using lab-456 context
	req := httptest.NewRequest(http.MethodGet, "/orders/order-123", nil)
	req = setLaboratoryIDInContextForOrder(req, "lab-456")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() with wrong laboratory ID status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestOrderHandler_StatusWorkflowTransitions(t *testing.T) {
	router, _, orderRepo, clientRepo, labRepo := setupOrderTestRouter()
	createTestLaboratoryForOrder(labRepo, "lab-123")
	createTestClientForOrder(clientRepo, "client-123", "lab-123")
	createTestOrder(orderRepo, "order-123", "client-123", "lab-123")

	// Test valid transitions
	transitions := []struct {
		from   ord.Status
		to     ord.Status
		valid  bool
		status int
	}{
		{ord.StatusReceived, ord.StatusInProduction, true, http.StatusOK},
		{ord.StatusInProduction, ord.StatusQualityCheck, true, http.StatusOK},
		{ord.StatusQualityCheck, ord.StatusReady, true, http.StatusOK},
		{ord.StatusReady, ord.StatusDelivered, true, http.StatusOK},
		{ord.StatusReceived, ord.StatusReady, false, http.StatusBadRequest},
		{ord.StatusDelivered, ord.StatusReady, false, http.StatusBadRequest},
	}

	for _, trans := range transitions {
		// Set order to the "from" status
		o, _ := orderRepo.GetByID(nil, "order-123")
		o.Status = trans.from
		_ = orderRepo.Update(nil, o)

		reqBody := dto.UpdateOrderStatusRequest{
			Status: string(trans.to),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPatch, "/orders/order-123/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = setLaboratoryIDInContextForOrder(req, "lab-123")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != trans.status {
			t.Errorf("UpdateStatus() from %v to %v status = %d, want %d", trans.from, trans.to, rec.Code, trans.status)
		}
	}
}
