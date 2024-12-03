package commands

import (
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strconv"
)

func UserInfoCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("UserInfo command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)
	sender := ctx.EffectiveSender

	username := sender.Username()
	if username == "" {
		username = "_Username Not Set_"
	}
	senderName := sender.Name()
	if senderName == "" {
		senderName = "_Name Not Set_"
	}

	debugText := "***User Information***\n"
	debugText += "Telegram ID: " + strconv.FormatInt(sender.Id(), 10) + "\n"
	debugText += "Username: " + utils.EscapeMarkdown(username) + "\n"
	debugText += "Name: " + utils.EscapeMarkdown(senderName) + "\n"
	debugText += "Is Bot: " + strconv.FormatBool(sender.IsBot()) + "\n"

	// Additional check if user is admin in app
	debugText += "\n***App Permissions***"
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

	opts := &gotgbot.SendMessageOpts{ParseMode: "Markdown", ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text: "Copy ID to Clipboard",
					CopyText: &gotgbot.CopyTextButton{
						Text: strconv.FormatInt(sender.Id(), 10),
					}},
			},
		},
	}}

	return utils.ReplyUserWithOpts(bot, ctx, debugText, opts)
}
