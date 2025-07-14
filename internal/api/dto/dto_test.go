package dto

import (
	"testing"
	"time"

	"pack-calculator/internal/domain/model"
)

func TestToCalculationResponse(t *testing.T) {
	// Create test calculation
	packSizes := []int{250, 500}
	orderQuantity := 300
	distribution := model.PackDistribution{500: 1}
	calculationTime := 5 * time.Millisecond

	calc := model.NewCalculation(packSizes, orderQuantity, distribution, calculationTime)
	calc.ID = "test-id-123"

	// Convert to response
	response := ToCalculationResponse(calc)

	// Verify response
	if response == nil {
		t.Fatal("Response is nil")
	}

	if response.ID != calc.ID {
		t.Errorf("Expected ID %s, got %s", calc.ID, response.ID)
	}

	if response.TotalItems != calc.TotalItems {
		t.Errorf("Expected TotalItems %d, got %d", calc.TotalItems, response.TotalItems)
	}

	if response.TotalPacks != calc.TotalPacks {
		t.Errorf("Expected TotalPacks %d, got %d", calc.TotalPacks, response.TotalPacks)
	}

	if response.ItemsOverage != calc.ItemsOverage {
		t.Errorf("Expected ItemsOverage %d, got %d", calc.ItemsOverage, response.ItemsOverage)
	}

	if !response.Success {
		t.Errorf("Expected Success to be true")
	}

	// Verify packs_used conversion
	if response.PacksUsed[500] != 1 {
		t.Errorf("Expected PacksUsed[500] to be 1, got %d", response.PacksUsed[500])
	}

	// Verify calculation time string conversion
	if response.CalculationTime == "" {
		t.Errorf("CalculationTime should not be empty")
	}
}
