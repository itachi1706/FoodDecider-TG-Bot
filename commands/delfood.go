package commands

import (
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
)

func DelFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DelFood command called by " + ctx.EffectiveSender.Username())

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)
	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Please provide a food id to delete\n\nFormat: /delfood <id>")
	}

	foodId, err := uuid.Parse(messageOpts[0])
	if err != nil {
		return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if food name already exists
	food := repo.FindFoodById(foodId)
	message := "An error has occurred. Please try again later"
	if food == nil {
		// New Food
		message = "Food with ID " + foodId.String() + " does not exist"
	} else {
		food.Status = "D"
		food.UpdatedBy = userId
		db.Save(&food)
		message = "Food " + food.Name + " deleted from database"
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
