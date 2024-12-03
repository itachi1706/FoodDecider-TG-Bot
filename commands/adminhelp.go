package commands

import (
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func AdminHelpCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AdminHelp command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	message := ""
	message += "***General Commands***\n"
	message += "/adminhelp - Show this message\n\n"

	message += "***Admin Management***\n"
	message += "/addadmin <telegram id> <nickname> - Add an admin\n"
	message += "/listadmins - List all admins\n"
	message += "/deladmin <telegram id> - Delete an admin\n\n"

	message += "***Food Management***\n"
	message += "/addfood <name> - Add a food to a group\n"
	message += "/updatefood <food id> <name/description> <value> - Update a food in a group\n"
	message += "/delfood <food id> - Delete a food from a group\n\n"

	message += "***Group Management***\n"
	message += "/addgroup <food id> <group name> - Add a group\n"
	message += "/renamegroup <group name> - Rename a group\n"
	message += "/removegroup <food id> <group name> - Remove a group\n\n"

	message += "***Location Management***\n"
	message += "/addcoordinate <food id> <latitude> <longitude> - Add a coordinate to a group\n"
	message += "/addlocation <food id> <name> - Add a location to a group\n"
	message += "/dellocation <location id> - Delete a coordinate from a group"

	return utils.BasicReplyToUserWithMarkdown(bot, ctx, message)
}
