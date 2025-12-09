package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/outbound/persistence/memory"
	clientapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/client"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

// mockIDGenerator is a mock ID generator for testing
type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) Generate() string {
	return m.id
}

func setupTestRouter() (*gin.Engine, *clientapp.Service, *memory.ClientRepository, *memory.LaboratoryRepository) {
	gin.SetMode(gin.TestMode)

	clientRepo := memory.NewClientRepository()
	labRepo := memory.NewLaboratoryRepository()
	idGen := &mockIDGenerator{id: "test-id-123"}
	svc := clientapp.NewService(clientRepo, labRepo, idGen)
	handler := NewClientHandler(svc)

	r := gin.New()
	r.POST("/clients", handler.Create)
	r.GET("/clients/:id", handler.Get)
	r.PUT("/clients/:id", handler.Update)
	r.GET("/clients", handler.List)
	r.DELETE("/clients/:id", handler.Delete)

	return r, svc, clientRepo, labRepo
}

func createTestLaboratory(repo *memory.LaboratoryRepository, id string) {
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

func createTestClient(repo *memory.ClientRepository, id, laboratoryID string) {
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

// addLaboratoryIDQueryParam adds laboratory_id query parameter to request URL
func addLaboratoryIDQueryParam(url string, laboratoryID string) string {
	if laboratoryID == "" {
		return url
	}
	separator := "?"
	if strings.Contains(url, "?") {
		separator = "&"
	}
	return url + separator + "laboratory_id=" + laboratoryID
}

func TestClientHandler_Create_Success(t *testing.T) {
	router, _, _, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")

	reqBody := dto.CreateClientRequest{
		Name:  "New Client",
		Email: "new@example.com",
		Phone: "+5511888888888",
		Address: dto.ClientAddressRequest{
			Street:     "New Street",
			City:       "New City",
			State:      "RJ",
			PostalCode: "98765-432",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	url := addLaboratoryIDQueryParam("/clients", "lab-123")
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Create() status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp dto.ClientResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Name != reqBody.Name {
		t.Errorf("Create() Name = %v, want %v", resp.Name, reqBody.Name)
	}
}

func TestClientHandler_Create_MissingLaboratoryID(t *testing.T) {
	router, _, _, _ := setupTestRouter()

	reqBody := dto.CreateClientRequest{
		Name:  "New Client",
		Email: "new@example.com",
		Phone: "+5511888888888",
		Address: dto.ClientAddressRequest{
			Street:     "New Street",
			City:       "New City",
			State:      "RJ",
			PostalCode: "98765-432",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No laboratory_id query parameter

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestClientHandler_Create_InvalidBody(t *testing.T) {
	router, _, _, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")

	url := addLaboratoryIDQueryParam("/clients", "lab-123")
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestClientHandler_Get_Success(t *testing.T) {
	router, _, clientRepo, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")
	createTestClient(clientRepo, "client-123", "lab-123")

	url := addLaboratoryIDQueryParam("/clients/client-123", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp dto.ClientResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ID != "client-123" {
		t.Errorf("Get() ID = %v, want %v", resp.ID, "client-123")
	}
}

func TestClientHandler_Get_NotFound(t *testing.T) {
	router, _, _, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")

	url := addLaboratoryIDQueryParam("/clients/non-existent", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestClientHandler_Update_Success(t *testing.T) {
	router, _, clientRepo, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")
	createTestClient(clientRepo, "client-123", "lab-123")

	reqBody := dto.UpdateClientRequest{
		Name:  "Updated Client",
		Email: "updated@example.com",
		Phone: "+5511777777777",
		Address: dto.ClientAddressRequest{
			Street:     "Updated Street",
			City:       "Updated City",
			State:      "MG",
			PostalCode: "11111-111",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	url := addLaboratoryIDQueryParam("/clients/client-123", "lab-123")
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Update() status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp dto.ClientResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Name != reqBody.Name {
		t.Errorf("Update() Name = %v, want %v", resp.Name, reqBody.Name)
	}
}

func TestClientHandler_Update_NotFound(t *testing.T) {
	router, _, _, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")

	reqBody := dto.UpdateClientRequest{
		Name:  "Updated Client",
		Email: "updated@example.com",
		Phone: "+5511777777777",
		Address: dto.ClientAddressRequest{
			Street:     "Updated Street",
			City:       "Updated City",
			State:      "MG",
			PostalCode: "11111-111",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	url := addLaboratoryIDQueryParam("/clients/non-existent", "lab-123")
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Update() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestClientHandler_List_Success(t *testing.T) {
	router, _, clientRepo, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")
	createTestClient(clientRepo, "client-1", "lab-123")

	// Create a second client
	c2 := &client.Client{
		ID:           "client-2",
		LaboratoryID: "lab-123",
		Name:         "Test Client 2",
		Email:        "test2@example.com",
		Phone:        "+5511888888888",
		Address: client.Address{
			Street:     "Test Street 2",
			City:       "Test City",
			State:      "SP",
			PostalCode: "01234-567",
			Country:    "Brazil",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = clientRepo.Create(nil, c2)

	url := addLaboratoryIDQueryParam("/clients", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.ClientResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("List() got %d clients, want 2", len(resp))
	}
}

func TestClientHandler_List_Empty(t *testing.T) {
	router, _, _, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")

	url := addLaboratoryIDQueryParam("/clients", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.ClientResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 0 {
		t.Errorf("List() got %d clients, want 0", len(resp))
	}
}

func TestClientHandler_Delete_Success(t *testing.T) {
	router, _, clientRepo, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")
	createTestClient(clientRepo, "client-123", "lab-123")

	url := addLaboratoryIDQueryParam("/clients/client-123", "lab-123")
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Verify it's deleted
	url = addLaboratoryIDQueryParam("/clients/client-123", "lab-123")
	req = httptest.NewRequest(http.MethodGet, url, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() after Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestClientHandler_Delete_NotFound(t *testing.T) {
	router, _, _, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")

	url := addLaboratoryIDQueryParam("/clients/non-existent", "lab-123")
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestClientHandler_LaboratoryScopedAccess(t *testing.T) {
	router, _, clientRepo, labRepo := setupTestRouter()
	createTestLaboratory(labRepo, "lab-123")
	createTestLaboratory(labRepo, "lab-456")
	createTestClient(clientRepo, "client-123", "lab-123")

	// Try to access client from lab-123 using lab-456 query parameter
	url := addLaboratoryIDQueryParam("/clients/client-123", "lab-456")
	req := httptest.NewRequest(http.MethodGet, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() with wrong laboratory ID status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
