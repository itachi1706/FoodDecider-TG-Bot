package services

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"errors"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
	"slices"
	"strconv"
	"strings"
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
			decidedLocation = &locationObj.ID
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

func rollCommon(foods []model.Food, count int64) (uuid.UUID, string, int64, error) {
	log.Println("Food Count to decide between: " + strconv.FormatInt(count, 10))

	if len(foods) == 0 {
		return uuid.Nil, "", count, errors.New("no food items to decide between")
	}

	if len(foods) == 1 {
		return foods[0].ID, foods[0].Name, count, nil // There is only 1 option. LOL
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

func rollForGeneralFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	log.Println("Rolling for general food")

	foods := foodRepo.FindAllActiveFood()
	count := foodRepo.GetFoodCount()

	return rollCommon(foods, count)
}

func rollForLocationFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	log.Println("Rolling for location food")

	// Get current location
	latitude := rollInfo.Latitude
	longitude := rollInfo.Longitude
	if latitude == 0 && longitude == 0 {
		return uuid.Nil, "", 0, errors.New("location not provided")
	}

	// Get all locations
	locations := foodRepo.FindAllActiveLocations()

	// Filter locations based on distance
	locationDistances := filterLocationsByDistance(locations, latitude, longitude)
	slices.SortFunc(locationDistances, func(a, b model.LocationDistance) int {
		if a.Distance < b.Distance {
			return -1
		}
		if a.Distance > b.Distance {
			return 1
		}
		return 0
	})

	// Get nearest 10 locations if there are 10 or more locations, else get all locations
	myNearestLocations := make([]model.LocationDistance, 0)
	if len(locationDistances) < 10 {
		myNearestLocations = locationDistances
	} else {
		myNearestLocations = locationDistances[:10]
	}

	// Randomize from this list
	randomIndex, err := GetTrueRandomNumber(0, int64(len(myNearestLocations)-1))
	if err != nil {
		return uuid.Nil, "", 0, errors.New("error getting random number. Please try again later")
	}

	// Get food from location and distance
	nearestLocation := myNearestLocations[randomIndex]

	// Get all food object from locations
	foodIds := make([]uuid.UUID, 0)
	for _, location := range myNearestLocations {
		foodIds = append(foodIds, location.Location.FoodID)
	}

	foodChoice := foodRepo.FindAllFoodsByIds(foodIds)
	// List choice name in a text string
	choiceString := ""
	for i, food := range foodChoice {
		choiceString += strconv.Itoa(i+1) + ". " + food.Name + "\n"
	}

	return nearestLocation.Location.ID, choiceString, int64(len(myNearestLocations)), nil
}

func filterLocationsByDistance(locations []model.Locations, latitude float64, longitude float64) []model.LocationDistance {
	// Return an array of locationdistance
	finalData := make([]model.LocationDistance, 0)
	for _, location := range locations {
		// Calculate distance via Vincenty method
		distance := utils.VincentyDistance(latitude, longitude, location.Latitude, location.Longitude)
		finalData = append(finalData, model.LocationDistance{
			Location: location,
			Distance: distance,
		})
	}

	return finalData
}

func rollForGroupFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	log.Println("Rolling for foods in groups")

	groupStr := rollInfo.GroupName
	groups := strings.Split(groupStr, "\n")

	foods := foodRepo.FindAllFoodsFromGroups(groups)
	count := len(foods)
	return rollCommon(foods, int64(count))
}

func rollForLocationGroupFood(foodRepo repository.FoodRepository, rollInfo model.Rolls) (uuid.UUID, string, int64, error) {
	return rollForGeneralFood(foodRepo, rollInfo) // Coming Soon
}

func RerollCommon(bot *gotgbot.Bot, ctx *ext.Context, rollType constants.DecisionType) (*model.Rolls, int64, *gotgbot.CallbackQuery, error) {
	RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)

	// Strip out "reroll-GROUP-" from callback data to get roll UUID
	rollIdStr := cb.Data[len("reroll-"+rollType.String()+"-"):]
	rollId, err := uuid.Parse(rollIdStr)
	if err != nil {
		_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: constants.RerollError,
		})
		return nil, 0, nil, err
	}

	// Get roll info
	db := utils.GetDbConnection()
	rollRepo := repository.NewRollsRepository(db)
	rollInfo := rollRepo.FindRollsById(rollId)

	// Check if roll exists
	if rollInfo == nil {
		_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: constants.RerollError,
		})
		return nil, 0, nil, err
	}

	// Update roll info
	rollInfo.UpdatedBy = ctx.EffectiveSender.Id()
	rollHistory, count, err := RollForFood(*rollInfo)
	if err != nil {
		log.Println("Error re-rolling decision: " + err.Error())
		_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: constants.RerollError,
		})
		return nil, 0, nil, err
	}

	rollInfo.DecidedFoodID = rollHistory.DecidedFoodID
	rollInfo.DecidedLocationID = rollHistory.DecidedLocationID
	db.Save(&rollInfo)
	db.Save(&rollHistory)

	return rollInfo, count, cb, nil
}
