package model

import (
    "github.com/google/uuid"
    "time"
)

type Locations struct {
    ID        int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
    FoodID    uuid.UUID `gorm:"column:food_id;NOT NULL"`
    Name      string    `gorm:"column:name;NOT NULL"`
    Latitude  float64   `gorm:"column:latitude;NOT NULL"`
    Longitude float64   `gorm:"column:longitude;NOT NULL"`
    CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
    UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
    CreatedBy int64     `gorm:"column:created_by;NOT NULL"`
    UpdatedBy int64     `gorm:"column:updated_by"`
    Status    string    `gorm:"column:status;default:A"`
}
