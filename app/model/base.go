package model

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-" json:"createdAt"`
	UpdatedAt time.Time      `json:"-" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" json:"deletedAt"`
}
