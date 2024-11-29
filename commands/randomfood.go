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
	"github.com/google/uuid"
	"log"
)

func RandomFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomFood command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	sender := ctx.EffectiveSender
	userId := sender.Id()
	chatId := ctx.EffectiveChat.Id

	db := utils.GetDbConnection()

	// Store in database
	rollInfo := model.Rolls{
		ID:                uuid.New(),
		Type:              constants.GENERAL,
		DecidedLocationID: nil,
		CreatedBy:         userId,
		UpdatedBy:         userId,
		ChatId:            chatId,
	}

	rollHistory, count, err := services.RollForFood(rollInfo)
	if err != nil {
		return utils.BasicReplyToUser(bot, ctx, utils.Capitalize(err.Error()))
	}
	rollInfo.DecidedFoodID = rollHistory.DecidedFoodID

	// Save both roll and roll history
	db.Save(&rollInfo)
	db.Save(&rollHistory)

	// Send message to user with reroll button
	message, hasLoc := sendWithRerollButtonGeneral(rollInfo, sender, count, false)
	messageOpts := utils.GenerateRerollKeysSend(constants.GENERAL, rollInfo, hasLoc)

	return utils.ReplyUserWithOpts(bot, ctx, message, messageOpts)
}

func sendWithRerollButtonGeneral(rollInfo model.Rolls, trigger *gotgbot.Sender, count int64, reroll bool) (string, bool) {
	foodRepo := repository.NewFoodsRepository(utils.GetDbConnection())
	food := foodRepo.FindFoodById(rollInfo.DecidedFoodID)
	locationCnt := foodRepo.FindAllLocationsForFoodCount(rollInfo.DecidedFoodID)

	messageFmt := "Food Decision made 🎉\n\n"
	messageFmt += "Selected food: %s\n"
	messageFmt += "Description: %s\n\n"
	if locationCnt > 0 {
		messageFmt += "There are %d locations found for this option. Click the button below to view more\n\n"
	} else {
		messageFmt += "There are no locations found for this option. Please go online to find your nearest location yourself!\n\n"
	}
	messageFmt += "This was randomized from a list of %d food options\n\n"
	if reroll {
		messageFmt += "Decision was re-ran on %s by %s (%s)"
	} else {
		messageFmt += "Decision was ran on %s by %s (%s)"
	}

	updatedTime := rollInfo.UpdatedAt
	// Format the time to be more readable
	updatedTimeStr := updatedTime.Format(constants.DateTimeFormat)

	var message string
	if locationCnt > 0 {
		message = fmt.Sprintf(messageFmt, food.Name, food.Description, locationCnt, count, updatedTimeStr, trigger.Name(), trigger.Username())
	} else {
		message = fmt.Sprintf(messageFmt, food.Name, food.Description, count, updatedTimeStr, trigger.Name(), trigger.Username())
	}

	return message, locationCnt > 0
}

func RandomFoodCommandReroll(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomFood reroll button clicked by " + ctx.EffectiveSender.Username())
	rollInfo, count, cb, err := services.RerollCommon(bot, ctx, constants.GENERAL)
	if err != nil {
		return err
	}

	// Send message to user with reroll button
	message, hasLoc := sendWithRerollButtonGeneral(*rollInfo, ctx.EffectiveSender, count, true)
	messageOpts := utils.GenerateRerollKeysEdit(constants.GENERAL, *rollInfo, hasLoc)

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: constants.RerollSuccess,
	})

	_, _, err = cb.Message.EditText(bot, message, messageOpts)

	return nil
}
