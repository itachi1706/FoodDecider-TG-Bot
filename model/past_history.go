package model

import (
	"github.com/google/uuid"
	"time"
)

type PastHistory struct {
	ID        int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	UserID    uuid.UUID `gorm:"column:user_id;NOT NULL"`
	Username  string    `gorm:"column:username;NOT NULL"`
	FullName  string    `gorm:"column:full_name;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (m *PastHistory) TableName() string {
	return "past_history"
}
