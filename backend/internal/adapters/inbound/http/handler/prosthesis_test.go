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
	prosthesisapp "github.com/JonatasP2A/dental-prosthesis/backend/internal/application/prosthesis"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/laboratory"
	"github.com/JonatasP2A/dental-prosthesis/backend/internal/domain/prosthesis"
)

// mockIDGenerator is a mock ID generator for testing
type mockProsthesisIDGenerator struct {
	id string
}

func (m *mockProsthesisIDGenerator) Generate() string {
	return m.id
}

func setupProsthesisTestRouter() (*gin.Engine, *prosthesisapp.Service, *memory.ProsthesisRepository, *memory.LaboratoryRepository) {
	gin.SetMode(gin.TestMode)

	prosthesisRepo := memory.NewProsthesisRepository()
	labRepo := memory.NewLaboratoryRepository()
	idGen := &mockProsthesisIDGenerator{id: "test-id-123"}
	svc := prosthesisapp.NewService(prosthesisRepo, labRepo, idGen)
	handler := NewProsthesisHandler(svc)

	r := gin.New()
	r.POST("/prostheses", handler.Create)
	r.GET("/prostheses/:id", handler.Get)
	r.PUT("/prostheses/:id", handler.Update)
	r.GET("/prostheses", handler.List)
	r.DELETE("/prostheses/:id", handler.Delete)

	return r, svc, prosthesisRepo, labRepo
}

func createTestLaboratoryForProsthesis(repo *memory.LaboratoryRepository, id string) {
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

func createTestProsthesis(repo *memory.ProsthesisRepository, id, laboratoryID string) {
	p, _ := prosthesis.NewProsthesis(id, laboratoryID, prosthesis.ProsthesisTypeCrown, "zirconia", "A1", "", "")
	_ = repo.Create(nil, p)
}

// addLaboratoryIDQueryParam adds laboratory_id query parameter to request URL
func addLaboratoryIDQueryParamProsthesis(url string, laboratoryID string) string {
	if laboratoryID == "" {
		return url
	}
	separator := "?"
	if strings.Contains(url, "?") {
		separator = "&"
	}
	return url + separator + "laboratory_id=" + laboratoryID
}

func TestProsthesisHandler_Create_Success(t *testing.T) {
	router, _, _, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")

	reqBody := dto.CreateProsthesisRequest{
		Type:           "crown",
		Material:       "zirconia",
		Shade:          "A1",
		Specifications: "Full coverage",
		Notes:          "High priority",
	}

	body, _ := json.Marshal(reqBody)
	url := addLaboratoryIDQueryParamProsthesis("/prostheses", "lab-123")
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Create() status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp dto.ProsthesisResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Type != reqBody.Type {
		t.Errorf("Create() Type = %v, want %v", resp.Type, reqBody.Type)
	}
	if resp.Material != reqBody.Material {
		t.Errorf("Create() Material = %v, want %v", resp.Material, reqBody.Material)
	}
}

func TestProsthesisHandler_Create_MissingLaboratoryID(t *testing.T) {
	router, _, _, _ := setupProsthesisTestRouter()

	reqBody := dto.CreateProsthesisRequest{
		Type:     "crown",
		Material: "zirconia",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/prostheses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestProsthesisHandler_Create_InvalidType(t *testing.T) {
	router, _, _, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")

	reqBody := dto.CreateProsthesisRequest{
		Type:     "invalid_type",
		Material: "zirconia",
	}

	body, _ := json.Marshal(reqBody)
	url := addLaboratoryIDQueryParamProsthesis("/prostheses", "lab-123")
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestProsthesisHandler_Get_Success(t *testing.T) {
	router, _, prosthesisRepo, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")
	createTestProsthesis(prosthesisRepo, "prosthesis-123", "lab-123")

	url := addLaboratoryIDQueryParamProsthesis("/prostheses/prosthesis-123", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Get() status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp dto.ProsthesisResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ID != "prosthesis-123" {
		t.Errorf("Get() ID = %v, want %v", resp.ID, "prosthesis-123")
	}
}

func TestProsthesisHandler_Get_NotFound(t *testing.T) {
	router, _, _, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")

	url := addLaboratoryIDQueryParamProsthesis("/prostheses/non-existent", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Get() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestProsthesisHandler_Update_Success(t *testing.T) {
	router, _, prosthesisRepo, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")
	createTestProsthesis(prosthesisRepo, "prosthesis-123", "lab-123")

	reqBody := dto.UpdateProsthesisRequest{
		Type:           "bridge",
		Material:       "porcelain",
		Shade:          "A2",
		Specifications: "3-unit bridge",
		Notes:          "Updated",
	}

	body, _ := json.Marshal(reqBody)
	url := addLaboratoryIDQueryParamProsthesis("/prostheses/prosthesis-123", "lab-123")
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Update() status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp dto.ProsthesisResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Type != reqBody.Type {
		t.Errorf("Update() Type = %v, want %v", resp.Type, reqBody.Type)
	}
}

func TestProsthesisHandler_List_Success(t *testing.T) {
	router, _, prosthesisRepo, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")
	createTestProsthesis(prosthesisRepo, "prosthesis-1", "lab-123")

	p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
	_ = prosthesisRepo.Create(nil, p2)

	url := addLaboratoryIDQueryParamProsthesis("/prostheses", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.ProsthesisResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("List() count = %v, want %v", len(resp), 2)
	}
}

func TestProsthesisHandler_List_FilteredByType(t *testing.T) {
	router, _, prosthesisRepo, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")
	createTestProsthesis(prosthesisRepo, "prosthesis-1", "lab-123")

	p2, _ := prosthesis.NewProsthesis("prosthesis-2", "lab-123", prosthesis.ProsthesisTypeBridge, "porcelain", "", "", "")
	_ = prosthesisRepo.Create(nil, p2)

	url := addLaboratoryIDQueryParamProsthesis("/prostheses?type=crown", "lab-123")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []dto.ProsthesisResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 1 {
		t.Errorf("List() count = %v, want %v", len(resp), 1)
	}
	if resp[0].Type != "crown" {
		t.Errorf("List() Type = %v, want %v", resp[0].Type, "crown")
	}
}

func TestProsthesisHandler_Delete_Success(t *testing.T) {
	router, _, prosthesisRepo, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")
	createTestProsthesis(prosthesisRepo, "prosthesis-123", "lab-123")

	url := addLaboratoryIDQueryParamProsthesis("/prostheses/prosthesis-123", "lab-123")
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestProsthesisHandler_Delete_NotFound(t *testing.T) {
	router, _, _, labRepo := setupProsthesisTestRouter()
	createTestLaboratoryForProsthesis(labRepo, "lab-123")

	url := addLaboratoryIDQueryParamProsthesis("/prostheses/non-existent", "lab-123")
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
