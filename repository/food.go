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
    result := f.db.Where(&model.Food{ID: foodId, Status: "A"}).First(&food)
    if result.Error != nil {
        return nil
    }

    return &food
}

func (f FoodRepository) FindFoodByNameAll(name string) *model.Food {
    var food model.Food
    result := f.db.Where(&model.Food{Name: name}).First(&food)
    if result.Error != nil {
        return nil
    }

    return &food
}

func (f FoodRepository) FindAllActiveFoodPaginated(size int, offset int) []model.Food {
    var foods []model.Food
    f.db.Where(&model.Food{Status: "A"}).Limit(size).Offset(offset * size).Find(&foods)

    return foods
}

func (f FoodRepository) GetFoodCount() int64 {
    var count int64
    f.db.Model(&model.Food{}).Where(&model.Food{Status: "A"}).Count(&count)

    return count
}

func (f FoodRepository) GetFoodGroup(name string) *model.FoodGroups {
    var foodGroup model.FoodGroups
    result := f.db.Where(&model.FoodGroups{Name: name}).First(&foodGroup)
    if result.Error != nil {
        return nil
    }

    return &foodGroup
}

func (f FoodRepository) GetActiveFoodGroupById(id int) *model.FoodGroups {
    var foodGroup model.FoodGroups
    result := f.db.Where(&model.FoodGroups{ID: id, Status: "A"}).First(&foodGroup)
    if result.Error != nil {
        return nil
    }

    return &foodGroup
}

func (f FoodRepository) GetActiveFoodGroup(name string) *model.FoodGroups {
    var foodGroup model.FoodGroups
    result := f.db.Where(&model.FoodGroups{Name: name, Status: "A"}).First(&foodGroup)
    if result.Error != nil {
        return nil
    }

    return &foodGroup
}

func (f FoodRepository) GetFoodGroupLink(foodId uuid.UUID, groupId int) *model.FoodGroupsLink {
    var foodGroupLink model.FoodGroupsLink
    result := f.db.Where(&model.FoodGroupsLink{FoodID: foodId, GroupID: groupId}).First(&foodGroupLink)
    if result.Error != nil {
        return nil
    }

    return &foodGroupLink
}

func (f FoodRepository) FindAllGroupsForFoodPaginated(foodId uuid.UUID, size int, offset int) []model.FoodGroups {
    var groups []model.FoodGroups
    f.db.Joins("JOIN food_groups_link f ON f.group_id = food_groups.id").Where("f.food_id = ? AND f.status = ?", foodId, "A").Limit(size).Offset(offset * size).Find(&groups)

    return groups
}

func (f FoodRepository) GetFoodGroupForFoodCount(foodId uuid.UUID) int64 {
    var count int64
    f.db.Model(&model.FoodGroups{}).Joins("JOIN food_groups_link f ON f.group_id = food_groups.id").Where("f.food_id = ? AND f.status = ?", foodId, "A").Count(&count)

    return count
}

func (f FoodRepository) GetActiveFoodGroupLink(foodId uuid.UUID, groupId int) *model.FoodGroupsLink {
    var foodGroupLink model.FoodGroupsLink
    result := f.db.Where(&model.FoodGroupsLink{FoodID: foodId, GroupID: groupId, Status: "A"}).First(&foodGroupLink)
    if result.Error != nil {
        return nil
    }

    return &foodGroupLink
}
