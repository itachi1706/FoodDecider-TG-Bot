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
	"github.com/google/uuid"
	"log"
	"strings"
)

const (
	AddFoodLocation = "add-food-location"
)

func AddLocationCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AddLocation command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	_, foodId, messageOpts, err := services.FoodValidationParameterChecks(bot, ctx, 1, "Invalid Format\n\nFormat: /addcoordinate <food id> [name]")
	if err != nil {
		return err
	}

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
	services.RunPreCommandScripts(ctx)

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

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	location := repo.GetFoodLocation(food.ID, pinLocation.Latitude, pinLocation.Longitude)
	message := constants.ErrorMessage
	if location == nil {
		// New location
		log.Println("Creating new location for food " + food.ID.String())
		location = &model.Locations{
			FoodID:    food.ID,
			Name:      name,
			Latitude:  pinLocation.Latitude,
			Longitude: pinLocation.Longitude,
			CreatedBy: userId,
			UpdatedBy: userId,
			ID:        uuid.New(),
		}
		db.Create(&location)
		message = "Location added for " + food.Name
	} else {
		location.Name = name
		location.UpdatedBy = userId
		message = "Location updated for " + food.Name
		if location.Status != "A" {
			log.Println("Reactivating location for food " + food.ID.String())
			location.Status = "A"
			message = "Location added for " + food.Name
		}
		db.Save(&location)
	}

	_ = utils.BasicReplyToUser(bot, ctx, message)
	return handlers.EndConversation()
}
