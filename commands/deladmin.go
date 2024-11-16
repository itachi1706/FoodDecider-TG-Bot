package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strconv"
)

func DelAdminCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DelAdmin command called by " + ctx.EffectiveSender.Username())

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	// Get user information to add to admin list from message
	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)

	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Please provide a user id\n\nFormat: /deladmin <telegram id>")
	}

	userIdToDelStr := messageOpts[0]

	userIdToDel, err := strconv.ParseInt(userIdToDelStr, 10, 64)
	if err != nil {
		return utils.BasicReplyToUser(bot, ctx, "Invalid user ID format")
	}

	// Check if id to delete is superadmin
	if utils.CheckIfSuperAdmin(userIdToDel) {
		return utils.BasicReplyToUser(bot, ctx, "Cannot delete superadmin. Please disable super admin from DB")
	}

	db := utils.GetDbConnection()
	repo := repository.NewAdminsRepository(db)
	admin := repo.FindActiveAdmin(userIdToDel)
	message := constants.ErrorMessage
	if admin == nil {
		message = "User ID " + userIdToDelStr + " is not an administrator"
	} else {
		log.Println("Disabling user")
		admin.Status = "I"
		admin.UpdatedBy = userId

		db.Save(&admin)
		message = "User ID " + userIdToDelStr + " removed as admin"
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
