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
	techapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/technician"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/technician"
)

// mockTechIDGenerator is a mock ID generator for testing
type mockTechIDGenerator struct {
	id string
}

func (m *mockTechIDGenerator) Generate() string {
	return m.id
}

func setupTechTestRouter() (*gin.Engine, *techapp.Service, *memory.TechnicianRepository, *memory.LaboratoryRepository) {
	gin.SetMode(gin.TestMode)

	techRepo := memory.NewTechnicianRepository()
	labRepo := memory.NewLaboratoryRepository()
	idGen := &mockTechIDGenerator{id: "tech-123"}
	svc := techapp.NewService(techRepo, labRepo, idGen)
	handler := NewTechnicianHandler(svc)

	r := gin.New()
	r.POST("/technicians", handler.Create)
	r.GET("/technicians/:id", handler.Get)
	r.PUT("/technicians/:id", handler.Update)
	r.GET("/technicians", handler.List)
	r.DELETE("/technicians/:id", handler.Delete)

	return r, svc, techRepo, labRepo
}

func createTestLaboratoryForTech(repo *memory.LaboratoryRepository, id string) {
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

func createTestTechnician(repo *memory.TechnicianRepository, id, labID string) {
	tech := &technician.Technician{
		ID:           id,
		LaboratoryID: labID,
		Name:         "John Doe",
		Email:        "john@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(nil, tech)
}

func TestTechnicianHandler_Create_Success(t *testing.T) {
	router, _, _, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")

	reqBody := dto.CreateTechnicianRequest{
		Name:  "Jane Doe",
		Email: "jane@lab.com",
		Phone: "+5511888888888",
		Role:  "technician",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/technicians?laboratory_id=lab-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Create() status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp dto.TechnicianResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Name != reqBody.Name {
		t.Errorf("Create() Name = %v, want %v", resp.Name, reqBody.Name)
	}
}

func TestTechnicianHandler_Create_MissingLaboratoryID(t *testing.T) {
	router, _, _, _ := setupTechTestRouter()

	reqBody := dto.CreateTechnicianRequest{
		Name:  "Jane Doe",
		Email: "jane@lab.com",
		Phone: "+5511888888888",
		Role:  "technician",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/technicians", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTechnicianHandler_Create_InvalidRole(t *testing.T) {
	router, _, _, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")

	reqBody := dto.CreateTechnicianRequest{
		Name:  "Jane Doe",
		Email: "jane@lab.com",
		Phone: "+5511888888888",
		Role:  "invalid_role",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/technicians?laboratory_id=lab-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTechnicianHandler_Get_Success(t *testing.T) {
	router, _, techRepo, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")
	createTestTechnician(techRepo, "tech-123", "lab-123")

	req := httptest.NewRequest(http.MethodGet, "/technicians/tech-123?laboratory_id=lab-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp dto.TechnicianResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ID != "tech-123" {
		t.Errorf("Get() ID = %v, want %v", resp.ID, "tech-123")
	}
}

func TestTechnicianHandler_Get_NotFound(t *testing.T) {
	router, _, _, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")

	req := httptest.NewRequest(http.MethodGet, "/technicians/non-existent?laboratory_id=lab-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTechnicianHandler_List_Success(t *testing.T) {
	router, _, techRepo, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")
	createTestTechnician(techRepo, "tech-1", "lab-123")
	createTestTechnician(techRepo, "tech-2", "lab-123")

	req := httptest.NewRequest(http.MethodGet, "/technicians?laboratory_id=lab-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.TechnicianResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("List() got %d techs, want 2", len(resp))
	}
}

func TestTechnicianHandler_List_FilteredByRole(t *testing.T) {
	router, _, techRepo, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")
	createTestTechnician(techRepo, "tech-1", "lab-123")
	
	tech2 := &technician.Technician{
		ID:           "tech-2",
		LaboratoryID: "lab-123",
		Name:         "Senior Tech",
		Email:        "senior@lab.com",
		Phone:        "+5511999999999",
		Role:         technician.RoleSeniorTechnician,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	_ = techRepo.Create(nil, tech2)

	req := httptest.NewRequest(http.MethodGet, "/technicians?laboratory_id=lab-123&role=senior_technician", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.TechnicianResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 1 {
		t.Errorf("List() with role filter got %d techs, want 1", len(resp))
	}
	if resp[0].Role != "senior_technician" {
		t.Errorf("List() Role = %v, want %v", resp[0].Role, "senior_technician")
	}
}

func TestTechnicianHandler_Delete_Success(t *testing.T) {
	router, _, techRepo, labRepo := setupTechTestRouter()
	createTestLaboratoryForTech(labRepo, "lab-123")
	createTestTechnician(techRepo, "tech-123", "lab-123")

	req := httptest.NewRequest(http.MethodDelete, "/technicians/tech-123?laboratory_id=lab-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Verify it's deleted
	req = httptest.NewRequest(http.MethodGet, "/technicians/tech-123?laboratory_id=lab-123", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() after Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
