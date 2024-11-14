package model

import (
    "time"
)

type Admins struct {
    ID           int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
    TelegramID   int64     `gorm:"column:telegram_id;NOT NULL"`
    Name         string    `gorm:"column:name;NOT NULL"`
    CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
    UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
    CreatedBy    int64     `gorm:"column:created_by;NOT NULL"`
    UpdatedBy    int64     `gorm:"column:updated_by"`
    IsSuperadmin bool      `gorm:"column:is_superadmin;default:0"`
    Status       string    `gorm:"column:status;default:A"`
}
