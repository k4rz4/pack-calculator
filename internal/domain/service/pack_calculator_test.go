package service

import (
	"testing"

	"pack-calculator/internal/domain/model"
)

func TestPackCalculator_Calculate(t *testing.T) {
	calculator := NewPackCalculator()

	tests := []struct {
		name          string
		packSizes     []int
		orderQuantity int
		expected      model.PackDistribution
		expectError   bool
	}{
		{
			name:          "Single item",
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			orderQuantity: 1,
			expected:      model.PackDistribution{250: 1},
		},
		{
			name:          "Exact match",
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			orderQuantity: 250,
			expected:      model.PackDistribution{250: 1},
		},
		{
			name:          "Need larger pack",
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			orderQuantity: 251,
			expected:      model.PackDistribution{500: 1},
		},
		{
			name:          "Complex combination",
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			orderQuantity: 501,
			expected:      model.PackDistribution{250: 1, 500: 1},
		},
		{
			name:          "Large order",
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			orderQuantity: 12001,
			expected:      model.PackDistribution{250: 1, 2000: 1, 5000: 2},
		},
		{
			name:          "Edge case - specific requirement",
			packSizes:     []int{23, 31, 53},
			orderQuantity: 500000,
			expected:      model.PackDistribution{23: 2, 31: 7, 53: 9429},
		},
		// Error cases
		{
			name:          "Empty pack sizes",
			packSizes:     []int{},
			orderQuantity: 100,
			expectError:   true,
		},
		{
			name:          "Zero order quantity",
			packSizes:     []int{250, 500},
			orderQuantity: 0,
			expectError:   true,
		},
		{
			name:          "Invalid pack size",
			packSizes:     []int{250, 0, 500},
			orderQuantity: 100,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Calculate(tt.packSizes, tt.orderQuantity)

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

			// Verify the result matches expected
			for size, expectedCount := range tt.expected {
				if result[size] != expectedCount {
					t.Errorf(
						"For pack size %d, expected %d but got %d",
						size,
						expectedCount,
						result[size],
					)
				}
			}

			// Verify total items fulfills order
			totalItems := result.TotalItems()
			if totalItems < tt.orderQuantity {
				t.Errorf(
					"Total items %d is less than order quantity %d",
					totalItems,
					tt.orderQuantity,
				)
			}

			t.Logf("✓ %s: %v (total: %d items)", tt.name, result, totalItems)
		})
	}
}

func BenchmarkPackCalculator_EdgeCase(b *testing.B) {
	calculator := NewPackCalculator()
	packSizes := []int{23, 31, 53}
	orderQuantity := 500000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calculator.Calculate(packSizes, orderQuantity)
		if err != nil {
			b.Fatalf("Calculate failed: %v", err)
		}
	}
}
