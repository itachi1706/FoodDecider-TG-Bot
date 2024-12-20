package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strings"
)

func UpdateFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("UpdateFood command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId, foodId, messageOpts, err := services.FoodValidationParameterChecksAdmin(bot, ctx, 3, "Invalid update food format\n\nFormat: /updatefood <id> <name/description> [value]")
	if err != nil {
		return err
	}

	updateType := strings.ToLower(messageOpts[1])
	log.Printf("Update type: '%s'\n", updateType)
	if updateType != "name" && updateType != "description" {
		return utils.BasicReplyToUser(bot, ctx, "Invalid update type provided. Valid types: name, description")
	}

	updateValue := strings.Trim(strings.Join(messageOpts[2:], " "), " ")

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if food name already exists
	food := repo.FindFoodById(*foodId)
	message := constants.ErrorMessage
	if food == nil {
		// New Food
		message = "Food with ID " + foodId.String() + " does not exist"
	} else {
		if updateType == "name" {
			food.Name = updateValue
		} else {
			food.Description = updateValue
		}

		food.UpdatedBy = *userId
		db.Save(&food)
		message = "Food " + food.Name + " updated in database"
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
