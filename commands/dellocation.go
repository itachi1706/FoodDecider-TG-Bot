package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
)

func DelLocationCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DelLocation command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)
	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Invalid Format\n\nFormat: /dellocation <location id>")
	}

	locationId, err := uuid.Parse(messageOpts[0])
	if err != nil {
		return utils.BasicReplyToUser(bot, ctx, "Invalid location id provided")
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if location already exists
	location := repo.FindActiveLocationById(locationId)
	message := constants.ErrorMessage
	if location == nil {
		// New Food
		message = "Location ID " + locationId.String() + " does not exist\n\nUse /listlocations <food id> to get the location ID to delete"
	} else {
		log.Println("Deleting location " + locationId.String())
		location.Status = "D"
		db.Save(location)
		message = "Location ID " + locationId.String() + " deleted"
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
