package api

import (
	"bytes"
	"cushion-isa/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvestment(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	queue := make(chan models.Investment, 1)
	mockRepo := &mockRepository{}
	handler := NewHandler(mockRepo, queue)

	router := gin.New()
	router.POST("/invest", handler.CreateInvestment)

	// Test cases
	tests := []struct {
		name       string
		input      map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid investment",
			input: map[string]interface{}{
				"name":   "John Doe",
				"email":  "john@example.com",
				"fund":   "sustainable-growth",
				"amount": 1000.00,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid amount",
			input: map[string]interface{}{
				"name":   "John Doe",
				"email":  "john@example.com",
				"fund":   "sustainable-growth",
				"amount": -100.00,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing required fields",
			input: map[string]interface{}{
				"name": "John Doe",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/invest", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(resp, req)

			// Assert response
			assert.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
			}
		})
	}
}

// Mock repository for testing
type mockRepository struct{}

func (m *mockRepository) Create(inv *models.Investment) error {
	inv.ID = 1 // Simulate ID assignment
	return nil
}

func (m *mockRepository) GetByID(id int64) (*models.Investment, error) {
	return &models.Investment{ID: id}, nil
}

func (m *mockRepository) UpdateStatus(id int64, status string) error {
	return nil
}

func (m *mockRepository) List(limit, offset int) ([]models.Investment, error) {
	return []models.Investment{}, nil
}
