package service

import (
	"context"
	"testing"
)

func TestPackService_CalculateOptimal(t *testing.T) {
	service := NewPackService()
	ctx := context.Background()

	tests := []struct {
		name          string
		packSizes     []int
		orderQuantity int
		expectError   bool
	}{
		{
			name:          "Valid calculation",
			packSizes:     []int{250, 500, 1000},
			orderQuantity: 263,
			expectError:   false,
		},
		{
			name:          "Edge case",
			packSizes:     []int{23, 31, 53},
			orderQuantity: 500000,
			expectError:   false,
		},
		{
			name:          "Invalid input - empty pack sizes",
			packSizes:     []int{},
			orderQuantity: 100,
			expectError:   true,
		},
		{
			name:          "Invalid input - zero quantity",
			packSizes:     []int{250, 500},
			orderQuantity: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CalculateOptimal(ctx, tt.packSizes, tt.orderQuantity)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify result structure
			if result == nil {
				t.Errorf("Result is nil")
				return
			}

			if result.ID == "" {
				t.Errorf("Result ID is empty")
			}

			if result.OrderQuantity != tt.orderQuantity {
				t.Errorf(
					"Expected order quantity %d, got %d",
					tt.orderQuantity,
					result.OrderQuantity,
				)
			}

			if result.TotalItems < tt.orderQuantity {
				t.Errorf(
					"Total items %d is less than order quantity %d",
					result.TotalItems,
					tt.orderQuantity,
				)
			}

			if result.CalculationTime <= 0 {
				t.Errorf("Calculation time should be positive, got %v", result.CalculationTime)
			}

			t.Logf("✓ %s: TotalItems=%d, TotalPacks=%d, Time=%v",
				tt.name, result.TotalItems, result.TotalPacks, result.CalculationTime)
		})
	}
}
