package repository

import (
    "FoodDecider-TG-Bot/model"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

func FindFoodById(db *gorm.DB, foodId uuid.UUID) *model.Food {
    var food model.Food
    result := db.Where("id = ? AND status = ?", foodId, "A").First(&food)
    if result.Error != nil {
        return nil
    }

    return &food
}

func FindFoodByNameAll(db *gorm.DB, name string) *model.Food {
    var food model.Food
    result := db.Where("name = ?", name).First(&food)
    if result.Error != nil {
        return nil
    }

    return &food
}

func FindAllActiveFoodPaginated(db *gorm.DB, size int, offset int) []model.Food {
    var foods []model.Food
    db.Where("status = ?", "A").Limit(size).Offset(offset * size).Find(&foods)

    return foods
}
