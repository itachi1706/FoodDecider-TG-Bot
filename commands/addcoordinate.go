package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

func AddCoordinateCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AddCoordinate command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId, foodId, messageOpts, err := services.FoodValidationParameterChecksAdmin(bot, ctx, 3, "Invalid format\n\nFormat: /addcoordinate <food id> <latitude> <longitude> [name]")
	if err != nil {
		return err
	}

	latitude, err1 := strconv.ParseFloat(messageOpts[1], 64)
	longitude, err2 := strconv.ParseFloat(messageOpts[2], 64)
	if err1 != nil || err2 != nil {
		return utils.BasicReplyToUser(bot, ctx, "Invalid latitude or longitude provided")
	}

	// Ensure latitude and longitude within range
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return utils.BasicReplyToUser(bot, ctx, "Invalid latitude or longitude provided. Out of range")
	}

	// Use geocoding api
	geocodingAPI := services.NewGeocodingAPI()
	address, err := geocodingAPI.GetAddressFromLocation(latitude, longitude)
	if err != nil {
		log.Println("Failed to get address from location: " + err.Error())
		return utils.BasicReplyToUser(bot, ctx, "Failed to get address from location")
	}

	var friendlyName string
	if len(messageOpts) > 3 {
		// Has friendly name
		friendlyName = strings.Trim(strings.Join(messageOpts[3:], " "), " ")
	}
	log.Printf("Food ID: %v, Latitude: %v, Longitude: %v, Name: %v\n", foodId, latitude, longitude, friendlyName)
	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	location := repo.GetFoodLocation(*foodId, latitude, longitude)
	message := constants.ErrorMessage
	if location == nil {
		// New location
		log.Println("Creating new location for food " + foodId.String())
		location = &model.Locations{
			FoodID:    *foodId,
			Name:      friendlyName,
			Latitude:  latitude,
			Longitude: longitude,
			CreatedBy: *userId,
			UpdatedBy: *userId,
			ID:        uuid.New(),
			PlusCode:  address.PlusCode.GlobalCode,
			Address:   address.FormattedAddress,
		}
		db.Create(&location)
		message = "Location added for food " + foodId.String()
	} else {
		location.Name = friendlyName
		location.UpdatedBy = *userId
		message = "Location updated for food " + foodId.String()
		if location.Status != "A" {
			log.Println("Reactivating location for food " + foodId.String())
			location.Status = "A"
			message = "Location added for food " + foodId.String()
		}
		db.Save(&location)
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
