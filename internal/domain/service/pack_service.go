package service

import (
	"context"
	"fmt"
	"time"

	"pack-calculator/internal/domain/model"
	"pack-calculator/internal/infrastructure/logger"
)

type PackService struct {
	calculator *PackCalculator
}

func NewPackService() *PackService {
	return &PackService{
		calculator: NewPackCalculator(),
	}
}

func (ps *PackService) CalculateOptimal(
	ctx context.Context,
	packSizes []int,
	orderQuantity int,
) (*model.Calculation, error) {
	startTime := time.Now()

	logger.Debug("Starting pack calculation", map[string]interface{}{
		"pack_sizes":     packSizes,
		"order_quantity": orderQuantity,
	})

	// Perform calculation
	distribution, err := ps.calculator.Calculate(packSizes, orderQuantity)
	if err != nil {
		logger.Error("Pack calculation failed", map[string]interface{}{
			"pack_sizes":     packSizes,
			"order_quantity": orderQuantity,
			"error":          err.Error(),
		})
		return nil, err
	}

	calculationTime := time.Since(startTime)

	// Create result
	result := model.NewCalculation(packSizes, orderQuantity, distribution, calculationTime)
	result.ID = ps.generateID()

	logger.Debug("Pack calculation completed", map[string]interface{}{
		"calculation_id": result.ID,
		"total_items":    result.TotalItems,
		"total_packs":    result.TotalPacks,
		"items_overage":  result.ItemsOverage,
		"duration_ms":    calculationTime.Milliseconds(),
	})

	return result, nil
}

func (ps *PackService) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
