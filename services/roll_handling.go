package services

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
)

func RollForFood(rollInfo model.Rolls) (*model.RollsHistory, int64, error) {
	log.Println("Rolling for food based on type: " + rollInfo.Type.String())

	db := utils.GetDbConnection()
	foodRepo := repository.NewFoodsRepository(db)
	choices := ""
	decidedFood := uuid.Nil
	var decidedLocation *uuid.UUID
	err := error(nil)
	count := int64(0)
	// Check type
	switch rollInfo.Type {
	case constants.GENERAL:
		decidedFood, choices, count, err = rollForGeneralFood(foodRepo, rollInfo)
	case constants.LOCATION:
		// decidedFood will be LOCATION here
		decidedFood, choices, count, err = rollForLocationFood(foodRepo, rollInfo)
	case constants.GROUP:
		decidedFood, choices, count, err = rollForGroupFood(foodRepo, rollInfo)
	case constants.LOCATION_GROUP:
		// decidedFood will be LOCATION here
		decidedFood, choices, count, err = rollForLocationGroupFood(foodRepo, rollInfo)
	}

	if err != nil {
		return nil, count, err
	}

	if rollInfo.Type == constants.LOCATION || rollInfo.Type == constants.LOCATION_GROUP {
		// Get food from location
		locationObj := foodRepo.FindActiveLocationById(decidedFood)
		if locationObj == nil {
			log.Println("Location not found. Ignoring and treating as no locations")
		} else {
			decidedLocation = &decidedFood
			decidedFood = locationObj.FoodID
		}
	}

	rollHistory := model.RollsHistory{
		RollID:            rollInfo.ID,
		CreatedBy:         rollInfo.UpdatedBy,
		UpdatedBy:         rollInfo.UpdatedBy,
		Choices:           choices,
		DecidedFoodID:     decidedFood,
		DecidedLocationID: decidedLocation,
	}

	return &rollHistory, count, nil
}

func rollForGeneralFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	log.Println("Rolling for general food")

	foods := foodRepo.FindAllActiveFood()
	count := foodRepo.GetFoodCount()

	if len(foods) == 0 {
		return uuid.Nil, "", count, errors.New("no food items to decide between")
	}

	// List choice name in a text string
	choiceString := ""
	for i, food := range foods {
		choiceString += strconv.Itoa(i+1) + ". " + food.Name + "\n"
	}

	// Make a decision with true rng
	randomIndex, err := GetTrueRandomNumber(0, int64(len(foods)-1))
	if err != nil {
		return uuid.Nil, "", count, errors.New("error getting random number. Please try again later")
	}

	selectedFoodObject := foods[randomIndex]

	return selectedFoodObject.ID, choiceString, count, nil
}

func rollForLocationFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	return rollForGeneralFood(foodRepo, rollInfo) // Coming Soon
}

func rollForGroupFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	return rollForGeneralFood(foodRepo, rollInfo) // Coming Soon
}

func rollForLocationGroupFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	return rollForGeneralFood(foodRepo, rollInfo) // Coming Soon
}
