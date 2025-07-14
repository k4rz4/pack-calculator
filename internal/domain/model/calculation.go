package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// PackDistribution represents how many packs of each size are used
type PackDistribution map[int]int

// Calculation represents a pack calculation event
type Calculation struct {
	ID                string         `json:"id"                  gorm:"primaryKey;type:varchar(255)"`
	PackSizes         datatypes.JSON `json:"pack_sizes"          gorm:"type:jsonb"`
	OrderQuantity     int            `json:"order_quantity"      gorm:"not null;index"`
	Distribution      datatypes.JSON `json:"distribution"        gorm:"type:jsonb"`
	TotalItems        int            `json:"total_items"         gorm:"not null"`
	TotalPacks        int            `json:"total_packs"         gorm:"not null"`
	ItemsOverage      int            `json:"items_overage"       gorm:"not null"`
	CalculationTimeMs int64          `json:"calculation_time_ms" gorm:"not null"`
	CalculationTime   time.Duration  `json:"calculation_time"    gorm:"-"`
	CreatedAt         time.Time      `json:"created_at"          gorm:"not null;index"`
	UserID            string         `json:"user_id,omitempty"   gorm:"type:varchar(255);index"`
}

// NewCalculation creates a new calculation
func NewCalculation(
	packSizes []int,
	orderQuantity int,
	distribution PackDistribution,
	calculationTime time.Duration,
) *Calculation {
	totalItems := distribution.TotalItems()
	totalPacks := distribution.TotalPacks()
	overage := 0
	if totalItems > orderQuantity {
		overage = totalItems - orderQuantity
	}

	return &Calculation{
		PackSizes:         mustMarshalJSON(packSizes),
		OrderQuantity:     orderQuantity,
		Distribution:      mustMarshalJSON(distribution),
		TotalItems:        totalItems,
		TotalPacks:        totalPacks,
		ItemsOverage:      overage,
		CalculationTime:   calculationTime,
		CalculationTimeMs: calculationTime.Milliseconds(),
		CreatedAt:         time.Now(),
	}
}

// GetDistribution returns the distribution as PackDistribution
func (c *Calculation) GetDistribution() PackDistribution {
	var dist PackDistribution
	json.Unmarshal(c.Distribution, &dist)
	return dist
}

// TotalItems calculates total items in the distribution
func (pd PackDistribution) TotalItems() int {
	total := 0
	for packSize, quantity := range pd {
		total += packSize * quantity
	}
	return total
}

// TotalPacks calculates total number of packs in the distribution
func (pd PackDistribution) TotalPacks() int {
	total := 0
	for _, quantity := range pd {
		total += quantity
	}
	return total
}

// IsEmpty checks if distribution has no packs
func (pd PackDistribution) IsEmpty() bool {
	return len(pd) == 0 || pd.TotalPacks() == 0
}

// CanFulfill checks if distribution can fulfill the given order quantity
func (pd PackDistribution) CanFulfill(orderQuantity int) bool {
	return pd.TotalItems() >= orderQuantity
}

// mustMarshalJSON marshals data to JSON or panics
func mustMarshalJSON(v interface{}) datatypes.JSON {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return datatypes.JSON(data)
}
