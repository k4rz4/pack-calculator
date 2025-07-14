package dto

import (
	"pack-calculator/internal/domain/model"
)

// CalculationRequest represents API request for pack calculation
type CalculationRequest struct {
	PackSizes     []int `json:"pack_sizes"     validate:"required,min=1,dive,gt=0"`
	OrderQuantity int   `json:"order_quantity" validate:"required,gt=0"`
}

// CalculationResponse represents API response for pack calculation
type CalculationResponse struct {
	ID              string      `json:"id"`
	PacksUsed       map[int]int `json:"packs_used"`
	TotalItems      int         `json:"total_items"`
	TotalPacks      int         `json:"total_packs"`
	ItemsOverage    int         `json:"items_overage"`
	CalculationTime string      `json:"calculation_time"`
	Success         bool        `json:"success"`
}

// SimpleCalculationRequest for calculations using stored pack configurations
type SimpleCalculationRequest struct {
	OrderQuantity int `json:"order_quantity" validate:"required,gt=0"`
}

// ToCalculationResponse converts domain model to API response
func ToCalculationResponse(result *model.Calculation) *CalculationResponse {
	return &CalculationResponse{
		ID:              result.ID,
		PacksUsed:       map[int]int(result.GetDistribution()),
		TotalItems:      result.TotalItems,
		TotalPacks:      result.TotalPacks,
		ItemsOverage:    result.ItemsOverage,
		CalculationTime: result.CalculationTime.String(),
		Success:         true,
	}
}
