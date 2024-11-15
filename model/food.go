package model

import (
	"github.com/google/uuid"
	"time"
)

type Food struct {
	ID          uuid.UUID `gorm:"column:id;primary_key"`
	Name        string    `gorm:"column:name;NOT NULL"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy   int64     `gorm:"column:created_by;NOT NULL"`
	UpdatedBy   int64     `gorm:"column:updated_by"`
	Status      string    `gorm:"column:status;default:A"`
}

func (m *Food) TableName() string {
	return "food"
}
