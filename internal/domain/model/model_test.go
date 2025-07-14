package model

import (
	"testing"
	"time"
)

func TestPackDistribution_Methods(t *testing.T) {
	tests := []struct {
		name          string
		distribution  PackDistribution
		expectedItems int
		expectedPacks int
		canFulfill100 bool
		isEmpty       bool
	}{
		{
			name:          "Empty distribution",
			distribution:  PackDistribution{},
			expectedItems: 0,
			expectedPacks: 0,
			canFulfill100: false,
			isEmpty:       true,
		},
		{
			name:          "Single pack",
			distribution:  PackDistribution{250: 1},
			expectedItems: 250,
			expectedPacks: 1,
			canFulfill100: true,
			isEmpty:       false,
		},
		{
			name:          "Multiple packs",
			distribution:  PackDistribution{250: 2, 500: 1},
			expectedItems: 1000,
			expectedPacks: 3,
			canFulfill100: true,
			isEmpty:       false,
		},
		{
			name:          "Edge case distribution",
			distribution:  PackDistribution{23: 2, 31: 7, 53: 9429},
			expectedItems: 23*2 + 31*7 + 53*9429,
			expectedPacks: 2 + 7 + 9429,
			canFulfill100: true,
			isEmpty:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test TotalItems
			if items := tt.distribution.TotalItems(); items != tt.expectedItems {
				t.Errorf("TotalItems: expected %d, got %d", tt.expectedItems, items)
			}

			// Test TotalPacks
			if packs := tt.distribution.TotalPacks(); packs != tt.expectedPacks {
				t.Errorf("TotalPacks: expected %d, got %d", tt.expectedPacks, packs)
			}

			// Test CanFulfill
			if canFulfill := tt.distribution.CanFulfill(100); canFulfill != tt.canFulfill100 {
				t.Errorf("CanFulfill(100): expected %v, got %v", tt.canFulfill100, canFulfill)
			}

			// Test IsEmpty
			if isEmpty := tt.distribution.IsEmpty(); isEmpty != tt.isEmpty {
				t.Errorf("IsEmpty: expected %v, got %v", tt.isEmpty, isEmpty)
			}
		})
	}
}

func TestNewCalculation(t *testing.T) {
	packSizes := []int{250, 500}
	orderQuantity := 300
	distribution := PackDistribution{500: 1}
	calculationTime := 5 * time.Millisecond

	calc := NewCalculation(packSizes, orderQuantity, distribution, calculationTime)

	if calc.OrderQuantity != orderQuantity {
		t.Errorf("Expected order quantity %d, got %d", orderQuantity, calc.OrderQuantity)
	}

	if calc.TotalItems != 500 {
		t.Errorf("Expected total items 500, got %d", calc.TotalItems)
	}

	if calc.ItemsOverage != 200 {
		t.Errorf("Expected overage 200, got %d", calc.ItemsOverage)
	}

	if calc.CalculationTime != calculationTime {
		t.Errorf("Expected calculation time %v, got %v", calculationTime, calc.CalculationTime)
	}
}

func TestNewPack(t *testing.T) {
	size := 250
	name := "Test Pack"

	pack := NewPack(size, name)

	if pack.Size != size || pack.Name != name || !pack.Active {
		t.Errorf("Pack not created correctly")
	}

	if pack.CreatedAt.IsZero() || pack.UpdatedAt.IsZero() {
		t.Errorf("Timestamps not set")
	}
}

func TestPack_Operations(t *testing.T) {
	pack := NewPack(250, "Test")

	// Test IsValid
	if !pack.IsValid() {
		t.Errorf("Valid pack should return true for IsValid()")
	}

	// Test invalid pack
	invalidPack := &Pack{Size: 0}
	if invalidPack.IsValid() {
		t.Errorf("Invalid pack should return false for IsValid()")
	}

	// Test Deactivate
	pack.Deactivate()
	if pack.Active {
		t.Errorf("Pack should be inactive after deactivation")
	}

	// Test Update
	pack.Update(500, "Updated")
	if pack.Size != 500 || pack.Name != "Updated" {
		t.Errorf("Pack not updated correctly")
	}
}

func TestPack_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		expected bool
	}{
		{"Valid pack", 250, true},
		{"Zero size", 0, false},
		{"Negative size", -10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pack := &Pack{Size: tt.size}
			if pack.IsValid() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, pack.IsValid())
			}
		})
	}
}

func TestPack_Deactivate(t *testing.T) {
	pack := NewPack(250, "Test")
	originalUpdatedAt := pack.UpdatedAt

	// Small delay to ensure time difference
	time.Sleep(time.Millisecond)
	pack.Deactivate()

	if pack.Active {
		t.Errorf("Pack should be inactive after deactivation")
	}

	if !pack.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt should be updated after deactivation")
	}
}

func TestPack_Update(t *testing.T) {
	pack := NewPack(250, "Original")
	originalUpdatedAt := pack.UpdatedAt

	newSize := 500
	newName := "Updated"

	// Small delay to ensure time difference
	time.Sleep(time.Millisecond)
	pack.Update(newSize, newName)

	if pack.Size != newSize {
		t.Errorf("Expected size %d, got %d", newSize, pack.Size)
	}

	if pack.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, pack.Name)
	}

	if !pack.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt should be updated after update")
	}
}
