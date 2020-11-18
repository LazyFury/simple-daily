package models

import (
	"time"

	"gorm.io/gorm"
)

// Model Model
type Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// GetObjectOrNotFound GetObjectOrNotFound
func GetObjectOrNotFound(model ModelType, query interface{}, args ...interface{}) (data interface{}, err error) {
	data = model.Object()
	if err = DB.Where(query, args...).Find(data).Error; err != nil {
		return
	}
	return
}
