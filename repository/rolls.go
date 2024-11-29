package repository

import (
	"FoodDecider-TG-Bot/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RollsRepository struct {
	db *gorm.DB
}

func NewRollsRepository(db *gorm.DB) RollsRepository {
	return RollsRepository{db: db}
}

func (f RollsRepository) FindRollsById(rollId uuid.UUID) *model.Rolls {
	var rolls model.Rolls
	result := f.db.Where(&model.Rolls{ID: rollId}).First(&rolls)
	if result.Error != nil {
		return nil
	}

	return &rolls
}

func (f RollsRepository) FindLastRollByChatId(chatId int64) *model.Rolls {
	var rolls model.Rolls
	result := f.db.Where(&model.Rolls{ChatId: chatId}).Order(clause.OrderByColumn{Column: clause.Column{Name: "updated_at"}, Desc: true}).First(&rolls)
	if result.Error != nil {
		return nil
	}

	return &rolls
}
