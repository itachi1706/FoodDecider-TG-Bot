package model

import (
    "github.com/google/uuid"
    "time"
)

type FoodGroupsLink struct {
    ID        int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
    FoodID    uuid.UUID `gorm:"column:food_id;NOT NULL"`
    GroupID   int       `gorm:"column:group_id;NOT NULL"`
    CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
    UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
    CreatedBy int64     `gorm:"column:created_by;NOT NULL"`
    UpdatedBy int64     `gorm:"column:updated_by"`
    Status    string    `gorm:"column:status;default:A"`
}

func (m *FoodGroupsLink) TableName() string {
    return "food_groups_link"
}
