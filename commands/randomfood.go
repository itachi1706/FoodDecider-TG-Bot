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
	message, hasLoc := sendWithRerollButton(rollInfo, sender, count, false)
	messageOpts := utils.GenerateRerollKeysSend(constants.GENERAL, rollInfo, hasLoc)

	return utils.ReplyUserWithOpts(bot, ctx, message, messageOpts)
}

func sendWithRerollButton(rollInfo model.Rolls, trigger *gotgbot.Sender, count int64, reroll bool) (string, bool) {
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
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)

	// Strip out "reroll-GENERAL-" from callback data to get roll UUID
	rollIdStr := cb.Data[len("reroll-GENERAL-"):]
	rollId, err := uuid.Parse(rollIdStr)
	if err != nil {
		_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: constants.RerollError,
		})
		return err
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
		return err
	}

	// Update roll info
	rollInfo.UpdatedBy = ctx.EffectiveSender.Id()
	rollHistory, count, err := services.RollForFood(*rollInfo)
	if err != nil {
		log.Println("Error re-rolling decision: " + err.Error())
		_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: constants.RerollError,
		})
		return err
	}

	rollInfo.DecidedFoodID = rollHistory.DecidedFoodID
	rollInfo.DecidedLocationID = nil
	db.Save(&rollInfo)
	db.Save(&rollHistory)

	// Send message to user with reroll button
	message, hasLoc := sendWithRerollButton(*rollInfo, ctx.EffectiveSender, count, true)
	messageOpts := utils.GenerateRerollKeysEdit(constants.GENERAL, *rollInfo, hasLoc)

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: constants.RerollSuccess,
	})

	_, _, err = cb.Message.EditText(bot, message, messageOpts)

	return nil
}
