package repository

import (
	"FoodDecider-TG-Bot/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
