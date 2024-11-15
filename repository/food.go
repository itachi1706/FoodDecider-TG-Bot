package repository

import (
    "FoodDecider-TG-Bot/model"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type FoodRepository struct {
    db *gorm.DB
}

func NewFoodsRepository(db *gorm.DB) FoodRepository {
    return FoodRepository{db: db}
}

func (f FoodRepository) FindFoodById(foodId uuid.UUID) *model.Food {
    var food model.Food
    result := f.db.Where("id = ? AND status = ?", foodId, "A").First(&food)
    if result.Error != nil {
        return nil
    }

    return &food
}

func (f FoodRepository) FindFoodByNameAll(name string) *model.Food {
    var food model.Food
    result := f.db.Where("name = ?", name).First(&food)
    if result.Error != nil {
        return nil
    }

    return &food
}

func (f FoodRepository) FindAllActiveFoodPaginated(size int, offset int) []model.Food {
    var foods []model.Food
    f.db.Where("status = ?", "A").Limit(size).Offset(offset * size).Find(&foods)

    return foods
}

func (f FoodRepository) GetFoodCount() int64 {
    var count int64
    f.db.Model(&model.Food{}).Where("status = ?", "A").Count(&count)

    return count
}

func (f FoodRepository) GetFoodGroup(name string) *model.FoodGroups {
    var foodGroup model.FoodGroups
    result := f.db.Where("name = ?", name).First(&foodGroup)
    if result.Error != nil {
        return nil
    }

    return &foodGroup
}

func (f FoodRepository) GetActiveFoodGroupById(id int) *model.FoodGroups {
    var foodGroup model.FoodGroups
    result := f.db.Where("id = ? AND status = ?", id, "A").First(&foodGroup)
    if result.Error != nil {
        return nil
    }

    return &foodGroup
}

func (f FoodRepository) GetActiveFoodGroup(name string) *model.FoodGroups {
    var foodGroup model.FoodGroups
    result := f.db.Where("name = ? AND status = ?", name, "A").First(&foodGroup)
    if result.Error != nil {
        return nil
    }

    return &foodGroup
}

func (f FoodRepository) GetFoodGroupLink(foodId uuid.UUID, groupId int) *model.FoodGroupsLink {
    var foodGroupLink model.FoodGroupsLink
    result := f.db.Where("food_id = ? AND group_id = ?", foodId, groupId).First(&foodGroupLink)
    if result.Error != nil {
        return nil
    }

    return &foodGroupLink

}
