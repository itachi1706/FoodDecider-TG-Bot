package model

import (
	"FoodDecider-TG-Bot/constants"
	"github.com/google/uuid"
	"time"
)

type Rolls struct {
	ID                uuid.UUID              `gorm:"column:id;primary_key"`
	Type              constants.DecisionType `gorm:"column:type;NOT NULL"`
	ChatId            int64                  `gorm:"column:chat_id;NOT NULL"`
	Latitude          float64                `gorm:"column:latitude"`
	Longitude         float64                `gorm:"column:longitude"`
	GroupName         string                 `gorm:"column:group_name"`
	DecidedFoodID     uuid.UUID              `gorm:"column:decided_food_id;NOT NULL"`
	DecidedLocationID *uuid.UUID             `gorm:"column:decided_location_id"`
	CreatedAt         time.Time              `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time              `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy         int64                  `gorm:"column:created_by;NOT NULL"`
	UpdatedBy         int64                  `gorm:"column:updated_by"`
}
