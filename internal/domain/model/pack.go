package model

import "time"

// Pack represents a pack configuration with a specific size
type Pack struct {
	ID        string    `json:"id"         gorm:"primaryKey;type:varchar(255)"`
	Size      int       `json:"size"       gorm:"not null;index"`
	Name      string    `json:"name"       gorm:"type:varchar(255)"`
	Active    bool      `json:"active"     gorm:"not null;default:true;index"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// NewPack creates a new pack instance
func NewPack(size int, name string) *Pack {
	now := time.Now()
	return &Pack{
		Size:      size,
		Name:      name,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsValid checks if pack has valid configuration
func (p *Pack) IsValid() bool {
	return p.Size > 0
}

// Deactivate marks the pack as inactive
func (p *Pack) Deactivate() {
	p.Active = false
	p.UpdatedAt = time.Now()
}

// Update modifies pack properties
func (p *Pack) Update(size int, name string) {
	p.Size = size
	p.Name = name
	p.UpdatedAt = time.Now()
}
