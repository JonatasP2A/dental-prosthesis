package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/inbound/http/dto"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/adapters/outbound/persistence/memory"
	labapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
)

// mockIDGenerator is a mock ID generator for testing
type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) Generate() string {
	return m.id
}

func setupTestRouter() (*gin.Engine, *labapp.Service, *memory.LaboratoryRepository) {
	gin.SetMode(gin.TestMode)

	repo := memory.NewLaboratoryRepository()
	idGen := &mockIDGenerator{id: "test-id-123"}
	svc := labapp.NewService(repo, idGen)
	handler := NewLaboratoryHandler(svc)

	r := gin.New()
	r.POST("/laboratories", handler.Create)
	r.GET("/laboratories/:id", handler.Get)
	r.PUT("/laboratories/:id", handler.Update)
	r.GET("/laboratories", handler.List)
	r.DELETE("/laboratories/:id", handler.Delete)

	return r, svc, repo
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

func TestLaboratoryHandler_Create_Success(t *testing.T) {
	router, _, _ := setupTestRouter()

	reqBody := dto.CreateLaboratoryRequest{
		Name:  "New Lab",
		Email: "new@lab.com",
		Phone: "+5511888888888",
		Address: dto.AddressRequest{
			Street:     "New Street",
			City:       "New City",
			State:      "RJ",
			PostalCode: "98765-432",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/laboratories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Create() status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp dto.LaboratoryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Name != reqBody.Name {
		t.Errorf("Create() Name = %v, want %v", resp.Name, reqBody.Name)
	}
}

func TestLaboratoryHandler_Create_InvalidBody(t *testing.T) {
	router, _, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/laboratories", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLaboratoryHandler_Get_Success(t *testing.T) {
	router, _, repo := setupTestRouter()
	createTestLaboratory(repo, "lab-123")

	req := httptest.NewRequest(http.MethodGet, "/laboratories/lab-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp dto.LaboratoryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ID != "lab-123" {
		t.Errorf("Get() ID = %v, want %v", resp.ID, "lab-123")
	}
}

func TestLaboratoryHandler_Get_NotFound(t *testing.T) {
	router, _, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/laboratories/non-existent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestLaboratoryHandler_Update_Success(t *testing.T) {
	router, _, repo := setupTestRouter()
	createTestLaboratory(repo, "lab-123")

	reqBody := dto.UpdateLaboratoryRequest{
		Name:  "Updated Lab",
		Email: "updated@lab.com",
		Phone: "+5511777777777",
		Address: dto.AddressRequest{
			Street:     "Updated Street",
			City:       "Updated City",
			State:      "MG",
			PostalCode: "11111-111",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/laboratories/lab-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Update() status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp dto.LaboratoryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Name != reqBody.Name {
		t.Errorf("Update() Name = %v, want %v", resp.Name, reqBody.Name)
	}
}

func TestLaboratoryHandler_Update_NotFound(t *testing.T) {
	router, _, _ := setupTestRouter()

	reqBody := dto.UpdateLaboratoryRequest{
		Name:  "Updated Lab",
		Email: "updated@lab.com",
		Phone: "+5511777777777",
		Address: dto.AddressRequest{
			Street:     "Updated Street",
			City:       "Updated City",
			State:      "MG",
			PostalCode: "11111-111",
			Country:    "Brazil",
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/laboratories/non-existent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Update() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestLaboratoryHandler_List_Success(t *testing.T) {
	router, _, repo := setupTestRouter()
	createTestLaboratory(repo, "lab-1")

	// Create a second lab with different email
	lab2 := &laboratory.Laboratory{
		ID:    "lab-2",
		Name:  "Test Lab 2",
		Email: "test2@lab.com",
		Phone: "+5511888888888",
		Address: laboratory.Address{
			Street:     "Test Street 2",
			City:       "Test City",
			State:      "SP",
			PostalCode: "01234-567",
			Country:    "Brazil",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(nil, lab2)

	req := httptest.NewRequest(http.MethodGet, "/laboratories", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.LaboratoryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("List() got %d labs, want 2", len(resp))
	}
}

func TestLaboratoryHandler_List_Empty(t *testing.T) {
	router, _, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/laboratories", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.LaboratoryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 0 {
		t.Errorf("List() got %d labs, want 0", len(resp))
	}
}

func TestLaboratoryHandler_Delete_Success(t *testing.T) {
	router, _, repo := setupTestRouter()
	createTestLaboratory(repo, "lab-123")

	req := httptest.NewRequest(http.MethodDelete, "/laboratories/lab-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Verify it's deleted
	req = httptest.NewRequest(http.MethodGet, "/laboratories/lab-123", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() after Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestLaboratoryHandler_Delete_NotFound(t *testing.T) {
	router, _, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodDelete, "/laboratories/non-existent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
