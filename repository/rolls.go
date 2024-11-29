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

func (f RollsRepository) FindAllRollsByChatIdOrderRecentPaginated(chatId int64, size int, offset int) []model.Rolls {
	var rolls []model.Rolls
	f.db.Where(&model.Rolls{ChatId: chatId}).Order(clause.OrderByColumn{Column: clause.Column{Name: "updated_at"}, Desc: true}).Limit(size).Offset(offset * size).Find(&rolls)

	return rolls
}

func (f RollsRepository) GetRollsCountForChatId(chatId int64) int64 {
	var count int64
	f.db.Model(&model.Rolls{}).Where(&model.Rolls{ChatId: chatId}).Count(&count)

	return count
}

func (f RollsRepository) GetAllHistoryForRolls(rollId uuid.UUID) []model.RollsHistory {
	var history []model.RollsHistory

	// First record will be latest
	f.db.Where(&model.RollsHistory{RollID: rollId}).Order(clause.OrderByColumn{Column: clause.Column{Name: "updated_at"}, Desc: true}).Find(&history)

	return history

}
