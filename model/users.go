package model

import (
	"github.com/google/uuid"
	"time"
)

type Users struct {
	ID         uuid.UUID `gorm:"column:id;primary_key"`
	TelegramID int64     `gorm:"column:telegram_id;NOT NULL"`
	Username   string    `gorm:"column:username;NOT NULL"`
	FullName   string    `gorm:"column:full_name;NOT NULL"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy  int64     `gorm:"column:created_by;NOT NULL"`
	UpdatedBy  int64     `gorm:"column:updated_by"`
	Version    int64     `gorm:"column:version;default:1"`
	RawData    string    `gorm:"column:raw_data"`
}
