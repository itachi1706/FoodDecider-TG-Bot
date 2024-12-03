package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func DelFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DelFood command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId, foodId, _, err := services.FoodValidationParameterChecksAdmin(bot, ctx, 1, "Please provide a food id to delete\n\nFormat: /delfood <id>")
	if err != nil {
		return err
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if food name already exists
	food := repo.FindFoodById(*foodId)
	message := constants.ErrorMessage
	if food == nil {
		// New Food
		message = "Food with ID " + foodId.String() + " does not exist"
	} else {
		food.Status = "D"
		food.UpdatedBy = *userId
		db.Save(&food)
		message = "Food " + food.Name + " deleted from database"
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
