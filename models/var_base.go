package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Model
type Model struct {
	*gorm.Model
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}