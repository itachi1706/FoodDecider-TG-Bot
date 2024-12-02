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
	"github.com/google/uuid"
	"log"
)

const (
	SelectLocationRandom = "select-location-random"
)

func RandomLocationCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomLocation command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	err := utils.BasicReplyToUser(bot, ctx, "Please reply to this message with a location pin")
	if err != nil {
		log.Println("Error replying to user: " + err.Error())
		return handlers.EndConversation()
	}

	return handlers.NextConversationState(SelectLocationRandom)
}

func RandomLocationCommandLocationPin(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomLocation command with new location called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CONVERSATION)

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

	message, rollInfo := searchAndReplyWithLocation(ctx, pinLocation)
	if rollInfo == nil {
		_ = utils.BasicReplyToUser(bot, ctx, message)
		return handlers.EndConversation()
	}

	messageOpts := utils.GenerateRerollKeysSend(constants.LOCATION, *rollInfo, false)

	_ = utils.ReplyUserWithOpts(bot, ctx, message, messageOpts)
	return handlers.EndConversation()
}

func RandomLocationCommandReroll(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomLocation reroll button clicked by " + ctx.EffectiveSender.Username())

	rollInfo, count, cb, err := services.RerollCommon(bot, ctx, constants.LOCATION)
	if err != nil {
		return err
	}

	location := &gotgbot.Location{
		Latitude:  rollInfo.Latitude,
		Longitude: rollInfo.Longitude,
	}

	// Send message to user with reroll button
	message := sendWithRerollButtonLocation(*rollInfo, ctx.EffectiveSender, count, true, location)
	messageOpts := utils.GenerateRerollKeysEdit(constants.LOCATION, *rollInfo, false)

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: constants.RerollSuccess,
	})

	_, _, err = cb.Message.EditText(bot, message, messageOpts)

	return nil
}

func searchAndReplyWithLocation(ctx *ext.Context, pinLocation *gotgbot.Location) (string, *model.Rolls) {
	// Search groups and reply randomly
	sender := ctx.EffectiveSender
	userId := sender.Id()
	chatId := ctx.EffectiveChat.Id

	db := utils.GetDbConnection()

	log.Printf("Chat ID: %v, Latitude: %v, Longitude: %v\n", chatId, pinLocation.Latitude, pinLocation.Longitude)

	// Store in database
	rollInfo := model.Rolls{
		ID:                uuid.New(),
		Type:              constants.LOCATION,
		DecidedLocationID: nil,
		CreatedBy:         userId,
		UpdatedBy:         userId,
		ChatId:            chatId,
		Longitude:         pinLocation.Longitude,
		Latitude:          pinLocation.Latitude,
	}

	rollHistory, count, err := services.RollForFood(rollInfo)
	if err != nil {
		return utils.Capitalize(err.Error()), nil
	}
	rollInfo.DecidedFoodID = rollHistory.DecidedFoodID
	rollInfo.DecidedLocationID = rollHistory.DecidedLocationID

	// Save both roll and roll history
	db.Save(&rollInfo)
	db.Save(&rollHistory)

	// Send message to user with reroll button
	message := sendWithRerollButtonLocation(rollInfo, sender, count, false, pinLocation)

	return message, &rollInfo
}

func sendWithRerollButtonLocation(rollInfo model.Rolls, trigger *gotgbot.Sender, count int64, reroll bool, location *gotgbot.Location) string {
	foodRepo := repository.NewFoodsRepository(utils.GetDbConnection())
	food := foodRepo.FindFoodById(rollInfo.DecidedFoodID)
	loc := foodRepo.FindActiveLocationById(*rollInfo.DecidedLocationID)

	messageFmt := "Food Decision made 🎉\n\n"
	messageFmt += "Selected food: %s\n"
	messageFmt += "Description: %s\n"
	messageFmt += "Location: %f,%f\n\n"
	messageFmt += "This was randomized from a list of %d food options with the following coordinates:\n(%f,%f)\n\n"
	if reroll {
		messageFmt += "Decision was re-ran on %s by %s (%s)"
	} else {
		messageFmt += "Decision was ran on %s by %s (%s)"
	}

	updatedTime := rollInfo.UpdatedAt
	// Format the time to be more readable
	updatedTimeStr := updatedTime.Format(constants.DateTimeFormat)

	message := fmt.Sprintf(messageFmt, food.Name, food.Description, loc.Latitude, loc.Longitude, count, location.Latitude, location.Longitude, updatedTimeStr, trigger.Name(), trigger.Username())

	return message
}
