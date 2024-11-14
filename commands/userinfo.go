package commands

import (
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "log"
    "strconv"
)

func UserInfoCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("UserInfo command called by " + ctx.EffectiveSender.Username())
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
    debugText += "Username: " + username + "\n"
    debugText += "Name: " + senderName + "\n"
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

    return utils.BasicReplyToUserWithMarkdown(bot, ctx, debugText)
}
