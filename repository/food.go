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

func (f FoodRepository) FindAllActiveFood() []model.Food {
	var foods []model.Food
	f.db.Where(&model.Food{Status: "A"}).Find(&foods)

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

func (f FoodRepository) GetFoodLocation(foodId uuid.UUID, lat float64, lng float64) *model.Locations {
	var location model.Locations
	result := f.db.Where(&model.Locations{FoodID: foodId, Latitude: lat, Longitude: lng}).First(&location)
	if result.Error != nil {
		return nil
	}

	return &location
}

func (f FoodRepository) FindAllLocationsForFoodPaginated(foodId uuid.UUID, size int, offset int) []model.Locations {
	var locations []model.Locations
	f.db.Where(&model.Locations{FoodID: foodId, Status: "A"}).Limit(size).Offset(offset * size).Find(&locations)

	return locations
}

func (f FoodRepository) FindAllLocationsForFoodCount(foodId uuid.UUID) int64 {
	var count int64
	f.db.Model(&model.Locations{}).Where(&model.Locations{FoodID: foodId, Status: "A"}).Count(&count)

	return count
}

func (f FoodRepository) FindActiveLocationById(id uuid.UUID) *model.Locations {
	var location model.Locations
	result := f.db.Where(&model.Locations{ID: id, Status: "A"}).First(&location)
	if result.Error != nil {
		return nil
	}

	return &location
}

func (f FoodRepository) FindAllFoodsFromGroups(groups []string) []model.Food {
	var foods []model.Food
	f.db.Joins("JOIN food_groups_link f ON f.food_id = food.id").Joins("JOIN food_groups g ON g.id = f.group_id").Where("g.name IN (?) AND f.status = ?", groups, "A").Distinct("food.id").Find(&foods)

	return foods
}
