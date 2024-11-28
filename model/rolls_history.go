package model

import (
	"github.com/google/uuid"
	"time"
)

type RollsHistory struct {
	ID                int        `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	RollID            uuid.UUID  `gorm:"column:roll_id;NOT NULL"`
	DecidedFoodID     uuid.UUID  `gorm:"column:decided_food_id;NOT NULL"`
	DecidedLocationID *uuid.UUID `gorm:"column:decided_location_id"`
	Choices           string     `gorm:"column:choices;NOT NULL"`
	CreatedAt         time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy         int64      `gorm:"column:created_by;NOT NULL"`
	UpdatedBy         int64      `gorm:"column:updated_by"`
}

func (m *RollsHistory) TableName() string {
	return "rolls_history"
}
