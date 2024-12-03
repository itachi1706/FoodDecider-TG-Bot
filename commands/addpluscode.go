package commands

import (
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
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
	message := services.AddLocationIfExist(*foodId, latitude, longitude, friendlyName, *userId, address)

	return utils.BasicReplyToUser(bot, ctx, message)
}
