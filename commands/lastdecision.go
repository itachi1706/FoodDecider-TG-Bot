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
	"log"
)

func LastDecisionCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("LastDecision command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	chatId := ctx.EffectiveChat.Id

	db := utils.GetDbConnection()
	rollsRepo := repository.NewRollsRepository(db)
	foodRepo := repository.NewFoodsRepository(db)

	lastRoll := rollsRepo.FindLastRollByChatId(chatId)
	if lastRoll == nil {
		return utils.BasicReplyToUser(bot, ctx, "No previous decision was found for this chat")
	}

	food := foodRepo.FindFoodById(lastRoll.DecidedFoodID)
	var location *model.Locations
	if lastRoll.DecidedLocationID != nil {
		location = foodRepo.FindActiveLocationById(*lastRoll.DecidedLocationID)
	}

	decisionMaker := lastRoll.UpdatedBy
	userRepo := repository.NewUserRepository(db)
	user := userRepo.FindUser(decisionMaker)

	messageFmt := "Last Decision Made: %s\n"
	if location != nil {
		messageFmt += "Location: %s (%f,%f)\n"
	}
	messageFmt += "\nMade by: %s (%s) on %s\nDecision ID: %s"
	
	var message string
	if location != nil {
		message = fmt.Sprintf(messageFmt, food.Name, location.Name, location.Latitude, location.Latitude, user.FullName, user.Username, lastRoll.UpdatedAt.Format(constants.DateTimeFormat), lastRoll.ID)
	} else {
		message = fmt.Sprintf(messageFmt, food.Name, user.FullName, user.Username, lastRoll.UpdatedAt.Format(constants.DateTimeFormat), lastRoll.ID)
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
