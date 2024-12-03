package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"log"
	"strings"
)

const (
	AddFoodLocation = "add-food-location"
)

func AddLocationCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AddLocation command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	_, foodId, messageOpts, err := services.FoodValidationParameterChecksAdmin(bot, ctx, 1, "Invalid Format\n\nFormat: /addlocation <food id> [name]")
	if err != nil {
		return err
	}
	log.Println("Why continue?")

	var friendlyName string
	if len(messageOpts) > 1 {
		// Has friendly name
		friendlyName = strings.Trim(strings.Join(messageOpts[1:], " "), " ")
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	food := repo.FindFoodById(*foodId)
	if food == nil {
		return utils.BasicReplyToUser(bot, ctx, "Food does not exist")
	}

	message := "Please reply to this messsage with a location pin for " + food.Name

	id, _ := conversation.KeyStrategySenderAndChat(ctx)
	err = utils.BasicReplyToUser(bot, ctx, message)
	services.SetString(fmt.Sprintf("foodloc-%s-id", id), food)
	services.SetString(fmt.Sprintf("foodloc-%s-name", id), friendlyName)
	return handlers.NextConversationState(AddFoodLocation)
}

func AddLocationCommandLocationPin(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AddLocation command with new location called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CONVERSATION)

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		_ = utils.BasicReplyToUser(bot, ctx, "This operation can only be done by an administrator. Use /cancel to cancel this operation")
		return handlers.NextConversationState(AddFoodLocation)
	}

	// Check if text is a location
	if ctx.EffectiveMessage.Location == nil {
		_ = utils.BasicReplyToUser(bot, ctx, "Invalid location provided. Please provide a location pin or use /cancel to cancel this operation")
		return handlers.NextConversationState(AddFoodLocation)
	}

	pinLocation := ctx.EffectiveMessage.Location

	// Make sure location is NOT a live location
	if pinLocation.LivePeriod != 0 {
		_ = utils.BasicReplyToUser(bot, ctx, "Invalid location provided. Do not use a live location. Please provide a location pin or use /cancel to cancel this operation")
		return handlers.NextConversationState(AddFoodLocation)
	}

	cid, _ := conversation.KeyStrategySenderAndChat(ctx)

	idKey := fmt.Sprintf("foodloc-%s-id", cid)
	nameKey := fmt.Sprintf("foodloc-%s-name", cid)

	foodIf, get := services.GetString(idKey)
	if !get {
		_ = utils.BasicReplyToUser(bot, ctx, constants.ErrorMessage)
		return handlers.EndConversation()
	}
	services.DeleteString(idKey)
	food := foodIf.(*model.Food)

	nameIf, get := services.GetString(nameKey)
	if !get {
		_ = utils.BasicReplyToUser(bot, ctx, constants.ErrorMessage)
		return handlers.EndConversation()
	}
	services.DeleteString(nameKey)
	name := nameIf.(string)

	log.Printf("Food ID: %v, Latitude: %v, Longitude: %v, Name: %v\n", food.ID, pinLocation.Latitude, pinLocation.Longitude, name)

	// Use geocoding api
	geocodingAPI := services.NewGeocodingAPI()
	address, err := geocodingAPI.GetAddressFromLocation(pinLocation.Latitude, pinLocation.Longitude)
	if err != nil {
		log.Println("Failed to get address from location: " + err.Error())
		return utils.BasicReplyToUser(bot, ctx, "Failed to get address from location")
	}

	log.Println("Address: " + address.FormattedAddress)

	message := services.AddLocationIfExist(food.ID, pinLocation.Latitude, pinLocation.Longitude, name, userId, address)

	_ = utils.BasicReplyToUser(bot, ctx, message)
	return handlers.EndConversation()
}
