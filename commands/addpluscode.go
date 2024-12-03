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
	"strings"
)

func AddPlusCodeCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AddPlusCode command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId, foodId, messageOpts, err := services.FoodValidationParameterChecksAdmin(bot, ctx, 2, "Invalid format\n\nFormat: /addpluscode <food id> <plus code>")
	if err != nil {
		return err
	}

	plusCode := strings.Join(messageOpts[1:], " ")

	// Use geocoding api to reverse geocode and then geocode again
	geocodingAPI := services.NewGeocodingAPI()
	latitude, longitude, err := geocodingAPI.GetLocationFromPlusCode(plusCode)
	if err != nil {
		log.Println("Failed to get location from plus code: " + err.Error())
		return utils.BasicReplyToUser(bot, ctx, "Failed to get location from plus code. Please check you have a valid plus code")
	}

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
	log.Printf("Food ID: %v, Latitude: %v, Longitude: %v, Name: %v, Plus Code: %v\n", foodId, latitude, longitude, friendlyName, plusCode)
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
