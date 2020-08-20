package model

import (
	"gorm.io/gorm"
	"time"
)

// Base - Revamped the current default base GORM model to provide more flexibility when encoding to JSON
// This reduces extra DisplayWrappers in order to obtain the ID
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-" json:"createdAt"`
	UpdatedAt time.Time      `json:"-" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" json:"deletedAt"`
}
