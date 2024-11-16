package commands

import (
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func DelFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DelFood command called by " + ctx.EffectiveSender.Username())

	userId, foodId, _, err := services.FoodValidationParameterChecks(bot, ctx, 1, "Please provide a food id to delete\n\nFormat: /delfood <id>")
	if err != nil {
		return err
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if food name already exists
	food := repo.FindFoodById(*foodId)
	message := "An error has occurred. Please try again later"
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
