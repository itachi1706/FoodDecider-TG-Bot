package commands

import (
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strconv"
)

func ListAdminsCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListAdmins command called by " + ctx.EffectiveSender.Username())

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	db := utils.GetDbConnection()
	repo := repository.NewAdminsRepository(db)
	admins := repo.FindAllActiveAdmins()

	message := "No admins found"
	if len(admins) > 0 {
		message = "App Admins:\n"
		for idx, admin := range admins {
			message += strconv.Itoa(idx+1) + ": " + utils.EscapeMarkdown(admin.Name)
			if admin.IsSuperadmin {
				message += " ***(SUPERADMIN)***"
			}
			message += "\n"
		}
	}

	return utils.BasicReplyToUserWithMarkdown(bot, ctx, message)
}
