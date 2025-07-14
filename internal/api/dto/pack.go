package dto

import (
	"time"

	"pack-calculator/internal/domain/model"
)

// CreatePackRequest represents API request for creating a pack
type CreatePackRequest struct {
	Size int    `json:"size" validate:"required,gt=0"`
	Name string `json:"name" validate:"required,min=1"`
}

// UpdatePackRequest represents API request for updating a pack
type UpdatePackRequest struct {
	Size int    `json:"size" validate:"required,gt=0"`
	Name string `json:"name" validate:"required,min=1"`
}

// PackResponse represents API response for pack operations
type PackResponse struct {
	ID        string    `json:"id"`
	Size      int       `json:"size"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PackListResponse represents API response for multiple packs
type PackListResponse struct {
	Packs []PackResponse `json:"packs"`
	Total int            `json:"total"`
}

// ToPackResponse converts domain model to API response
func ToPackResponse(pack *model.Pack) *PackResponse {
	return &PackResponse{
		ID:        pack.ID,
		Size:      pack.Size,
		Name:      pack.Name,
		Active:    pack.Active,
		CreatedAt: pack.CreatedAt,
		UpdatedAt: pack.UpdatedAt,
	}
}

// ToPackListResponse converts domain models to API response
func ToPackListResponse(packs []*model.Pack) *PackListResponse {
	responses := make([]PackResponse, len(packs))
	for i, pack := range packs {
		responses[i] = *ToPackResponse(pack)
	}

	return &PackListResponse{
		Packs: responses,
		Total: len(responses),
	}
}
