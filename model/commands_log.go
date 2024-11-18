package model

import (
	"FoodDecider-TG-Bot/constants"
	"github.com/google/uuid"
	"time"
)

type CommandsLog struct {
	ID        int                   `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	UserID    uuid.UUID             `gorm:"column:user_id;NOT NULL"`
	Command   string                `gorm:"column:command"`
	Arguments string                `gorm:"column:arguments"`
	Type      constants.MessageType `gorm:"column:type;NOT NULL"`
	ExtraData string                `gorm:"column:extra_data"`
	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy int64                 `gorm:"column:created_by;NOT NULL"`
	UpdatedBy int64                 `gorm:"column:updated_by"`
	Version   int                   `gorm:"column:version;default:1"`
	RawData   string                `gorm:"column:raw_data"`
}

func (m *CommandsLog) TableName() string {
	return "commands_log"
}
