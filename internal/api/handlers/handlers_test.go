package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pack-calculator/internal/api/dto"
	"pack-calculator/internal/domain/service"
)

func TestCalculationHandler_Calculate(t *testing.T) {
	packService := service.NewPackService()
	handler := NewCalculationHandler(packService)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "Valid calculation request",
			requestBody: dto.CalculationRequest{
				PackSizes:     []int{250, 500, 1000},
				OrderQuantity: 263,
			},
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name: "Edge case request",
			requestBody: dto.CalculationRequest{
				PackSizes:     []int{23, 31, 53},
				OrderQuantity: 500000,
			},
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name: "Invalid request - empty pack sizes",
			requestBody: dto.CalculationRequest{
				PackSizes:     []int{},
				OrderQuantity: 100,
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name: "Invalid request - zero order quantity",
			requestBody: dto.CalculationRequest{
				PackSizes:     []int{250, 500},
				OrderQuantity: 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.Background())

			rr := httptest.NewRecorder()
			handler.Calculate(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			success, ok := response["success"].(bool)
			if !ok {
				t.Errorf("Response missing success field")
			}

			if success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, success)
			}

			if tt.expectSuccess {
				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Errorf("Response missing data field")
					return
				}

				requiredFields := []string{"id", "packs_used", "total_items", "total_packs"}
				for _, field := range requiredFields {
					if _, exists := data[field]; !exists {
						t.Errorf("Response data missing field: %s", field)
					}
				}

				// Verify edge case result
				if tt.name == "Edge case request" {
					packsUsed, ok := data["packs_used"].(map[string]interface{})
					if !ok {
						t.Errorf("packs_used is not a map")
						return
					}

					expectedPacks := map[string]float64{"23": 2, "31": 7, "53": 9429}
					for packSize, expectedCount := range expectedPacks {
						if count, exists := packsUsed[packSize]; !exists || count != expectedCount {
							t.Errorf("Expected %s: %v, got %v", packSize, expectedCount, count)
						}
					}
				}
			}
		})
	}
}

func TestHealthHandler_Health(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response dto.HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest("GET", "/ready", nil)
	rr := httptest.NewRecorder()

	handler.Ready(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response dto.HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "ready" {
		t.Errorf("Expected status 'ready', got '%s'", response.Status)
	}
}
