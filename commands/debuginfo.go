package commands

import (
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strconv"
)

func DebugInfoCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DebugInfo command called by " + ctx.EffectiveSender.Username())
	sender := ctx.EffectiveSender
	chat := ctx.EffectiveChat

	username := sender.Username()
	if username == "" {
		username = "Unknown Username"
	}
	senderName := sender.Name()
	if senderName != "" {
		username += " (" + senderName + ")"
	}

	group := username
	forum := false
	chatType := chat.Type
	if chatType != "private" {
		group = chat.Title
	}
	if chatType == "supergroup" || chatType == "group" {
		forum = chat.IsForum
	}

	debugText := "***Debug Information***\n"
	debugText += "Command Caller ID: " + strconv.FormatInt(sender.Id(), 10) + "\n"
	debugText += "Command Caller: " + utils.EscapeMarkdown(username) + "\n"
	debugText += "Chat ID: " + strconv.FormatInt(chat.Id, 10) + "\n"
	debugText += "Chat Type: " + chatType + "\n"
	debugText += "Is Forum Mode: " + strconv.FormatBool(forum) + "\n"
	debugText += "Chat Info: " + utils.EscapeMarkdown(group)

	// Additional check if user is admin in app
	debugText += "\n\n***App Permissions***"
	userId := sender.Id()
	if utils.CheckIfAdmin(userId) {
		debugText += "\nApp Admin: Yes"
	} else {
		debugText += "\nApp Admin: No"
	}

	if utils.CheckIfSuperAdmin(userId) {
		debugText += "\nApp Super Admin: Yes"
	} else {
		debugText += "\nApp Super Admin: No"
	}

	return utils.BasicReplyToUserWithMarkdown(bot, ctx, debugText)
}
