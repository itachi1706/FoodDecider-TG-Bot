package repository

import (
	"FoodDecider-TG-Bot/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (u UserRepository) FindUser(id int64) *model.Users {
	var user model.Users
	result := u.db.Where(&model.Users{TelegramID: id}).First(&user)
	if result.Error != nil {
		return nil
	}

	return &user
}
